package pantech

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

const sampleDynamicIPJSON = `{
	"instance_id": "vm_1",
	"schedule": "daily",
	"connections": "keep",
	"on_demand": true,
	"state": "active",
	"current_address": "198.51.100.42",
	"next_rotation_at": "2026-11-25T12:00:00Z",
	"last_rotated_at": "2026-11-24T12:00:00Z",
	"surcharge_percent": 95,
	"on_demand_used": 2,
	"on_demand_limit": 10,
	"failure_reason": null,
	"events": [
		{
			"kind": "on_demand",
			"from_address": "198.51.100.10",
			"to_address": "198.51.100.42",
			"degraded": false,
			"occurred_at": "2026-11-24T12:00:00Z"
		},
		{
			"kind": "enabled",
			"from_address": null,
			"to_address": "198.51.100.10",
			"degraded": false,
			"occurred_at": "2026-11-23T12:00:00Z"
		}
	]
}`

func TestInstancesGetDynamicIP(t *testing.T) {
	rec := newRecorder(t, ok(sampleDynamicIPJSON))
	c := newTestClientAt(t, rec.srv.URL)

	dip, err := c.Instances.GetDynamicIP(context.Background(), "vm_1")
	if err != nil {
		t.Fatal(err)
	}
	req := rec.request(0)
	if req.Method != "GET" || req.URL.Path != "/public/v1/instances/vm_1/address-rotation" {
		t.Fatalf("got %s %s", req.Method, req.URL.Path)
	}
	if dip.InstanceID != "vm_1" || dip.Schedule != "daily" || dip.Connections != "keep" || !dip.OnDemand || dip.State != "active" {
		t.Fatalf("unexpected DynamicIP: %+v", dip)
	}
	if dip.CurrentAddress == nil || *dip.CurrentAddress != "198.51.100.42" {
		t.Fatalf("CurrentAddress = %v", dip.CurrentAddress)
	}
	wantNext := time.Date(2026, 11, 25, 12, 0, 0, 0, time.UTC)
	wantLast := time.Date(2026, 11, 24, 12, 0, 0, 0, time.UTC)
	if dip.NextRotationAt == nil || !dip.NextRotationAt.Equal(wantNext) {
		t.Fatalf("NextRotationAt = %v, want %v", dip.NextRotationAt, wantNext)
	}
	if dip.LastRotatedAt == nil || !dip.LastRotatedAt.Equal(wantLast) {
		t.Fatalf("LastRotatedAt = %v, want %v", dip.LastRotatedAt, wantLast)
	}
	if dip.SurchargePercent != 95 || dip.OnDemandUsed != 2 || dip.OnDemandLimit != 10 || dip.FailureReason != nil {
		t.Fatalf("unexpected metrics/failure: %+v", dip)
	}
	if len(dip.Events) != 2 {
		t.Fatalf("len(Events) = %d, want 2", len(dip.Events))
	}
	ev0 := dip.Events[0]
	if ev0.Kind != "on_demand" || ev0.FromAddress == nil || *ev0.FromAddress != "198.51.100.10" || ev0.ToAddress == nil || *ev0.ToAddress != "198.51.100.42" || ev0.Degraded || !ev0.OccurredAt.Equal(wantLast) {
		t.Fatalf("unexpected event[0]: %+v", ev0)
	}
	ev1 := dip.Events[1]
	if ev1.Kind != "enabled" || ev1.FromAddress != nil || ev1.ToAddress == nil || *ev1.ToAddress != "198.51.100.10" {
		t.Fatalf("unexpected event[1]: %+v", ev1)
	}
}

func TestInstancesSetDynamicIP(t *testing.T) {
	t.Run("enables with 201", func(t *testing.T) {
		rec := newRecorder(t, status(201, `{
			"instance_id": "vm_1",
			"schedule": "every_12_hours",
			"connections": "move",
			"on_demand": true,
			"state": "enabling",
			"current_address": null,
			"next_rotation_at": null,
			"last_rotated_at": null,
			"surcharge_percent": 95,
			"on_demand_used": 0,
			"on_demand_limit": 10,
			"failure_reason": null,
			"events": []
		}`))
		c := newTestClientAt(t, rec.srv.URL)

		dip, err := c.Instances.SetDynamicIP(context.Background(), "vm_1", PutDynamicIPRequest{
			Schedule:    "every_12_hours",
			Connections: Ptr("move"),
			OnDemand:    Ptr(true),
		})
		if err != nil {
			t.Fatal(err)
		}
		req := rec.request(0)
		if req.Method != "PUT" || req.URL.Path != "/public/v1/instances/vm_1/address-rotation" {
			t.Fatalf("got %s %s", req.Method, req.URL.Path)
		}
		if req.Header.Get("Idempotency-Key") == "" {
			t.Fatal("missing Idempotency-Key")
		}
		wantBody := `{"schedule":"every_12_hours","connections":"move","on_demand":true}`
		if rec.bodies[0] != wantBody {
			t.Fatalf("body = %s, want %s", rec.bodies[0], wantBody)
		}
		if dip.State != "enabling" || dip.CurrentAddress != nil || dip.NextRotationAt != nil || dip.LastRotatedAt != nil {
			t.Fatalf("unexpected DynamicIP: %+v", dip)
		}
	})

	t.Run("updates with 200", func(t *testing.T) {
		rec := newRecorder(t, ok(sampleDynamicIPJSON))
		c := newTestClientAt(t, rec.srv.URL)

		dip, err := c.Instances.SetDynamicIP(context.Background(), "vm_1", PutDynamicIPRequest{
			Schedule: "weekly",
		})
		if err != nil {
			t.Fatal(err)
		}
		if rec.bodies[0] != `{"schedule":"weekly"}` {
			t.Fatalf("body = %s, want {\"schedule\":\"weekly\"}", rec.bodies[0])
		}
		if dip.InstanceID != "vm_1" || dip.State != "active" {
			t.Fatalf("unexpected DynamicIP: %+v", dip)
		}
	})
}

