package pantech

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

const certificateJSON = `{"id":"cert_1","kind":"public","names":["example.com","api.example.org"],"common_name":"example.com","description":"",
"key_type":"ecdsa-p256","exportable":true,"status":"%s","pending_reason":%s,"failure_code":%s,
"validations":[{"name":"example.com","method":"pantech_dns","zone_id":"dnsz_1","cname_name":null,"cname_target":null,"status":"ready"},
{"name":"api.example.org","method":"cname","zone_id":null,"cname_name":"_acme-challenge.api.example.org","cname_target":"k3x9.acme.pantechdynamics.com","status":"pending"}],
"validation_deadline_at":"2026-11-13T10:00:00Z","issuer":null,"serial_number":null,"fingerprint_sha256":null,"not_before":null,"not_after":null,
"issued_at":null,"renew_after":null,"renewal_status":null,"in_use_by":[],"revocation_pending":false,"revoked_at":null,
"created_at":"2026-11-10T10:00:00Z","updated_at":"2026-11-10T10:00:00Z","version":1}`

func cert(status, pending, failure string) string {
	s := strings.Replace(certificateJSON, "%s", status, 1)
	s = strings.Replace(s, "%s", pending, 1)
	return strings.Replace(s, "%s", failure, 1)
}

const testKeyPEM = "-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgSECRETKEYMATERIAL\n-----END PRIVATE KEY-----\n"

func TestCertificatesRequest(t *testing.T) {
	rec := newRecorder(t, status(202, cert("pending_validation", `"awaiting_cname"`, "null")))
	c := newTestClientAt(t, rec.srv.URL)
	got, err := c.Certificates.Request(context.Background(), RequestCertificateRequest{Names: []string{"example.com", "api.example.org"}, Exportable: Ptr(true)})
	if err != nil {
		t.Fatal(err)
	}
	r := rec.request(0)
	if r.Method != "POST" || r.URL.Path != "/public/v1/certificates" || r.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	if rec.bodies[0] != `{"names":["example.com","api.example.org"],"exportable":true}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
	v := got.Validations[1]
	if got.Status != CertificatePendingValidation || *got.PendingReason != "awaiting_cname" || v.Method != ValidationCNAME || *v.CNAMEName != "_acme-challenge.api.example.org" {
		t.Fatalf("certificate %+v", got)
	}
}

func TestCertificatesListFilters(t *testing.T) {
	rec := newRecorder(t, ok(`{"data":[`+cert("issued", "null", "null")+`],"next_cursor":null}`))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.Certificates.List(context.Background(), &ListCertificatesParams{Q: "example", Status: CertificateIssued, Kind: CertificatePublic}); err != nil {
		t.Fatal(err)
	}
	q := rec.request(0).URL.Query()
	if q.Get("q") != "example" || q.Get("status") != "issued" || q.Get("kind") != "public" {
		t.Fatalf("query %s", rec.request(0).URL.RawQuery)
	}
}

func TestCertificatesWaitForCertificate(t *testing.T) {
	rec := newRecorder(t,
		ok(cert("pending_validation", `"awaiting_cname"`, "null")),
		ok(cert("issuing", `"awaiting_ca"`, "null")),
		ok(cert("issued", "null", "null")),
	)
	c := newTestClientAt(t, rec.srv.URL)
	got, err := c.Certificates.WaitForCertificate(context.Background(), "cert_1", nil)
	if err != nil || got.Status != CertificateIssued || rec.count() != 3 || rec.request(0).URL.Path != "/public/v1/certificates/cert_1" {
		t.Fatalf("certificate %+v err %v", got, err)
	}
}

func TestCertificatesWaitForCertificateFails(t *testing.T) {
	rec := newRecorder(t, ok(cert("failed", "null", `"validation_timeout"`)))
	c := newTestClientAt(t, rec.srv.URL)
	got, err := c.Certificates.WaitForCertificate(context.Background(), "cert_1", nil)
	var failed *CertificateFailureError
	if !errors.As(err, &failed) || failed.FailureCode != "validation_timeout" || got == nil || !strings.Contains(err.Error(), "cert_1") {
		t.Fatalf("err %v", err)
	}
}

func TestCertificatesExportKeyNeverPrints(t *testing.T) {
	rec := newRecorder(t, ok(`{"certificate_pem":"-----BEGIN CERTIFICATE-----\nLEAF\n-----END CERTIFICATE-----","chain_pem":"-----BEGIN CERTIFICATE-----\nINTER\n-----END CERTIFICATE-----\n",
"private_key_pem":"-----BEGIN PRIVATE KEY-----\nSECRETKEYMATERIAL\n-----END PRIVATE KEY-----\n","not_after":"2027-02-08T10:00:00Z"}`))
	c := newTestClientAt(t, rec.srv.URL)
	exp, err := c.Certificates.Export(context.Background(), "cert_1")
	if err != nil {
		t.Fatal(err)
	}
	if r := rec.request(0); r.Method != "POST" || r.URL.Path != "/public/v1/certificates/cert_1/export" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	if !strings.Contains(exp.PrivateKeyPEM, "SECRETKEYMATERIAL") {
		t.Fatal("no key")
	}
	if fc := exp.FullChainPEM(); !strings.HasPrefix(fc, "-----BEGIN CERTIFICATE-----\nLEAF\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nINTER") {
		t.Fatalf("full chain %q", fc)
	}
	holder := struct{ E *ExportedCertificate }{exp}
	imp := ImportCertificateRequest{CertificatePEM: "x", PrivateKeyPEM: testKeyPEM}
	re := ReimportCertificateRequest{CertificatePEM: "x", PrivateKeyPEM: testKeyPEM}
	for _, s := range []string{fmt.Sprint(exp), fmt.Sprintf("%v %+v %#v", *exp, *exp, exp), fmt.Sprintf("%+v", holder),
		fmt.Sprintf("%v %+v %#v", imp, imp, imp), fmt.Sprintf("%v %+v %#v", re, re, &re)} {
		if strings.Contains(s, "SECRETKEYMATERIAL") {
			t.Fatalf("printed the key: %s", s)
		}
	}
}

func TestCertificatesImportSendsTheKey(t *testing.T) {
	rec := newRecorder(t, status(201, cert("issued", "null", "null")))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.Certificates.Import(context.Background(), ImportCertificateRequest{CertificatePEM: "LEAF", PrivateKeyPEM: testKeyPEM}); err != nil {
		t.Fatal(err)
	}
	if r := rec.request(0); r.URL.Path != "/public/v1/certificates/import" || !strings.Contains(rec.bodies[0], `"private_key_pem"`) || strings.Contains(rec.bodies[0], "chain_pem") {
		t.Fatalf("sent %s %s", r.URL.Path, rec.bodies[0])
	}
}

func TestCertificatesRevokeAndInUse(t *testing.T) {
	rec := newRecorder(t,
		status(202, strings.Replace(cert("issued", "null", "null"), `"revocation_pending":false`, `"revocation_pending":true`, 1)),
		status(409, `{"title":"In use","code":"CERTIFICATE_IN_USE","detail":"A load balancer uses this certificate."}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	got, err := c.Certificates.Revoke(context.Background(), "cert_1", RevokeCertificateRequest{Reason: Ptr("key_compromise")})
	if err != nil || !got.RevocationPending || rec.bodies[0] != `{"reason":"key_compromise"}` {
		t.Fatalf("certificate %+v err %v body %s", got, err, rec.bodies[0])
	}
	if err := c.Certificates.Delete(context.Background(), "cert_1"); !IsCode(err, CodeCertificateInUse) {
		t.Fatalf("err %v", err)
	}
}
