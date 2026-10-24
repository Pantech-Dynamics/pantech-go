package pantech

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestNewClientNeedsAKey(t *testing.T) {
	t.Setenv(EnvAPIKey, "")
	if _, err := NewClient(); err == nil || !strings.Contains(err.Error(), EnvAPIKey) {
		t.Fatalf("NewClient() error = %v, want one naming %s", err, EnvAPIKey)
	}
}

func TestNewClientReadsTheEnvironment(t *testing.T) {
	t.Setenv(EnvAPIKey, "PAN_from_env")
	t.Setenv(EnvAPIURL, "https://api.example.com/public/v1/")
	c, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	if c.BaseURL() != "https://api.example.com/public/v1" || c.apiKey() != "PAN_from_env" {
		t.Fatalf("base URL %q key %q", c.BaseURL(), c.apiKey())
	}
	c, err = NewClient(WithAPIKey("PAN_option"), WithBaseURL("https://api.example.org/public/v1"))
	if err != nil {
		t.Fatal(err)
	}
	if c.BaseURL() != "https://api.example.org/public/v1" || c.apiKey() != "PAN_option" {
		t.Fatalf("options did not win over the environment: %q %q", c.BaseURL(), c.apiKey())
	}
}

func TestNewClientDefaultsToProduction(t *testing.T) {
	t.Setenv(EnvAPIURL, "")
	c, err := NewClient(WithAPIKey(testKey))
	if err != nil {
		t.Fatal(err)
	}
	if c.BaseURL() != DefaultBaseURL {
		t.Fatalf("base URL = %q", c.BaseURL())
	}
}

func TestNewClientRejectsABadBaseURL(t *testing.T) {
	for _, u := range []string{"ftp://example.com", "not a url", "https://"} {
		if _, err := NewClient(WithAPIKey(testKey), WithBaseURL(u)); err == nil {
			t.Errorf("NewClient(WithBaseURL(%q)) succeeded", u)
		}
	}
}

func TestTheKeyIsNeverPrinted(t *testing.T) {
	c := newTestClientAt(t, "http://localhost:1")
	for _, s := range []string{fmt.Sprintf("%v", c), fmt.Sprintf("%+v", *c), fmt.Sprintf("%#v", *c)} {
		if strings.Contains(s, testKey) {
			t.Fatalf("printing the client shows the key: %s", s)
		}
	}
}

