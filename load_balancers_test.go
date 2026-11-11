package pantech

import (
	"context"
	"testing"
)

func TestLoadBalancersCreateHTTPS(t *testing.T) {
	rec := newRecorder(t, status(202, `{"operation_id":"op_1","resource_id":"lb_1","status":"submitting"}`))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.LoadBalancers.Create(context.Background(), CreateLoadBalancerRequest{
		Name: "web", ZoneID: Ptr("zone_1"), Kind: Ptr("appliance"), PublicPort: 443, PrivatePort: Ptr(8080),
		Protocol: Ptr("https"), CertificateID: Ptr("cert_1"), TLSPolicy: Ptr("tls13"), RedirectHTTP: Ptr(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `{"name":"web","zone_id":"zone_1","kind":"appliance","public_port":443,"private_port":8080,"protocol":"https","certificate_id":"cert_1","tls_policy":"tls13","redirect_http":true}`
	if rec.bodies[0] != want {
		t.Fatalf("body %s", rec.bodies[0])
	}
}

func TestLoadBalancerHTTPSFieldsAndTLSError(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"id":"lb_1","protocol":"https","certificate_id":"cert_1","certificate_fingerprint_sha256":"ab12","tls_policy":"modern","redirect_http":true}`),
		status(422, `{"title":"Invalid","code":"VALIDATION_FAILED","detail":"protocol: https is offered on appliance load balancers only.","errors":[{"field":"protocol","code":"LOAD_BALANCER_TLS_NOT_SUPPORTED","message":"https is offered on appliance load balancers only"}]}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	lb, err := c.LoadBalancers.Get(context.Background(), "lb_1")
	if err != nil || *lb.CertificateFingerprintSHA256 != "ab12" || *lb.TLSPolicy != "modern" || !lb.RedirectHTTP {
		t.Fatalf("lb %+v err %v", lb, err)
	}
	_, err = c.LoadBalancers.Create(context.Background(), CreateLoadBalancerRequest{Name: "x", PublicPort: 443, Protocol: Ptr("https")})
	if e, ok := AsError(err); !ok || !e.HasField("protocol", CodeLoadBalancerTLSNotSupported) {
		t.Fatalf("err %v", err)
	}
}
