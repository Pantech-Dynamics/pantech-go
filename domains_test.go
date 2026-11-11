package pantech

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

const domainJSON = `{"id":"dom_1","name":"example.com","tld":"com","status":"registering","failure_code":null,"auto_renew":true,
"transfer_lock":true,"privacy":true,"name_server_mode":"pantech","name_servers":[],"zone_id":"dnsz_1","registered_at":null,
"expires_at":null,"contact_verified":false,"settings_pending":false,"auth_code_pending":false,"suspended_reason":null,
"created_at":"2026-11-10T10:04:41Z","updated_at":"2026-11-10T10:04:41Z","version":1}`

func domainOrder(status, failure string, refunded bool) string {
	return fmt.Sprintf(`{"id":"dor_1","domain_id":"dom_1","domain_name":"example.com","kind":"register","years":2,"amount_minor":3000,
"currency":"USD","status":"%s","failure_code":%s,"refunded":%t,"created_at":"2026-11-10T10:04:41Z","completed_at":null}`, status, failure, refunded)
}

var testContact = DomainContact{
	Type: "person", FirstName: "Ada", LastName: "Obi", Email: "ada@example.com", Phone: "+234.8012345678",
	AddressLine1: "1 Marina", City: "Lagos", CountryCode: "NG",
}

func TestDomainsSearchSendsTheQuery(t *testing.T) {
	rec := newRecorder(t, ok(`{"data":[{"name":"example.com","tld":"com","status":"available",
"price":{"tld":"com","currency":"USD","register_minor":1500,"renew_minor":1500,"transfer_minor":1500,"featured":true}},
{"name":"example.zz","tld":"zz","status":"unknown","price":null}]}`))
	c := newTestClientAt(t, rec.srv.URL)
	results, err := c.Domains.Search(context.Background(), "example")
	if err != nil {
		t.Fatal(err)
	}
	r := rec.request(0)
	if r.URL.Path != "/public/v1/domains/search" || r.URL.Query().Get("query") != "example" {
		t.Fatalf("sent %s?%s", r.URL.Path, r.URL.RawQuery)
	}
	if len(results) != 2 || results[0].Status != DomainAvailable || results[0].Price.RegisterMinor != 1500 || results[1].Price != nil {
		t.Fatalf("results %+v", results)
	}
}

func TestDomainsCheckTransfer(t *testing.T) {
	rec := newRecorder(t, ok(`{"name":"example.com","transferable":false,"reason":"domain_locked","price":null}`))
	c := newTestClientAt(t, rec.srv.URL)
	check, err := c.Domains.CheckTransfer(context.Background(), "example.com")
	if err != nil {
		t.Fatal(err)
	}
	if rec.request(0).URL.Query().Get("name") != "example.com" || check.Transferable || *check.Reason != "domain_locked" {
		t.Fatalf("check %+v", check)
	}
}

func TestTLDOfferPrices(t *testing.T) {
	o := TLDOffer{RegisterMinor: 1500, RenewMinor: 1200, TransferMinor: 1000}
	if got := o.RegisterPrice(3); got != 1500+2*1200 {
		t.Fatalf("register %d", got)
	}
	if got := o.TransferPrice(1); got != 1000 {
		t.Fatalf("transfer %d", got)
	}
	if got := o.RenewPrice(2); got != 2400 {
		t.Fatalf("renew %d", got)
	}
	if got := o.RegisterPrice(0); got != 1500 {
		t.Fatalf("register for 0 years %d", got)
	}
}