func TestRequestHeaders(t *testing.T) {
	rec := newRecorder(t, ok(`{"organization_id":"org_1","project_id":"prj_1","api_key":{"id":"key_1","name":"ci","scopes":["read"],"expires_at":null,"kind":"ai_app","client_name":"Research assistant"}}`))
	c := newTestClientAt(t, rec.srv.URL, WithHeader("Pantech-Client", "mcp/1.2.3"), WithUserAgent("pantech-mcp/1.2.3"))
	me, err := c.Me.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if me.APIKey.Kind != "ai_app" || me.APIKey.ClientName == nil || *me.APIKey.ClientName != "Research assistant" {
		t.Fatalf("me = %+v", me.APIKey)
	}
	r := rec.request(0)
	if r.URL.Path != "/public/v1/me" || r.Method != http.MethodGet {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	checks := map[string]string{
		"Authorization":  "Bearer " + testKey,
		"User-Agent":     "pantech-mcp/1.2.3",
		"Accept":         "application/json",
		"Pantech-Client": "mcp/1.2.3",
	}
	for k, want := range checks {
		if got := r.Header.Get(k); got != want {
			t.Errorf("%s = %q, want %q", k, got, want)
		}
	}
	if r.Header.Get("Idempotency-Key") != "" {
		t.Error("a GET carried an Idempotency-Key")
	}
}

func TestDefaultUserAgent(t *testing.T) {
	rec := newRecorder(t, ok(`{}`))
	c := newTestClientAt(t, rec.srv.URL)
	if err := c.Do(context.Background(), "GET", "/regions", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if got := rec.request(0).Header.Get("User-Agent"); got != "pantech-go/"+Version {
		t.Fatalf("User-Agent = %q", got)
	}
}

var uuidV4 = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestWritesCarryAnIdempotencyKeyAndJSON(t *testing.T) {
	rec := newRecorder(t, status(202, `{"order_id":"ord_1","instance_id":"vm_1","status":"awaiting_payment","amount_minor":1504000,"currency":"NGN"}`))
	c := newTestClientAt(t, rec.srv.URL)
	order, err := c.Instances.Create(context.Background(), CreateInstanceRequest{
		Name: "web-1", PlanSlug: "individual", ImageSlug: "ubuntu-24-04", SSHKeyID: "sshk_1",
		Encrypted: Ptr(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if order.OrderID != "ord_1" || order.AmountMinor != 1504000 {
		t.Fatalf("order = %+v", order)
	}
	r := rec.request(0)
	if key := r.Header.Get("Idempotency-Key"); !uuidV4.MatchString(key) {
		t.Errorf("Idempotency-Key = %q, want a UUIDv4", key)
	}
	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q", r.Header.Get("Content-Type"))
	}
	var sent map[string]any
	if err := json.Unmarshal([]byte(rec.bodies[0]), &sent); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"name": "web-1", "plan_slug": "individual", "image_slug": "ubuntu-24-04", "ssh_key_id": "sshk_1", "encrypted": true}
	if fmt.Sprint(sent) != fmt.Sprint(want) {
		t.Errorf("body = %v, want %v (optional fields left nil must not be sent)", sent, want)
	}
}

func TestWithIdempotencyKeyAndRequestHeader(t *testing.T) {
	rec := newRecorder(t, status(202, `{}`))
	c := newTestClientAt(t, rec.srv.URL)
	err := c.Do(context.Background(), "post", "instances/vm_1/start", nil, nil, nil,
		WithIdempotencyKey("my-key"), WithRequestHeader("Authorization", "Bearer PAN_other"))
	if err != nil {
		t.Fatal(err)
	}
	r := rec.request(0)
	if r.Method != "POST" || r.URL.Path != "/public/v1/instances/vm_1/start" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	if r.Header.Get("Idempotency-Key") != "my-key" || r.Header.Get("Authorization") != "Bearer PAN_other" {
		t.Fatalf("headers = %v", r.Header)
	}
}

func TestDoSendsQueryAndRawJSON(t *testing.T) {
	rec := newRecorder(t, ok(`{"x":1}`))
	c := newTestClientAt(t, rec.srv.URL)
	var out map[string]int
	err := c.Do(context.Background(), "PUT", "/things", url.Values{"a": {"b c"}}, json.RawMessage(`{"raw":true}`), &out)
	if err != nil {
		t.Fatal(err)
	}
	if out["x"] != 1 || rec.request(0).URL.RawQuery != "a=b+c" || rec.bodies[0] != `{"raw":true}` {
		t.Fatalf("out %v query %q body %q", out, rec.request(0).URL.RawQuery, rec.bodies[0])
	}
}

func TestPathParametersAreEscaped(t *testing.T) {
	rec := newRecorder(t, ok(`{}`))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.Instances.Get(context.Background(), "../me"); err != nil {
		t.Fatal(err)
	}
	if got := rec.request(0).URL.EscapedPath(); got != "/public/v1/instances/..%2Fme" {
		t.Fatalf("path = %q", got)
	}
}

func TestRawBodyIntoBytes(t *testing.T) {
	rec := newRecorder(t, cannedResponse{status: 200, body: "%PDF-1.7 ...", header: map[string]string{"Content-Type": "application/pdf"}})
	c := newTestClientAt(t, rec.srv.URL)
	pdf, err := c.Billing.GetReceiptPDF(context.Background(), "rcp_1")
	if err != nil {
		t.Fatal(err)
	}
	if string(pdf) != "%PDF-1.7 ..." || rec.request(0).Header.Get("Accept") != "application/pdf" {
		t.Fatalf("pdf %q accept %q", pdf, rec.request(0).Header.Get("Accept"))
	}
}

func TestNoContent(t *testing.T) {
	rec := newRecorder(t, status(204, ""))
	c := newTestClientAt(t, rec.srv.URL)
	if err := c.Me.Revoke(context.Background()); err != nil {
		t.Fatal(err)
	}
	if r := rec.request(0); r.Method != "DELETE" || r.URL.Path != "/public/v1/me" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
}

func TestResponseSizeIsCapped(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chunk := []byte(strings.Repeat("a", 1<<20))
		for i := 0; i <= maxResponseBytes>>20; i++ {
			if _, err := w.Write(chunk); err != nil {
				return
			}
		}
	}))
	defer srv.Close()
	c := newTestClientAt(t, srv.URL, WithMaxRetries(0))
	err := c.Do(context.Background(), "GET", "/big", nil, nil, nil)
	if !errors.Is(err, errTooLarge) {
		t.Fatalf("err = %v, want errTooLarge", err)
	}
}

func TestContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()
	c := newTestClientAt(t, srv.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := c.Instances.Get(ctx, "vm_1")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
	if time.Since(start) > 5*time.Second {
		t.Fatal("the cancelled call did not return promptly")
	}
}

func TestOptionalPointersDecode(t *testing.T) {
	rec := newRecorder(t, ok(`{"id":"vm_1","name":"web","plan_id":"plan_1","image_id":"img_1","desired_state":"running","observed_state":"running","tags":{},"public_ipv4":null,"private_ipv4":"10.0.0.5","failure":{"code":"INSTANCE_BUSY","reason":"Busy."},"created_at":"2026-10-03T10:00:00.123456789Z"}`))
	c := newTestClientAt(t, rec.srv.URL)
	vm, err := c.Instances.Get(context.Background(), "vm_1")
	if err != nil {
		t.Fatal(err)
	}
	if vm.PublicIPv4 != nil || vm.PrivateIPv4 == nil || *vm.PrivateIPv4 != "10.0.0.5" || vm.Failure.Code != "INSTANCE_BUSY" || vm.CreatedAt.Year() != 2026 {
		t.Fatalf("instance = %+v", vm)
	}
}

func TestTypedMethodsTakeRequestOptions(t *testing.T) {
	rec := newRecorder(t, status(503, `{}`), status(202, `{}`))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.Instances.Rename(context.Background(), "vm_1", RenameRequest{Name: "web-2"}, WithIdempotencyKey("rename-web-2"), WithNoRetry())
	if !IsCode(err, "") || rec.count() != 1 {
		t.Fatalf("WithNoRetry ignored: err %v, %d requests", err, rec.count())
	}
	if got := rec.request(0).Header.Get("Idempotency-Key"); got != "rename-web-2" {
		t.Fatalf("Idempotency-Key = %q", got)
	}
}
