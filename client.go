// Package pantech is the official Go SDK for the Pantech Dynamics public API
// (https://docs.pantechdynamics.com/api).
//
// Create a client with NewClient, then call the typed services on it:
//
//	c, err := pantech.NewClient() // reads PANTECH_API_KEY
//	if err != nil { ... }
//	me, err := c.Me.Get(ctx)
//
// Every write carries an Idempotency-Key that is reused when the request is
// retried, so a retry never makes a change twice. Errors from the API are
// *Error values carrying the problem's stable code; branch on it with IsCode.
package pantech

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Version is this SDK's version. It goes into the default User-Agent.
const Version = "0.1.0"

// DefaultBaseURL is the production public API. Every path is relative to it.
const DefaultBaseURL = "https://api.pantechdynamics.com/public/v1"

// DefaultMaxRetries is how many times a request is retried after the first
// try, unless WithMaxRetries says otherwise.
const DefaultMaxRetries = 3

// Environment variables NewClient reads when no option sets the value.
const (
	EnvAPIKey = "PANTECH_API_KEY"
	EnvAPIURL = "PANTECH_API_URL"
)

const (
	// defaultTimeout bounds each attempt. Some API calls take tens of seconds.
	defaultTimeout = 90 * time.Second

	// maxResponseBytes bounds how much of a response is read. Receipts as PDF
	// are the largest bodies the API returns.
	maxResponseBytes = 32 << 20

	defaultPollInterval = 2 * time.Second
)

// Client talks to the Pantech Dynamics public API. It is safe for concurrent
// use. Create it with NewClient.
type Client struct {
	Me                  *MeService
	Billing             *BillingService
	Catalog             *CatalogService
	Instances           *InstancesService
	Networks            *NetworksService
	Subnets             *SubnetsService
	FirewallRules       *FirewallRulesService
	PublicIPs           *PublicIPsService
	PortForwardingRules *PortForwardingRulesService
	SecurityGroups      *SecurityGroupsService
	LoadBalancers       *LoadBalancersService
	Volumes             *VolumesService
	Snapshots           *SnapshotsService
	Databases           *DatabasesService
	Kubernetes          *KubernetesService
	SSHKeys             *SSHKeysService
	Webhooks            *WebhooksService
	Operations          *OperationsService

	baseURL    string
	userAgent  string
	maxRetries int
	headers    http.Header
	http       *http.Client

	// apiKey is a function, not a string, so printing the Client with %+v
	// shows a function address instead of the key.
	apiKey func() string

	// sleep and jitter are fields so tests can run retries and waits
	// without waiting.
	sleep        func(ctx context.Context, d time.Duration) error
	jitter       func(d time.Duration) time.Duration
	pollInterval time.Duration
}

// Option configures a Client in NewClient.
type Option func(*settings)

type settings struct {
	apiKey     string
	baseURL    string
	userAgent  string
	maxRetries int
	headers    http.Header
	http       *http.Client
}

// WithAPIKey sets the API key (PAN_…). Without it NewClient reads
// PANTECH_API_KEY.
func WithAPIKey(key string) Option {
	return func(s *settings) { s.apiKey = key }
}

// WithBaseURL sets the API's base URL, including the /public/v1 prefix, for
// example https://api.example.com/public/v1. Without it NewClient
// reads PANTECH_API_URL, and then uses DefaultBaseURL.
func WithBaseURL(u string) Option {
	return func(s *settings) { s.baseURL = u }
}

// WithHTTPClient sets the HTTP client requests go through. Its Timeout bounds
// each attempt; the default client allows 90 seconds.
func WithHTTPClient(h *http.Client) Option {
	return func(s *settings) { s.http = h }
}

// WithUserAgent replaces the User-Agent header, which is pantech-go/<Version>
// by default. Tools built on the SDK set their own, such as pantech-mcp/1.0.0.
func WithUserAgent(ua string) Option {
	return func(s *settings) { s.userAgent = ua }
}