func TestDomainsRegisterSendsTheOrder(t *testing.T) {
	rec := newRecorder(t, status(202, `{"domain":`+domainJSON+`,"order":`+domainOrder("pending", "null", false)+`}`))
	c := newTestClientAt(t, rec.srv.URL)
	out, err := c.Domains.Register(context.Background(), RegisterDomainRequest{
		Name: "example.com", Years: Ptr(2), Contact: testContact, Privacy: Ptr(false), AcceptTerms: true,
	}, WithIdempotencyKey("order-1"))
	if err != nil {
		t.Fatal(err)
	}
	r := rec.request(0)
	if r.Method != "POST" || r.URL.Path != "/public/v1/domains" || r.Header.Get("Idempotency-Key") != "order-1" {
		t.Fatalf("sent %s %s key %q", r.Method, r.URL.Path, r.Header.Get("Idempotency-Key"))
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(rec.bodies[0]), &body); err != nil {
		t.Fatal(err)
	}
	contact, _ := body["contact"].(map[string]any)
	if body["accept_terms"] != true || body["years"] != float64(2) || body["privacy"] != false || body["auto_renew"] != nil || contact["country_code"] != "NG" {
		t.Fatalf("body %s", rec.bodies[0])
	}
	if _, sent := contact["organization_name"]; sent {
		t.Fatalf("empty optional contact field sent: %s", rec.bodies[0])
	}
	if out.Domain.ID != "dom_1" || out.Order.ID != "dor_1" || out.Order.AmountMinor != 3000 || out.Order.Status != DomainOrderPending {
		t.Fatalf("accepted %+v", out)
	}
}

func TestDomainContactDoesNotPrint(t *testing.T) {
	req := RegisterDomainRequest{Name: "example.com", Contact: testContact}
	for _, s := range []string{fmt.Sprint(testContact), fmt.Sprintf("%+v", req), fmt.Sprintf("%#v", testContact)} {
		if strings.Contains(s, "ada@example.com") || strings.Contains(s, "Marina") {
			t.Fatalf("printed the contact: %s", s)
		}
	}
	tr := TransferDomainRequest{Name: "example.com", AuthCode: "s3cret-code", Contact: testContact}
	for _, s := range []string{fmt.Sprint(tr), fmt.Sprintf("%+v", tr), fmt.Sprintf("%#v", tr)} {
		if strings.Contains(s, "s3cret-code") || strings.Contains(s, "ada@example.com") {
			t.Fatalf("printed the transfer request: %s", s)
		}
	}
}

func TestDomainsSetNameServers(t *testing.T) {
	rec := newRecorder(t, status(202, domainJSON))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.Domains.SetNameServers(context.Background(), "dom_1", SetDomainNameServersRequest{NameServers: []string{"ns1.example.net", "ns2.example.net"}}); err != nil {
		t.Fatal(err)
	}
	r := rec.request(0)
	if r.Method != "PUT" || r.URL.Path != "/public/v1/domains/dom_1/name-servers" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	if rec.bodies[0] != `{"use_pantech_dns":false,"name_servers":["ns1.example.net","ns2.example.net"]}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
}

func TestDomainsUpdateSendsOnlyWhatIsSet(t *testing.T) {
	rec := newRecorder(t, ok(domainJSON))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.Domains.Update(context.Background(), "dom_1", UpdateDomainRequest{TransferLock: Ptr(false)}); err != nil {
		t.Fatal(err)
	}
	if rec.request(0).Method != "PATCH" || rec.bodies[0] != `{"transfer_lock":false}` {
		t.Fatalf("sent %s %s", rec.request(0).Method, rec.bodies[0])
	}
}

func TestDomainsWaitForOrder(t *testing.T) {
	rec := newRecorder(t, ok(domainOrder("paid", "null", false)), ok(domainOrder("submitted", "null", false)), ok(domainOrder("succeeded", "null", false)))
	c := newTestClientAt(t, rec.srv.URL)
	order, err := c.Domains.WaitForOrder(context.Background(), "dor_1", nil)
	if err != nil || order.Status != DomainOrderSucceeded || rec.count() != 3 {
		t.Fatalf("order %+v err %v", order, err)
	}
	if rec.request(0).URL.Path != "/public/v1/domain-orders/dor_1" {
		t.Fatalf("polled %s", rec.request(0).URL.Path)
	}
}

func TestDomainsWaitForOrderFails(t *testing.T) {
	rec := newRecorder(t, ok(domainOrder("failed", `"domain_unavailable"`, true)))
	c := newTestClientAt(t, rec.srv.URL)
	order, err := c.Domains.WaitForOrder(context.Background(), "dor_1", nil)
	var failed *DomainOrderError
	if !errors.As(err, &failed) || failed.FailureCode != "domain_unavailable" || !failed.Refunded || order == nil {
		t.Fatalf("err = %v", err)
	}
	if !strings.Contains(err.Error(), "returned to your credit") || !strings.Contains(err.Error(), "example.com") {
		t.Fatalf("message %q", err)
	}
}

func TestDomainsAuthCodeWaitsWhilePending(t *testing.T) {
	rec := newRecorder(t,
		status(202, domainJSON),
		status(202, `{"status":"pending"}`),
		ok(`{"auth_code":"Zx9-abc","expires_at":"2026-11-10T11:00:00Z"}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	ctx := context.Background()
	if _, err := c.Domains.RequestAuthCode(ctx, "dom_1"); err != nil {
		t.Fatal(err)
	}
	code, err := c.Domains.WaitForAuthCode(ctx, "dom_1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if code.Pending() || code.AuthCode != "Zx9-abc" || code.ExpiresAt == nil || rec.count() != 3 {
		t.Fatalf("code pending=%v after %d requests", code.Pending(), rec.count())
	}
	if r := rec.request(0); r.Method != "POST" || r.URL.Path != "/public/v1/domains/dom_1/auth-code" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	if r := rec.request(1); r.Method != "GET" || r.URL.Path != "/public/v1/domains/dom_1/auth-code" {
		t.Fatalf("polled %s %s", r.Method, r.URL.Path)
	}
	if strings.Contains(fmt.Sprintf("%v %+v %#v", code, *code, code), "Zx9-abc") {
		t.Fatal("printed the auth code")
	}
}

