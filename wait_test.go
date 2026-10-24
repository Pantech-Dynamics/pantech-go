package pantech

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestWaitSucceeds(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"id":"op_1","kind":"start_instance","status":"submitting"}`),
		ok(`{"id":"op_1","kind":"start_instance","status":"submitted"}`),
		ok(`{"id":"op_1","kind":"start_instance","status":"succeeded","resource_id":"vm_1"}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	var slept []time.Duration
	c.sleep = func(ctx context.Context, d time.Duration) error { slept = append(slept, d); return nil }
	op, err := c.Operations.Wait(context.Background(), "op_1", &WaitOptions{Interval: 3 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if op.Status != OperationSucceeded || op.ResourceID != "vm_1" || rec.count() != 3 {
		t.Fatalf("op %+v after %d polls", op, rec.count())
	}
	if len(slept) != 2 || slept[0] != 3*time.Second {
		t.Fatalf("slept %v", slept)
	}
	if rec.request(0).URL.Path != "/public/v1/operations/op_1" {
		t.Fatalf("polled %s", rec.request(0).URL.Path)
	}
}

func TestWaitReturnsTheFailure(t *testing.T) {
	rec := newRecorder(t, ok(`{"id":"op_2","kind":"attach_volume","status":"failed","failure":{"code":"INSTANCE_BUSY","reason":"The instance kept changing. Retry when it has settled."}}`))
	c := newTestClientAt(t, rec.srv.URL)
	op, err := c.Operations.Wait(context.Background(), "op_2", nil)
	var opErr *OperationError
	if !errors.As(err, &opErr) || opErr.Operation.Failure.Code != "INSTANCE_BUSY" || op == nil {
		t.Fatalf("err = %v", err)
	}
	if !strings.Contains(err.Error(), "INSTANCE_BUSY") || !strings.Contains(err.Error(), "op_2") {
		t.Fatalf("message = %q", err)
	}
}

func TestWaitHonoursTheTimeout(t *testing.T) {
	rec := newRecorder(t, ok(`{"id":"op_3","status":"submitted"}`))
	c := newTestClientAt(t, rec.srv.URL)
	c.sleep = sleepContext
	start := time.Now()
	_, err := c.Operations.Wait(context.Background(), "op_3", &WaitOptions{Interval: 5 * time.Millisecond, Timeout: 60 * time.Millisecond})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v", err)
	}
	if time.Since(start) > 5*time.Second {
		t.Fatal("the wait outlived its timeout")
	}
}

func TestWaitStopsOnAnAPIError(t *testing.T) {
	rec := newRecorder(t, status(404, `{"code":"RESOURCE_NOT_FOUND"}`))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.Operations.Wait(context.Background(), "op_x", nil); !IsCode(err, "RESOURCE_NOT_FOUND") {
		t.Fatalf("err = %v", err)
	}
}

func TestInstanceOrderWait(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"id":"ord_1","instance_id":"vm_1","status":"awaiting_payment","amount_minor":100,"currency":"NGN"}`),
		ok(`{"id":"ord_1","instance_id":"vm_1","status":"provisioning","amount_minor":100,"currency":"NGN"}`),
		ok(`{"id":"ord_1","instance_id":"vm_1","status":"provisioned","operation_id":"op_9","amount_minor":100,"currency":"NGN"}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	order, err := c.Instances.WaitForOrder(context.Background(), "ord_1", nil)
	if err != nil || order.Status != OrderProvisioned || *order.OperationID != "op_9" {
		t.Fatalf("order %+v err %v", order, err)
	}
	if rec.request(0).URL.Path != "/public/v1/instance-orders/ord_1" {
		t.Fatalf("polled %s", rec.request(0).URL.Path)
	}
}

func TestDatabaseOrderWaitFails(t *testing.T) {
	for _, st := range []string{OrderFailed, OrderPaymentFailed} {
		rec := newRecorder(t, ok(`{"id":"ord_2","database_id":"db_1","status":"`+st+`","failure_code":"card_declined","amount_minor":1,"currency":"USD"}`))
		c := newTestClientAt(t, rec.srv.URL)
		_, err := c.Databases.WaitForOrder(context.Background(), "ord_2", nil)
		var orderErr *OrderError
		if !errors.As(err, &orderErr) || orderErr.Status != st || orderErr.FailureCode != "card_declined" || orderErr.Kind != "database" {
			t.Fatalf("%s: err = %v", st, err)
		}
		if rec.request(0).URL.Path != "/public/v1/database-orders/ord_2" {
			t.Fatalf("polled %s", rec.request(0).URL.Path)
		}
	}
}