// WithMaxRetries sets how many times a failed request is retried after the
// first try. 0 turns retries off. The default is DefaultMaxRetries.
func WithMaxRetries(n int) Option {
	return func(s *settings) { s.maxRetries = max(n, 0) }
}

// WithHeader adds a header to every request, for example Pantech-Client. It
// is set after the SDK's own headers, so it can replace them.
func WithHeader(key, value string) Option {
	return func(s *settings) { s.headers.Set(key, value) }
}

// NewClient returns a client configured by opts and the environment. It fails
// when there is no API key or the base URL is not an http(s) URL.
func NewClient(opts ...Option) (*Client, error) {
	s := settings{
		apiKey:     os.Getenv(EnvAPIKey),
		baseURL:    os.Getenv(EnvAPIURL),
		userAgent:  "pantech-go/" + Version,
		maxRetries: DefaultMaxRetries,
		headers:    http.Header{},
	}
	for _, opt := range opts {
		opt(&s)
	}
	if s.baseURL == "" {
		s.baseURL = DefaultBaseURL
	}
	u, err := url.Parse(s.baseURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, fmt.Errorf("pantech: invalid base URL %q: use an http(s) URL such as %s", s.baseURL, DefaultBaseURL)
	}
	key := strings.TrimSpace(s.apiKey)
	if key == "" {
		return nil, errors.New("pantech: no API key: pass WithAPIKey or set " + EnvAPIKey)
	}
	if s.http == nil {
		s.http = &http.Client{Timeout: defaultTimeout}
	}

	c := &Client{
		baseURL:      strings.TrimRight(s.baseURL, "/"),
		userAgent:    s.userAgent,
		maxRetries:   s.maxRetries,
		headers:      s.headers,
		http:         s.http,
		apiKey:       func() string { return key },
		sleep:        sleepContext,
		jitter:       halfJitter,
		pollInterval: defaultPollInterval,
	}
	c.Me = &MeService{c}
	c.Billing = &BillingService{c}
	c.Catalog = &CatalogService{c}
	c.Instances = &InstancesService{c}
	c.Networks = &NetworksService{c}
	c.Subnets = &SubnetsService{c}
	c.FirewallRules = &FirewallRulesService{c}
	c.PublicIPs = &PublicIPsService{c}
	c.PortForwardingRules = &PortForwardingRulesService{c}
	c.SecurityGroups = &SecurityGroupsService{c}
	c.LoadBalancers = &LoadBalancersService{c}
	c.Volumes = &VolumesService{c}
	c.Snapshots = &SnapshotsService{c}
	c.Databases = &DatabasesService{c}
	c.Kubernetes = &KubernetesService{c}
	c.SSHKeys = &SSHKeysService{c}
	c.Webhooks = &WebhooksService{c}
	c.Operations = &OperationsService{c}
	return c, nil
}

// BaseURL is the base URL the client sends requests to.
func (c *Client) BaseURL() string { return c.baseURL }

// RequestOption changes one call made with Do.
type RequestOption func(*requestOptions)

type requestOptions struct {
	idempotencyKey string
	noRetry        bool
	headers        http.Header
}

// WithIdempotencyKey sends key as the Idempotency-Key instead of a new random
// one. Use it to make a change that is safe to repeat across processes: the API
// returns the original result for a key it has seen.
func WithIdempotencyKey(key string) RequestOption {
	return func(o *requestOptions) { o.idempotencyKey = key }
}

// WithRequestHeader sets a header on this call only. It is applied last, so it
// replaces any header the SDK or WithHeader set.
func WithRequestHeader(key, value string) RequestOption {
	return func(o *requestOptions) {
		if o.headers == nil {
			o.headers = http.Header{}
		}
		o.headers.Set(key, value)
	}
}

// WithNoRetry sends the call once, whatever the answer.
func WithNoRetry() RequestOption {
	return func(o *requestOptions) { o.noRetry = true }
}