func TestDomainsLockedIsACode(t *testing.T) {
	rec := newRecorder(t, status(409, `{"title":"Domain locked","code":"DOMAIN_LOCKED","detail":"Turn the transfer lock off first."}`))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.Domains.RequestAuthCode(context.Background(), "dom_1")
	if !IsCode(err, CodeDomainLocked) || rec.count() != 1 {
		t.Fatalf("err %v after %d requests", err, rec.count())
	}
}

func TestDomainsListOrdersAndTLDs(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"data":[`+domainOrder("succeeded", "null", false)+`]}`),
		ok(`{"data":[{"tld":"com","currency":"NGN","register_minor":2500000,"renew_minor":2500000,"transfer_minor":2500000,"featured":true}]}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	orders, err := c.Domains.ListOrders(context.Background(), "dom_1")
	if err != nil || len(orders) != 1 || rec.request(0).URL.Path != "/public/v1/domains/dom_1/orders" {
		t.Fatalf("orders %+v err %v", orders, err)
	}
	tlds, err := c.Domains.ListTLDs(context.Background())
	if err != nil || len(tlds) != 1 || tlds[0].Currency != "NGN" || rec.request(1).URL.Path != "/public/v1/domains/tlds" {
		t.Fatalf("tlds %+v err %v", tlds, err)
	}
}

func awaitingPaymentOrder() string {
	return strings.Replace(domainOrder("awaiting_payment", "null", false), `"refunded":false`,
		`"refunded":false,"payment_url":"https://pay.example.test/o/dor_1","payment_expires_at":"2026-11-10T11:04:41Z"`, 1)
}

func TestDomainsWaitForOrderStopsForPayment(t *testing.T) {
	rec := newRecorder(t, ok(domainOrder("pending", "null", false)), ok(awaitingPaymentOrder()))
	c := newTestClientAt(t, rec.srv.URL)
	order, err := c.Domains.WaitForOrder(context.Background(), "dor_1", nil)
	var pay *DomainPaymentRequiredError
	if !errors.As(err, &pay) || !errors.Is(err, ErrPaymentRequired) || rec.count() != 2 {
		t.Fatalf("err %v after %d polls", err, rec.count())
	}
	if pay.PaymentURL != "https://pay.example.test/o/dor_1" || pay.PaymentExpiresAt == nil || order.Status != DomainOrderAwaitingPayment {
		t.Fatalf("payment %+v order %+v", pay, order)
	}
	if !strings.Contains(err.Error(), "https://pay.example.test/o/dor_1") {
		t.Fatalf("message %q", err)
	}
}