func TestInstancesDisableDynamicIP(t *testing.T) {
	rec := newRecorder(t, ok(`{
		"instance_id": "vm_1",
		"schedule": "daily",
		"connections": "keep",
		"on_demand": false,
		"state": "disabling",
		"current_address": "198.51.100.42",
		"next_rotation_at": null,
		"last_rotated_at": "2026-11-24T12:00:00Z",
		"surcharge_percent": 0,
		"on_demand_used": 0,
		"on_demand_limit": 10,
		"failure_reason": null,
		"events": []
	}`))
	c := newTestClientAt(t, rec.srv.URL)

	dip, err := c.Instances.DisableDynamicIP(context.Background(), "vm_1")
	if err != nil {
		t.Fatal(err)
	}
	req := rec.request(0)
	if req.Method != "DELETE" || req.URL.Path != "/public/v1/instances/vm_1/address-rotation" {
		t.Fatalf("got %s %s", req.Method, req.URL.Path)
	}
	if req.Header.Get("Idempotency-Key") == "" {
		t.Fatal("missing Idempotency-Key")
	}
	if dip.State != "disabling" || dip.SurchargePercent != 0 {
		t.Fatalf("unexpected DynamicIP: %+v", dip)
	}
}

func TestInstancesRotateDynamicIP(t *testing.T) {
	rec := newRecorder(t, ok(sampleDynamicIPJSON))
	c := newTestClientAt(t, rec.srv.URL)

	dip, err := c.Instances.RotateDynamicIP(context.Background(), "vm_1")
	if err != nil {
		t.Fatal(err)
	}
	req := rec.request(0)
	if req.Method != "POST" || req.URL.Path != "/public/v1/instances/vm_1/address-rotation/rotate" {
		t.Fatalf("got %s %s", req.Method, req.URL.Path)
	}
	if req.Header.Get("Idempotency-Key") == "" {
		t.Fatal("missing Idempotency-Key")
	}
	if dip.CurrentAddress == nil || *dip.CurrentAddress != "198.51.100.42" || dip.OnDemandUsed != 2 {
		t.Fatalf("unexpected DynamicIP: %+v", dip)
	}
}

func TestInstancesCreateWithDynamicIP(t *testing.T) {
	rec := newRecorder(t,
		status(202, `{"order_id":"ord_1","instance_id":"vm_1","resource_id":"vm_1","status":"awaiting_payment"}`),
		status(202, `{"order_id":"ord_2","instance_id":"vm_2","resource_id":"vm_2","status":"awaiting_payment"}`),
	)
	c := newTestClientAt(t, rec.srv.URL)

	_, err := c.Instances.Create(context.Background(), CreateInstanceRequest{
		Name:      "web-1",
		PlanSlug:  "standard-1",
		ImageSlug: "ubuntu-24-04",
		SSHKeyID:  "sshk_1",
		DynamicIP: &CreateInstanceDynamicIP{
			Schedule:    "daily",
			Connections: Ptr("move"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(rec.bodies[0]), &payload); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	ar, ok := payload["address_rotation"].(map[string]any)
	if !ok {
		t.Fatalf("address_rotation missing from body: %s", rec.bodies[0])
	}
	if ar["schedule"] != "daily" || ar["connections"] != "move" {
		t.Fatalf("address_rotation = %v", ar)
	}

	// Without DynamicIP, address_rotation is omitted.
	_, err = c.Instances.Create(context.Background(), CreateInstanceRequest{
		Name:      "web-2",
		PlanSlug:  "standard-1",
		ImageSlug: "ubuntu-24-04",
		SSHKeyID:  "sshk_1",
	})
	if err != nil {
		t.Fatal(err)
	}
	var payload2 map[string]any
	if err := json.Unmarshal([]byte(rec.bodies[1]), &payload2); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	if _, exists := payload2["address_rotation"]; exists {
		t.Fatalf("address_rotation should be omitted when nil: %s", rec.bodies[1])
	}
}