// Do sends one API call and decodes the answer. It is what every typed method
// uses, and it is exported for calls the typed services do not cover.
//
// path is relative to the base URL, such as "/instances" or
// "/instances/vm_123/start"; query may be nil. body is nil or a value that
// marshals to JSON (a json.RawMessage is sent as is). out is nil or a pointer:
// a 2xx JSON body is decoded into it, and a *[]byte receives the raw body.
//
// POST, PUT, PATCH and DELETE carry an Idempotency-Key, the same one on every
// retry. Network errors, 429 (after Retry-After), 500, 502, 503 and 504 are
// retried with jittered backoff, except on the routes the API does not replay
// (creating or deleting an SSH key, opening an instance console, and testing,
// deleting, rotating the secret of or redelivering a webhook), which are
// retried only on 429. A non-2xx answer is returned as an *Error.
func (c *Client) Do(ctx context.Context, method, path string, query url.Values, body, out any, opts ...RequestOption) error {
	var ro requestOptions
	for _, opt := range opts {
		opt(&ro)
	}
	method = strings.ToUpper(method)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	payload, err := encodeBody(body)
	if err != nil {
		return err
	}
	key := ro.idempotencyKey
	if key == "" && isWrite(method) {
		if key, err = newIdempotencyKey(); err != nil {
			return err
		}
	}

	target := c.baseURL + path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}
	replays := replaysOnRetry(method, path)

	for attempt := 0; ; attempt++ {
		res, err := c.send(ctx, method, target, payload, key, ro.headers)
		if err != nil && ctx.Err() != nil {
			// Cancelled or out of time: say so, not "connection failed".
			return fmt.Errorf("pantech: %s %s: %w", method, path, ctx.Err())
		}
		if err == nil && res.status >= 200 && res.status < 300 {
			return decodeBody(res.body, out)
		}
		var status int
		var retryAfter time.Duration
		if err == nil {
			status = res.status
			retryAfter = parseRetryAfter(res.header.Get("Retry-After"), time.Now())
		}
		if !ro.noRetry && attempt < c.maxRetries && shouldRetry(replays, status, err) {
			if err := c.sleep(ctx, c.backoff(attempt, retryAfter)); err != nil {
				return fmt.Errorf("pantech: %s %s: %w", method, path, err)
			}
			continue
		}
		if err != nil {
			return fmt.Errorf("pantech: %s %s: %w", method, path, err)
		}
		return newError(res.status, res.header, res.body)
	}
}

type response struct {
	status int
	header http.Header
	body   []byte
}

// send makes one attempt and reads the whole body, so the connection is free
// before the retry loop decides what to do.
func (c *Client) send(ctx context.Context, method, target string, payload []byte, idempotencyKey string, extra http.Header) (*response, error) {
	req, err := c.newRequest(ctx, method, target, payload, idempotencyKey, extra)
	if err != nil {
		return nil, err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := readCapped(res.Body)
	if err != nil {
		return nil, err
	}
	return &response{status: res.StatusCode, header: res.Header, body: data}, nil
}

// newRequest builds a fresh request per attempt, since a body is read once.
func (c *Client) newRequest(ctx context.Context, method, target string, payload []byte, idempotencyKey string, extra http.Header) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, target, body)
	if err != nil {
		return nil, fmt.Errorf("pantech: building the request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey())
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}
	for k, v := range c.headers {
		req.Header[k] = v
	}
	for k, v := range extra {
		req.Header[k] = v
	}
	return req, nil
}

var errTooLarge = fmt.Errorf("pantech: the response is larger than %d bytes", maxResponseBytes)

func readCapped(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxResponseBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxResponseBytes {
		return nil, errTooLarge
	}
	return data, nil
}

func encodeBody(body any) ([]byte, error) {
	switch b := body.(type) {
	case nil:
		return nil, nil
	case json.RawMessage:
		return b, nil
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("pantech: encoding the request body: %w", err)
	}
	return payload, nil
}

func decodeBody(body []byte, out any) error {
	if out == nil {
		return nil
	}
	if raw, ok := out.(*[]byte); ok {
		*raw = body
		return nil
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("pantech: decoding the response: %w", err)
	}
	return nil
}

func isWrite(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}
	return false
}
