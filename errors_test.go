package pantech

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestProblemParsing(t *testing.T) {
	body := `{"type":"https://api.pantechdynamics.com/problems/validation-failed","title":"Validation failed","status":422,"code":"VALIDATION_FAILED","detail":"plan_slug: No active plan has this slug.","request_id":"req_abc","errors":[{"field":"plan_slug","code":"PLAN_NOT_FOUND","message":"No active plan has this slug."}]}`
	rec := newRecorder(t, status(422, body))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.Instances.Resize(context.Background(), "vm_1", ResizeRequest{PlanSlug: "nope"})
	if !IsCode(err, "VALIDATION_FAILED") || IsCode(err, "NOT_IT") {
		t.Fatalf("IsCode failed for %v", err)
	}
	e, ok := AsError(fmt.Errorf("resizing: %w", err))
	if !ok {
		t.Fatal("AsError did not find the *Error through wrapping")
	}
	if e.Status != 422 || e.Title != "Validation failed" || e.RequestID != "req_abc" || e.Type == "" {
		t.Fatalf("error = %#v", e)
	}
	if len(e.Errors) != 1 || !e.HasField("plan_slug", "PLAN_NOT_FOUND") {
		t.Fatalf("field errors = %+v", e.Errors)
	}
	msg := err.Error()
	for _, want := range []string{"Validation failed", "VALIDATION_FAILED", "status 422", "req_abc", "plan_slug: No active plan has this slug."} {
		if !strings.Contains(msg, want) {
			t.Errorf("message %q lacks %q", msg, want)
		}
	}
	if strings.Count(msg, "No active plan") != 1 {
		t.Errorf("the field error is repeated: %q", msg)
	}
}

func TestProblemDetailLeads(t *testing.T) {
	rec := newRecorder(t, status(404, `{"title":"Resource not found","status":404,"code":"RESOURCE_NOT_FOUND","detail":"No instance with this id in your organization.","request_id":"req_1"}`))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.Instances.Get(context.Background(), "vm_x")
	want := "pantech: No instance with this id in your organization. (RESOURCE_NOT_FOUND, status 404, request req_1)"
	if err == nil || err.Error() != want {
		t.Fatalf("got %q\nwant %q", err, want)
	}
}

func TestNonProblemBodies(t *testing.T) {
	cases := []struct {
		name   string
		resp   cannedResponse
		detail string
	}{
		{"html page", cannedResponse{status: 502, body: "<html><body>Temporarily unavailable</body></html>"}, "The API is temporarily unavailable (502 Bad Gateway). Try again in a moment."},
		{"short text", cannedResponse{status: 400, body: "bad request"}, "bad request"},
		{"empty", cannedResponse{status: 403, body: ""}, ""},
	}
	for _, tc := range cases {
		rec := newRecorder(t, tc.resp)
		c := newTestClientAt(t, rec.srv.URL, WithMaxRetries(0))
		err := c.Do(context.Background(), "GET", "/x", nil, nil, nil)
		e, ok := AsError(err)
		if !ok {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if e.Detail != tc.detail || e.Code != "" || strings.Contains(err.Error(), "<html") {
			t.Errorf("%s: detail %q, error %q", tc.name, e.Detail, err)
		}
	}
}

func TestRequestIDFromHeader(t *testing.T) {
	rec := newRecorder(t, cannedResponse{status: 401, body: `{"code":"UNAUTHENTICATED","detail":"The API key is not valid."}`, header: map[string]string{"X-Request-Id": "req_hdr"}})
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.Me.Get(context.Background())
	if e, ok := AsError(err); !ok || e.RequestID != "req_hdr" || e.Status != 401 {
		t.Fatalf("err = %#v", err)
	}
}

func TestIsCodeOnOtherErrors(t *testing.T) {
	if IsCode(nil, "X") || IsCode(errors.New("x"), "X") {
		t.Fatal("IsCode matched a non-API error")
	}
}
