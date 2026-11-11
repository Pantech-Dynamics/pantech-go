package pantech

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

const zoneJSON = `{"id":"dnsz_1","name":"example.com","description":"","status":"%s","failure_code":%s,
"name_servers":["ns-1.example-dns.org","ns-2.example-dns.com","ns-3.example-dns.net","ns-4.example-dns.co.uk"],
"domain_id":null,"record_count":2,"last_checked_at":null,"activated_at":null,"suspended_reason":null,
"created_at":"2026-11-10T09:58:12Z","updated_at":"2026-11-10T10:00:00Z","version":1}`

func zone(status, failure string) string {
	return strings.Replace(strings.Replace(zoneJSON, "%s", status, 1), "%s", failure, 1)
}

const recordJSON = `{"id":"dnsr_1","zone_id":"dnsz_1","name":"www","fqdn":"www.example.com","type":"A","ttl":300,
"values":["203.0.113.10"],"allow_private":false,"comment":"","managed":false,"status":"%s","failure_code":%s,"version":2,
"created_at":"2026-11-10T10:01:00Z","updated_at":"2026-11-10T10:01:00Z"}`

func record(status, failure string) string {
	return strings.Replace(strings.Replace(recordJSON, "%s", status, 1), "%s", failure, 1)
}

func TestDNSCreateZoneSendsTheBody(t *testing.T) {
	rec := newRecorder(t, status(201, zone("pending", "null")))
	c := newTestClientAt(t, rec.srv.URL)
	z, err := c.DNS.CreateZone(context.Background(), CreateDNSZoneRequest{Name: "example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if z.ID != "dnsz_1" || z.Status != DNSZonePending || len(z.NameServers) != 4 || z.DomainID != nil {
		t.Fatalf("zone %+v", z)
	}
	r := rec.request(0)
	if r.Method != "POST" || r.URL.Path != "/public/v1/dns/zones" || r.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("sent %s %s key %q", r.Method, r.URL.Path, r.Header.Get("Idempotency-Key"))
	}
	if rec.bodies[0] != `{"name":"example.com"}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
}

func TestDNSUpdateZoneCanClearTheDescription(t *testing.T) {
	rec := newRecorder(t, ok(zone("active", "null")))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.DNS.UpdateZone(context.Background(), "dnsz_1", UpdateDNSZoneRequest{}); err != nil {
		t.Fatal(err)
	}
	if r := rec.request(0); r.Method != "PATCH" || r.URL.Path != "/public/v1/dns/zones/dnsz_1" || rec.bodies[0] != `{"description":""}` {
		t.Fatalf("sent %s %s %s", r.Method, r.URL.Path, rec.bodies[0])
	}
}

func TestDNSListRecordsSendsTheFilters(t *testing.T) {
	rec := newRecorder(t, ok(`{"data":[`+record("active", "null")+`],"next_cursor":null}`))
	c := newTestClientAt(t, rec.srv.URL)
	page, err := c.DNS.ListRecords(context.Background(), "dnsz_1", &ListDNSRecordsParams{Name: "www", Type: "A", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	q := rec.request(0).URL.Query()
	if q.Get("name") != "www" || q.Get("type") != "A" || q.Get("limit") != "10" || q.Has("q") || q.Has("cursor") {
		t.Fatalf("query %s", rec.request(0).URL.RawQuery)
	}
	if len(page.Data) != 1 || page.Data[0].FQDN != "www.example.com" || page.HasMore() {
		t.Fatalf("page %+v", page)
	}
}

func TestDNSAllZonesFollowsTheCursor(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"data":[`+zone("active", "null")+`],"next_cursor":"c2"}`),
		ok(`{"data":[`+zone("awaiting_delegation", "null")+`],"next_cursor":null}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	var statuses []string
	for z, err := range c.DNS.AllZones(context.Background(), &ListDNSZonesParams{Status: "active"}) {
		if err != nil {
			t.Fatal(err)
		}
		statuses = append(statuses, z.Status)
	}
	if len(statuses) != 2 || rec.request(1).URL.Query().Get("cursor") != "c2" || rec.request(1).URL.Query().Get("status") != "active" {
		t.Fatalf("statuses %v, second query %s", statuses, rec.request(1).URL.RawQuery)
	}
}

func TestDNSUpdateRecordSendsIfMatchAndOnlyWhatIsSet(t *testing.T) {
	rec := newRecorder(t, status(202, record("pending", "null")))
	c := newTestClientAt(t, rec.srv.URL)
	r, err := c.DNS.UpdateRecord(context.Background(), "dnsz_1", "dnsr_1", UpdateDNSRecordRequest{TTL: Ptr(60)}, WithIfMatch(2))
	if err != nil {
		t.Fatal(err)
	}
	if r.Status != DNSRecordPending {
		t.Fatalf("record %+v", r)
	}
	req := rec.request(0)
	if req.Method != "PATCH" || req.URL.Path != "/public/v1/dns/zones/dnsz_1/records/dnsr_1" || req.Header.Get("If-Match") != "2" {
		t.Fatalf("sent %s %s If-Match %q", req.Method, req.URL.Path, req.Header.Get("If-Match"))
	}
	if rec.bodies[0] != `{"ttl":60}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
}

