package pantech

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

// Webhook delivery headers.
const (
	WebhookSignatureHeader  = "Pantech-Signature"
	WebhookEventIDHeader    = "Pantech-Event-Id"
	WebhookEventTypeHeader  = "Pantech-Event-Type"
	WebhookDeliveryIDHeader = "Pantech-Delivery-Id"
)

// DefaultWebhookTolerance is how old a delivery's timestamp may be.
const DefaultWebhookTolerance = 5 * time.Minute

// Errors VerifyWebhook returns. Answer any of them with a 4xx and do not act
// on the body.
var (
	// ErrWebhookHeader means the Pantech-Signature header is missing or malformed.
	ErrWebhookHeader = errors.New("pantech: the webhook signature header is missing or malformed")
	// ErrWebhookTimestamp means the signature's timestamp is outside the tolerance.
	ErrWebhookTimestamp = errors.New("pantech: the webhook signature has expired")
	// ErrWebhookSignature means no v1 signature matches the body and secret.
	ErrWebhookSignature = errors.New("pantech: the webhook signature does not match")
)

// VerifyWebhook checks a delivery's Pantech-Signature header,
// "t=<unix seconds>,v1=<hex>[,v1=<hex>]", against the raw request body: one
// v1 must equal HMAC-SHA256(secret, t + "." + body), compared in constant time,
// and t must be within tolerance of now (DefaultWebhookTolerance when
// tolerance is 0). Pass the body exactly as received, before any JSON parsing.
// After a secret rotation there are two v1 values for 24 hours, one per secret.
func VerifyWebhook(secret string, header string, body []byte, tolerance time.Duration) error {
	return verifyWebhookAt(secret, header, body, tolerance, time.Now())
}

func verifyWebhookAt(secret, header string, body []byte, tolerance time.Duration, now time.Time) error {
	if tolerance <= 0 {
		tolerance = DefaultWebhookTolerance
	}
	if secret == "" {
		return errors.New("pantech: no webhook secret to verify with")
	}
	var ts string
	var sigs [][]byte
	for _, part := range strings.Split(header, ",") {
		k, v, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok {
			continue
		}
		switch k {
		case "t":
			ts = v
		case "v1":
			if sig, err := hex.DecodeString(v); err == nil {
				sigs = append(sigs, sig)
			}
		}
	}
	secs, err := strconv.ParseInt(ts, 10, 64)
	if err != nil || len(sigs) == 0 {
		return ErrWebhookHeader
	}
	age := now.Sub(time.Unix(secs, 0))
	if age > tolerance || age < -tolerance {
		return ErrWebhookTimestamp
	}
	want := webhookMAC(secret, ts, body)
	for _, sig := range sigs {
		if hmac.Equal(sig, want) {
			return nil
		}
	}
	return ErrWebhookSignature
}

// SignWebhook returns the Pantech-Signature header value for body signed with
// secret at t, as the API sends it. Use it to test your receiver.
func SignWebhook(secret string, t time.Time, body []byte) string {
	ts := strconv.FormatInt(t.Unix(), 10)
	return "t=" + ts + ",v1=" + hex.EncodeToString(webhookMAC(secret, ts, body))
}

// ParseWebhook verifies a delivery like VerifyWebhook, then decodes its body.
func ParseWebhook(secret string, header string, body []byte, tolerance time.Duration) (*WebhookEvent, error) {
	if err := VerifyWebhook(secret, header, body, tolerance); err != nil {
		return nil, err
	}
	var ev WebhookEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil, errors.New("pantech: the webhook body is not a JSON event")
	}
	return &ev, nil
}

func webhookMAC(secret, ts string, body []byte) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts))
	mac.Write([]byte("."))
	mac.Write(body)
	return mac.Sum(nil)
}
