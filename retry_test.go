package pantech

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const problem500 = `{"type":"about:blank","title":"Internal error","status":500,"code":"INTERNAL","detail":"Something failed on our side.","request_id":"req_1"}`

func TestRetriesTransientStatuses(t *testing.T) {
	for _, code := range []int{429, 500, 502, 503, 504} {
		rec := newRecorder(t, status(code, `{}`), status(code, `{}`), status(202, `{"operation_id":"op_1","resource_id":"vm_1","status":"submitting"}`))
		c := newTestClientAt(t, rec.srv.URL)
		acc, err := c.Instances.Start(context.Background(), "vm_1")
		if err != nil {
			t.Fatalf("%d: %v", code, err)
		}
		if acc.OperationID != "op_1" || rec.count() != 3 {
			t.Fatalf("%d: %d requests, accepted %+v", code, rec.count(), acc)
		}
	}
}

func TestIdempotencyKeyIsReusedAcrossRetries(t *testing.T) {
	rec := newRecorder(t, status(503, `{}`), status(500, problem500), status(202, `{}`))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.Volumes.Create(context.Background(), CreateVolumeRequest{Name: "data", DiskOfferingSlug: "small-5gb"}); err != nil {
		t.Fatal(err)
	}
	key := rec.request(0).Header.Get("Idempotency-Key")
	if key == "" {
		t.Fatal("no Idempotency-Key")
	}
	for i := 1; i < rec.count(); i++ {
		if got := rec.request(i).Header.Get("Idempotency-Key"); got != key {
			t.Fatalf("retry %d sent key %q, first try sent %q", i, got, key)
		}
		if rec.bodies[i] != rec.bodies[0] {
			t.Fatalf("retry %d sent body %q, first try %q", i, rec.bodies[i], rec.bodies[0])
		}
	}
	// A second logical call gets a new key.
	if _, err := c.Volumes.Create(context.Background(), CreateVolumeRequest{Name: "data", DiskOfferingSlug: "small-5gb"}); err != nil {
		t.Fatal(err)
	}
	if rec.request(rec.count()-1).Header.Get("Idempotency-Key") == key {
		t.Fatal("a new call reused the previous call's Idempotency-Key")
	}
}

func TestGivesUpAfterMaxRetries(t *testing.T) {
	rec := newRecorder(t, status(500, problem500))
	c := newTestClientAt(t, rec.srv.URL, WithMaxRetries(2))
	_, err := c.Instances.Get(context.Background(), "vm_1")
	if !IsCode(err, "INTERNAL") {
		t.Fatalf("err = %v", err)
	}
	if rec.count() != 3 {
		t.Fatalf("%d requests, want 3 (1 + 2 retries)", rec.count())
	}
}

func TestDefaultRetries(t *testing.T) {
	rec := newRecorder(t, status(503, `{}`))
	c := newTestClientAt(t, rec.srv.URL)
	_, _ = c.Instances.Get(context.Background(), "vm_1")
	if rec.count() != 1+DefaultMaxRetries {
		t.Fatalf("%d requests, want %d", rec.count(), 1+DefaultMaxRetries)
	}
}

func TestDoesNotRetryClientErrors(t *testing.T) {
	for _, code := range []int{400, 401, 403, 404, 409, 422} {
		rec := newRecorder(t, status(code, `{"code":"X","status":0}`))
		c := newTestClientAt(t, rec.srv.URL)
		if _, err := c.Instances.Start(context.Background(), "vm_1"); err == nil {
			t.Fatalf("%d: no error", code)
		}
		if rec.count() != 1 {
			t.Fatalf("%d: %d requests", code, rec.count())
		}
	}
}

func TestWithNoRetry(t *testing.T) {
	rec := newRecorder(t, status(429, `{}`))
	c := newTestClientAt(t, rec.srv.URL)
	_ = c.Do(context.Background(), "GET", "/regions", nil, nil, nil, WithNoRetry())
	if rec.count() != 1 {
		t.Fatalf("%d requests", rec.count())
	}
}

