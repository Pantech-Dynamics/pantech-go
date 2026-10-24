package pantech

import (
	"context"
	"testing"
)

func TestStreamYieldsOperations(t *testing.T) {
	body := ": keepalive\n\n" +
		"event: operation\ndata: {\"id\":\"op_1\",\"status\":\"submitted\"}\n\n" +
		": keepalive\n\n" +
		"event: operation\ndata: {\"id\":\"op_1\",\"status\":\"succeeded\"}\n\n"
	rec := newRecorder(t, cannedResponse{status: 200, body: body, header: map[string]string{"Content-Type": "text/event-stream"}})
	c := newTestClientAt(t, rec.srv.URL)
	var statuses []string
	for op, err := range c.Operations.Stream(context.Background(), "op_1") {
		if err != nil {
			t.Fatal(err)
		}
		statuses = append(statuses, op.Status)
	}
	if len(statuses) != 2 || statuses[1] != OperationSucceeded {
		t.Fatalf("statuses = %v", statuses)
	}
	r := rec.request(0)
	if r.URL.Path != "/public/v1/operations/op_1/events" || r.Header.Get("Accept") != "text/event-stream" {
		t.Fatalf("sent %s with Accept %q", r.URL.Path, r.Header.Get("Accept"))
	}
}

func TestStreamErrorEvent(t *testing.T) {
	body := "event: error\ndata: {\"code\":\"OPERATION_STREAM_TIMEOUT\",\"detail\":\"Reconnect.\",\"request_id\":\"req_1\"}\n\n"
	rec := newRecorder(t, cannedResponse{status: 200, body: body, header: map[string]string{"Content-Type": "text/event-stream"}})
	c := newTestClientAt(t, rec.srv.URL)
	var last error
	for _, err := range c.Operations.Stream(context.Background(), "op_1") {
		last = err
	}
	if !IsCode(last, "OPERATION_STREAM_TIMEOUT") {
		t.Fatalf("err = %v", last)
	}
}

func TestStreamRefused(t *testing.T) {
	rec := newRecorder(t, status(404, `{"code":"RESOURCE_NOT_FOUND"}`))
	c := newTestClientAt(t, rec.srv.URL)
	for _, err := range c.Operations.Stream(context.Background(), "op_x") {
		if !IsCode(err, "RESOURCE_NOT_FOUND") {
			t.Fatalf("err = %v", err)
		}
	}
}
