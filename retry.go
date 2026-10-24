package pantech

import (
	"context"
	"crypto/rand"
	"fmt"
	mathrand "math/rand/v2"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	baseDelay     = 500 * time.Millisecond
	maxDelay      = 8 * time.Second
	maxRetryAfter = 60 * time.Second
)

// noReplay lists the routes the API documents as not replaying a request with
// the same Idempotency-Key: a retry after a lost answer would act again (or
// lose a generated private key or a new secret), so they are not retried.
var noReplay = []struct {
	method string
	path   *regexp.Regexp
}{
	{http.MethodPost, regexp.MustCompile(`^/ssh-keys$`)},
	{http.MethodDelete, regexp.MustCompile(`^/ssh-keys/[^/]+$`)},
	{http.MethodPost, regexp.MustCompile(`^/instances/[^/]+/console$`)},
	{http.MethodPost, regexp.MustCompile(`^/webhooks/[^/]+/test$`)},
	{http.MethodDelete, regexp.MustCompile(`^/webhooks/[^/]+$`)},
	{http.MethodPost, regexp.MustCompile(`^/webhooks/[^/]+/rotate-secret$`)},
	{http.MethodPost, regexp.MustCompile(`^/webhooks/[^/]+/deliveries/[^/]+/redeliver$`)},
}

// replaysOnRetry reports whether a retried request is answered from the
// API's idempotency record rather than acted on again.
func replaysOnRetry(method, path string) bool {
	path = strings.TrimRight(path, "/")
	for _, r := range noReplay {
		if r.method == method && r.path.MatchString(path) {
			return false
		}
	}
	return true
}

// shouldRetry is the retry policy. A 429 was refused before anything ran, so
// it is retried on every route. Network errors and 500, 502, 503 and 504 are
// retried where the API replays, which is every GET and every write except
// the noReplay routes.
func shouldRetry(replays bool, status int, err error) bool {
	if err == nil && status == http.StatusTooManyRequests {
		return true
	}
	if !replays {
		return false
	}
	if err != nil {
		return true
	}
	switch status {
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}
	return false
}

// backoff is the server's Retry-After when it sent one (capped), else
// exponential backoff from 500ms up to 8s with jitter.
func (c *Client) backoff(attempt int, retryAfter time.Duration) time.Duration {
	if retryAfter > 0 {
		return min(retryAfter, maxRetryAfter)
	}
	d := baseDelay << attempt
	if d <= 0 || d > maxDelay {
		d = maxDelay
	}
	return c.jitter(d)
}

// parseRetryAfter reads Retry-After in seconds or as an HTTP date.
func parseRetryAfter(v string, now time.Time) time.Duration {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(v); err == nil {
		if secs <= 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil && t.After(now) {
		return t.Sub(now)
	}
	return 0
}

// sleepContext waits for d, or returns early with ctx's error.
func sleepContext(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// halfJitter is a random duration in [d/2, d], so clients throttled together
// do not retry together.
func halfJitter(d time.Duration) time.Duration {
	half := d / 2
	if half <= 0 {
		return d
	}
	return half + mathrand.N(half+1)
}

// newIdempotencyKey returns a random UUID (version 4).
func newIdempotencyKey() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("pantech: generating an Idempotency-Key: %w", err)
	}
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
