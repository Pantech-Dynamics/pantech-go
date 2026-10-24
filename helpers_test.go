package pantech

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

const testKey = "PAN_test_0123456789"

// newTestClientAt returns a client for a test server whose waits and retry
// backoffs return at once.
func newTestClientAt(t *testing.T, serverURL string, opts ...Option) *Client {
	t.Helper()
	opts = append([]Option{WithAPIKey(testKey), WithBaseURL(serverURL + "/public/v1")}, opts...)
	c, err := NewClient(opts...)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	c.sleep = func(ctx context.Context, d time.Duration) error { return ctx.Err() }
	c.jitter = func(d time.Duration) time.Duration { return d }
	return c
}

// recorder is a test server that answers each request with the next canned
// response and records what it was sent.
type recorder struct {
	mu        sync.Mutex
	requests  []*http.Request
	bodies    []string
	responses []cannedResponse
	srv       *httptest.Server
}

type cannedResponse struct {
	status int
	body   string
	header map[string]string
}

func newRecorder(t *testing.T, responses ...cannedResponse) *recorder {
	t.Helper()
	r := &recorder{responses: responses}
	r.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.mu.Lock()
		defer r.mu.Unlock()
		buf, _ := io.ReadAll(req.Body)
		r.requests = append(r.requests, req)
		r.bodies = append(r.bodies, string(buf))
		resp := cannedResponse{status: 200, body: `{}`}
		if n := len(r.requests) - 1; n < len(r.responses) {
			resp = r.responses[n]
		} else if len(r.responses) > 0 {
			resp = r.responses[len(r.responses)-1]
		}
		for k, v := range resp.header {
			w.Header().Set(k, v)
		}
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(resp.status)
		_, _ = w.Write([]byte(resp.body))
	}))
	t.Cleanup(r.srv.Close)
	return r
}

func (r *recorder) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.requests)
}

func (r *recorder) request(i int) *http.Request {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.requests[i]
}

func ok(body string) cannedResponse { return cannedResponse{status: 200, body: body} }

func status(code int, body string) cannedResponse { return cannedResponse{status: code, body: body} }
