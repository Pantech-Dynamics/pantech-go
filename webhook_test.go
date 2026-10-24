package pantech

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"testing"
	"time"
)

const whSecret = "whsec_test_secret"

var whBody = []byte(`{"id":"evt_1","type":"instance.created","api_version":"2026-10-01","created_at":"2026-10-24T12:00:00Z","data":{"resource_id":"vm_1"}}`)

func TestVerifyWebhookGood(t *testing.T) {
	header := SignWebhook(whSecret, time.Now(), whBody)
	if err := VerifyWebhook(whSecret, header, whBody, 0); err != nil {
		t.Fatalf("a good signature failed: %v", err)
	}
	ev, err := ParseWebhook(whSecret, header, whBody, 0)
	if err != nil || ev.ID != "evt_1" || ev.Type != "instance.created" || ev.Data["resource_id"] != "vm_1" {
		t.Fatalf("event %+v err %v", ev, err)
	}
}

func TestVerifyWebhookMatchesTheServer(t *testing.T) {
	// Computed independently: HMAC-SHA256("whsec_test_secret", "1700000000." + body).
	now := time.Unix(1700000000, 0)
	header := "t=1700000000,v1=" + hexMAC(whSecret, "1700000000."+string(whBody))
	if err := verifyWebhookAt(whSecret, header, whBody, 0, now); err != nil {
		t.Fatal(err)
	}
	if SignWebhook(whSecret, now, whBody) != header {
		t.Fatal("SignWebhook does not produce the server's header")
	}
}

func TestVerifyWebhookBad(t *testing.T) {
	now := time.Now()
	header := SignWebhook(whSecret, now, whBody)
	cases := map[string]struct {
		secret, header string
		body           []byte
	}{
		"wrong secret":  {"whsec_other", header, whBody},
		"changed body":  {whSecret, header, []byte(`{"id":"evt_2"}`)},
		"wrong v1":      {whSecret, "t=" + strconv.FormatInt(now.Unix(), 10) + ",v1=00ff", whBody},
		"other version": {whSecret, "t=" + strconv.FormatInt(now.Unix(), 10) + ",v0=" + hexMAC(whSecret, strconv.FormatInt(now.Unix(), 10)+"."+string(whBody)), whBody},
	}
	for name, tc := range cases {
		err := VerifyWebhook(tc.secret, tc.header, tc.body, 0)
		want := ErrWebhookSignature
		if name == "other version" {
			want = ErrWebhookHeader
		}
		if !errors.Is(err, want) {
			t.Errorf("%s: err = %v, want %v", name, err, want)
		}
	}
}

func TestVerifyWebhookMalformed(t *testing.T) {
	for _, h := range []string{"", "garbage", "t=abc,v1=00", "v1=00", "t=1700000000"} {
		if err := VerifyWebhook(whSecret, h, whBody, 0); !errors.Is(err, ErrWebhookHeader) {
			t.Errorf("header %q: err = %v", h, err)
		}
	}
	if err := VerifyWebhook("", SignWebhook(whSecret, time.Now(), whBody), whBody, 0); err == nil {
		t.Error("an empty secret verified")
	}
}

func TestVerifyWebhookExpired(t *testing.T) {
	old := SignWebhook(whSecret, time.Now().Add(-6*time.Minute), whBody)
	if err := VerifyWebhook(whSecret, old, whBody, 0); !errors.Is(err, ErrWebhookTimestamp) {
		t.Fatalf("a 6-minute-old delivery: err = %v", err)
	}
	if err := VerifyWebhook(whSecret, old, whBody, 10*time.Minute); err != nil {
		t.Fatalf("within a 10-minute tolerance: %v", err)
	}
	future := SignWebhook(whSecret, time.Now().Add(6*time.Minute), whBody)
	if err := VerifyWebhook(whSecret, future, whBody, 0); !errors.Is(err, ErrWebhookTimestamp) {
		t.Fatalf("a delivery from the future: err = %v", err)
	}
}

func TestVerifyWebhookMultipleV1(t *testing.T) {
	// After a rotation the API signs with both secrets for 24 hours.
	now := time.Now()
	ts := strconv.FormatInt(now.Unix(), 10)
	header := "t=" + ts + ",v1=" + hexMAC("whsec_new", ts+"."+string(whBody)) + ", v1=" + hexMAC("whsec_old", ts+"."+string(whBody))
	for _, secret := range []string{"whsec_new", "whsec_old"} {
		if err := VerifyWebhook(secret, header, whBody, 0); err != nil {
			t.Errorf("secret %s: %v", secret, err)
		}
	}
	if err := VerifyWebhook("whsec_third", header, whBody, 0); !errors.Is(err, ErrWebhookSignature) {
		t.Errorf("an unrelated secret: err = %v", err)
	}
}

// hexMAC is the signature computed from first principles, not through the
// code under test.
func hexMAC(secret, message string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