func TestDomainsWaitForPaymentWaitsThroughIt(t *testing.T) {
	rec := newRecorder(t, ok(awaitingPaymentOrder()), ok(awaitingPaymentOrder()), ok(domainOrder("paid", "null", false)), ok(domainOrder("succeeded", "null", false)))
	c := newTestClientAt(t, rec.srv.URL)
	order, err := c.Domains.WaitForPayment(context.Background(), "dor_1", nil)
	if err != nil || order.Status != DomainOrderSucceeded || rec.count() != 4 {
		t.Fatalf("order %+v err %v after %d polls", order, err, rec.count())
	}

	rec = newRecorder(t, ok(awaitingPaymentOrder()), ok(domainOrder("failed", `"insufficient_credit"`, false)))
	c = newTestClientAt(t, rec.srv.URL)
	_, err = c.Domains.WaitForPayment(context.Background(), "dor_1", nil)
	var failed *DomainOrderError
	if !errors.As(err, &failed) || failed.FailureCode != "insufficient_credit" {
		t.Fatalf("err %v", err)
	}
}

func TestDomainsRetryOrderPayment(t *testing.T) {
	rec := newRecorder(t, status(202, awaitingPaymentOrder()))
	c := newTestClientAt(t, rec.srv.URL)
	order, err := c.Domains.RetryOrderPayment(context.Background(), "dor_1", WithIdempotencyKey("retry-1"))
	if err != nil {
		t.Fatal(err)
	}
	r := rec.request(0)
	if r.Method != "POST" || r.URL.Path != "/public/v1/domain-orders/dor_1/retry-payment" || r.Header.Get("Idempotency-Key") != "retry-1" {
		t.Fatalf("sent %s %s key %q", r.Method, r.URL.Path, r.Header.Get("Idempotency-Key"))
	}
	if rec.bodies[0] != "" {
		t.Fatalf("body %q", rec.bodies[0])
	}
	if order.ID != "dor_1" || order.Status != DomainOrderAwaitingPayment || order.PaymentExpiresAt == nil {
		t.Fatalf("order %+v", order)
	}
}

func TestDomainsRetryOrderPaymentNotAwaiting(t *testing.T) {
	rec := newRecorder(t, status(409, `{"title":"Not awaiting payment","code":"DOMAIN_ORDER_NOT_AWAITING_PAYMENT","detail":"Place the order again."}`))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.Domains.RetryOrderPayment(context.Background(), "dor_1")
	if !IsCode(err, CodeDomainOrderNotAwaitingPayment) || rec.count() != 1 {
		t.Fatalf("err %v after %d requests", err, rec.count())
	}
}

func TestDomainPaymentRequiredWithoutLink(t *testing.T) {
	order := &DomainOrder{ID: "dor_1", Kind: DomainOrderRenew, DomainName: "example.com"}
	msg := (&DomainPaymentRequiredError{Order: order}).Error()
	if strings.Contains(msg, "confirm it at") || !strings.Contains(msg, "waits for a payment") {
		t.Fatalf("message %q", msg)
	}
}

func TestTLDOfferPrivacyAlwaysOn(t *testing.T) {
	rec := newRecorder(t, ok(`{"data":[{"tld":"com.ng","currency":"NGN","register_minor":550000,"renew_minor":550000,"transfer_minor":550000,"featured":true,"privacy_always_on":true}]}`))
	c := newTestClientAt(t, rec.srv.URL)
	tlds, err := c.Domains.ListTLDs(context.Background())
	if err != nil || !tlds[0].PrivacyAlwaysOn {
		t.Fatalf("tlds %+v err %v", tlds, err)
	}
}