func TestRoutesThatDoNotReplayAreNotRetried(t *testing.T) {
	ctx := context.Background()
	calls := map[string]func(c *Client) error{
		"POST /ssh-keys": func(c *Client) error {
			_, err := c.SSHKeys.Create(ctx, CreateSSHKeyRequest{Name: "k"})
			return err
		},
		"DELETE /ssh-keys/{id}":        func(c *Client) error { return c.SSHKeys.Delete(ctx, "sshk_1") },
		"POST /instances/{id}/console": func(c *Client) error { _, err := c.Instances.OpenConsole(ctx, "vm_1"); return err },
		"POST /webhooks/{id}/test":     func(c *Client) error { _, err := c.Webhooks.Test(ctx, "wh_1"); return err },
		"DELETE /webhooks/{id}":        func(c *Client) error { return c.Webhooks.Delete(ctx, "wh_1") },
		"POST /webhooks/{id}/rotate-secret": func(c *Client) error {
			_, err := c.Webhooks.RotateSecret(ctx, "wh_1")
			return err
		},
		"POST /webhooks/{id}/deliveries/{id}/redeliver": func(c *Client) error {
			_, err := c.Webhooks.Redeliver(ctx, "wh_1", "whd_1")
			return err
		},
	}
	for name, call := range calls {
		for _, code := range []int{500, 502, 503, 504} {
			rec := newRecorder(t, status(code, `{}`))
			if err := call(newTestClientAt(t, rec.srv.URL)); err == nil {
				t.Fatalf("%s %d: no error", name, code)
			}
			if rec.count() != 1 {
				t.Errorf("%s was retried after %d (%d requests)", name, code, rec.count())
			}
		}
		// A 429 was refused before anything ran, so it is retried.
		rec := newRecorder(t, status(429, `{}`), status(200, `{}`))
		if err := call(newTestClientAt(t, rec.srv.URL)); err != nil {
			t.Fatalf("%s after 429: %v", name, err)
		}
		if rec.count() != 2 {
			t.Errorf("%s 429: %d requests, want 2", name, rec.count())
		}
	}
	// The same paths with other methods still retry.
	rec := newRecorder(t, status(503, `{}`), ok(`{"id":"wh_1"}`))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.Webhooks.Get(ctx, "wh_1"); err != nil || rec.count() != 2 {
		t.Errorf("GET /webhooks/{id}: err %v after %d requests, want a retry", err, rec.count())
	}
}

func TestNetworkErrorsAreRetried(t *testing.T) {
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if n.Add(1) == 1 {
			hj, _ := w.(http.Hijacker)
			conn, _, _ := hj.Hijack()
			conn.Close() // drop the connection without an answer
			return
		}
		_, _ = w.Write([]byte(`{"operation_id":"op_1"}`))
	}))
	defer srv.Close()
	c := newTestClientAt(t, srv.URL)
	acc, err := c.Instances.Reboot(context.Background(), "vm_1")
	if err != nil || acc.OperationID != "op_1" || n.Load() != 2 {
		t.Fatalf("err %v, %d requests", err, n.Load())
	}
}

func TestRetryAfterIsHonoured(t *testing.T) {
	rec := newRecorder(t,
		cannedResponse{status: 429, body: `{"code":"RATE_LIMITED"}`, header: map[string]string{"Retry-After": "7"}},
		ok(`{}`))
	c := newTestClientAt(t, rec.srv.URL)
	var mu sync.Mutex
	var waits []time.Duration
	c.sleep = func(ctx context.Context, d time.Duration) error {
		mu.Lock()
		waits = append(waits, d)
		mu.Unlock()
		return nil
	}
	if err := c.Do(context.Background(), "GET", "/regions", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if len(waits) != 1 || waits[0] != 7*time.Second {
		t.Fatalf("waits = %v, want [7s]", waits)
	}
}

func TestRetryAfterOnTheError(t *testing.T) {
	rec := newRecorder(t, cannedResponse{status: 429, body: `{"code":"RATE_LIMITED","detail":"Slow down."}`, header: map[string]string{"Retry-After": "12"}})
	c := newTestClientAt(t, rec.srv.URL, WithMaxRetries(0))
	err := c.Do(context.Background(), "GET", "/regions", nil, nil, nil)
	e, ok := AsError(err)
	if !ok || e.RetryAfter != 12*time.Second || e.Status != 429 {
		t.Fatalf("err = %#v", err)
	}
}

func TestBackoffGrowsAndIsCapped(t *testing.T) {
	c := newTestClientAt(t, "http://localhost:1")
	want := []time.Duration{500 * time.Millisecond, time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 8 * time.Second}
	for i, w := range want {
		if got := c.backoff(i, 0); got != w {
			t.Errorf("backoff(%d) = %v, want %v", i, got, w)
		}
	}
	if got := c.backoff(0, 10*time.Minute); got != maxRetryAfter {
		t.Errorf("a long Retry-After was not capped: %v", got)
	}
	for i := 0; i < 100; i++ {
		if d := halfJitter(time.Second); d < 500*time.Millisecond || d > time.Second {
			t.Fatalf("halfJitter = %v", d)
		}
	}
}

func TestParseRetryAfter(t *testing.T) {
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	cases := map[string]time.Duration{
		"":                              0,
		"3":                             3 * time.Second,
		"-1":                            0,
		"soon":                          0,
		"Sat, 24 Oct 2026 12:00:30 GMT": 30 * time.Second,
	}
	for in, want := range cases {
		if got := parseRetryAfter(in, now); got != want {
			t.Errorf("parseRetryAfter(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestCancellationStopsRetrying(t *testing.T) {
	rec := newRecorder(t, status(503, `{}`))
	c := newTestClientAt(t, rec.srv.URL)
	ctx, cancel := context.WithCancel(context.Background())
	c.sleep = func(context.Context, time.Duration) error { cancel(); return context.Canceled }
	err := c.Do(ctx, "GET", "/regions", nil, nil, nil)
	if err == nil || rec.count() != 1 {
		t.Fatalf("err %v after %d requests", err, rec.count())
	}
}