func TestDNSVersionConflictIsACode(t *testing.T) {
	rec := newRecorder(t, cannedResponse{status: 412, body: `{"type":"about:blank","title":"Version conflict","code":"VERSION_CONFLICT","detail":"The record changed since you read it."}`,
		header: map[string]string{"Content-Type": "application/problem+json"}})
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.DNS.DeleteRecord(context.Background(), "dnsz_1", "dnsr_1", WithIfMatch(1))
	if !IsCode(err, CodeVersionConflict) || rec.count() != 1 {
		t.Fatalf("err %v after %d requests", err, rec.count())
	}
	if rec.request(0).Header.Get("If-Match") != "1" {
		t.Fatalf("If-Match %q", rec.request(0).Header.Get("If-Match"))
	}
}

func TestDNSBatchRecords(t *testing.T) {
	rec := newRecorder(t, status(202, `{"data":[`+record("pending", "null")+`,`+record("deleting", "null")+`]}`))
	c := newTestClientAt(t, rec.srv.URL)
	out, err := c.DNS.BatchRecords(context.Background(), "dnsz_1", DNSRecordBatchRequest{Changes: []DNSRecordChange{
		{Action: DNSChangeCreate, Record: &DNSRecordInput{Name: "www", Type: DNSRecordA, Values: []string{"203.0.113.10"}}},
		{Action: DNSChangeDelete, RecordID: "dnsr_2", Version: Ptr(3)},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 || out[1].Status != DNSRecordDeleting {
		t.Fatalf("records %+v", out)
	}
	if r := rec.request(0); r.Method != "POST" || r.URL.Path != "/public/v1/dns/zones/dnsz_1/records/batch" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	want := `{"changes":[{"action":"create","record":{"name":"www","type":"A","values":["203.0.113.10"]}},{"action":"delete","record_id":"dnsr_2","version":3}]}`
	if rec.bodies[0] != want {
		t.Fatalf("body %s", rec.bodies[0])
	}
}

func TestDNSExportZoneReturnsTheZoneFile(t *testing.T) {
	file := "$ORIGIN example.com.\nwww 300 IN A 203.0.113.10\n"
	rec := newRecorder(t, cannedResponse{status: 200, body: file, header: map[string]string{"Content-Type": "text/dns"}})
	c := newTestClientAt(t, rec.srv.URL)
	out, err := c.DNS.ExportZone(context.Background(), "dnsz_1")
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != file || rec.request(0).Header.Get("Accept") != "text/dns" || rec.request(0).URL.Path != "/public/v1/dns/zones/dnsz_1/export" {
		t.Fatalf("got %q, Accept %q", out, rec.request(0).Header.Get("Accept"))
	}
}

func TestDNSImportZone(t *testing.T) {
	rec := newRecorder(t, status(202, `{"data":[]}`))
	c := newTestClientAt(t, rec.srv.URL)
	out, err := c.DNS.ImportZone(context.Background(), "dnsz_1", ImportDNSZoneRequest{ZoneFile: "www 300 IN A 203.0.113.10\n", Replace: Ptr(true)})
	if err != nil {
		t.Fatal(err)
	}
	if out == nil || len(out) != 0 {
		t.Fatalf("records %v", out)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(rec.bodies[0]), &body); err != nil || body["replace"] != true || body["zone_file"] != "www 300 IN A 203.0.113.10\n" {
		t.Fatalf("body %s", rec.bodies[0])
	}
}

func TestDNSInvalidZoneFileNamesTheField(t *testing.T) {
	rec := newRecorder(t, status(422, `{"title":"Invalid zone file","code":"INVALID_ZONE_FILE","detail":"Line 3 does not parse.","errors":[{"field":"zone_file","code":"invalid","message":"Line 3 does not parse."}]}`))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.DNS.ImportZone(context.Background(), "dnsz_1", ImportDNSZoneRequest{ZoneFile: "garbage"})
	e, ok := AsError(err)
	if !ok || e.Code != CodeInvalidZoneFile || !e.HasField("zone_file", "invalid") {
		t.Fatalf("err %v", err)
	}
}

func TestDNSWaitForRecord(t *testing.T) {
	rec := newRecorder(t, ok(record("pending", "null")), ok(record("active", "null")))
	c := newTestClientAt(t, rec.srv.URL)
	r, err := c.DNS.WaitForRecord(context.Background(), "dnsz_1", "dnsr_1", nil)
	if err != nil || r.Status != DNSRecordActive || rec.count() != 2 {
		t.Fatalf("record %+v err %v after %d polls", r, err, rec.count())
	}
	if rec.request(0).URL.Path != "/public/v1/dns/zones/dnsz_1/records/dnsr_1" {
		t.Fatalf("polled %s", rec.request(0).URL.Path)
	}
}

func TestDNSWaitForRecordFails(t *testing.T) {
	rec := newRecorder(t, ok(record("failed", `"invalid_change"`)))
	c := newTestClientAt(t, rec.srv.URL)
	r, err := c.DNS.WaitForRecord(context.Background(), "dnsz_1", "dnsr_1", nil)
	var failed *DNSFailureError
	if !errors.As(err, &failed) || failed.Kind != "record" || failed.FailureCode != "invalid_change" || r == nil {
		t.Fatalf("err = %v", err)
	}
	if !strings.Contains(err.Error(), "dnsr_1") || !strings.Contains(err.Error(), "invalid_change") {
		t.Fatalf("message %q", err)
	}
}

func TestDNSWaitForZone(t *testing.T) {
	rec := newRecorder(t, ok(zone("pending", "null")), ok(zone("awaiting_delegation", "null")), ok(zone("active", "null")))
	c := newTestClientAt(t, rec.srv.URL)
	z, err := c.DNS.WaitForZone(context.Background(), "dnsz_1", nil)
	if err != nil || z.Status != DNSZoneActive || rec.count() != 3 {
		t.Fatalf("zone %+v err %v", z, err)
	}
}

func TestDNSWaitForZoneFails(t *testing.T) {
	rec := newRecorder(t, ok(zone("failed", `"delegation_timeout"`)))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.DNS.WaitForZone(context.Background(), "dnsz_1", nil)
	var failed *DNSFailureError
	if !errors.As(err, &failed) || failed.Kind != "zone" || failed.FailureCode != "delegation_timeout" {
		t.Fatalf("err = %v", err)
	}
}

func TestDNSZoneNotFoundIsNotRetried(t *testing.T) {
	rec := newRecorder(t, status(404, `{"title":"Not found","code":"DNS_ZONE_NOT_FOUND","detail":"No such zone."}`))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.DNS.GetZone(context.Background(), "dnsz_404")
	if !IsCode(err, CodeDNSZoneNotFound) || rec.count() != 1 {
		t.Fatalf("err %v after %d requests", err, rec.count())
	}
}

func TestDNSPathsAreEscaped(t *testing.T) {
	rec := newRecorder(t, ok(record("active", "null")))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.DNS.GetRecord(context.Background(), "dnsz_1/../x", "dnsr_1"); err != nil {
		t.Fatal(err)
	}
	if got := rec.request(0).URL.EscapedPath(); got != "/public/v1/dns/zones/dnsz_1%2F..%2Fx/records/dnsr_1" {
		t.Fatalf("path %s", got)
	}
}

func TestDNSQuoteZones(t *testing.T) {
	rec := newRecorder(t, ok(`{"free_zones":5,"live_zones":7,"billable_zones":2,"currency":"USD","monthly_per_zone_minor":50,"monthly_estimate_minor":100,"nominal_month_hours":730}`))
	c := newTestClientAt(t, rec.srv.URL)
	q, err := c.DNS.QuoteZones(context.Background())
	if err != nil || rec.request(0).URL.Path != "/public/v1/dns/zones/quote" {
		t.Fatalf("err %v", err)
	}
	if q.BillableZones != 2 || q.MonthlyEstimateMinor != 100 || q.FreeZones != 5 {
		t.Fatalf("quote %+v", q)
	}
}
