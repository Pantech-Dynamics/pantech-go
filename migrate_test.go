package pantech

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

const connectionJSON = `{"id":"cc_1","provider":"aws","name":"Production","status":"%s","failure":%s,
"role_arn":null,"aws_account_id":null,"verified_at":null,"created_at":"2026-10-27T10:00:00Z","updated_at":"2026-10-27T10:00:00Z",
"setup":{"issuer_url":"https://id.example.test","audience":"aud","subject":"cc_1","quick_create_url":"https://console.example.test/stack","template":"{}"}}`

func connection(status, failure string) string {
	return strings.Replace(strings.Replace(connectionJSON, "%s", status, 1), "%s", failure, 1)
}

func TestMigrateCreateConnectionSendsTheBody(t *testing.T) {
	rec := newRecorder(t, cannedResponse{status: 201, body: connection("awaiting_role", "null")})
	c := newTestClientAt(t, rec.srv.URL)
	conn, err := c.Migrate.CreateConnection(context.Background(), CreateCloudConnectionRequest{Name: "Production"})
	if err != nil {
		t.Fatal(err)
	}
	if conn.Status != CloudConnectionAwaitingRole || conn.Setup.QuickCreateURL == nil || conn.RoleARN != nil {
		t.Fatalf("connection %+v", conn)
	}
	r := rec.request(0)
	if r.Method != "POST" || r.URL.Path != "/public/v1/cloud-connections" || r.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("sent %s %s key %q", r.Method, r.URL.Path, r.Header.Get("Idempotency-Key"))
	}
	if rec.bodies[0] != `{"name":"Production"}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
}

func TestMigrateUpdateConnectionSendsOnlyWhatIsSet(t *testing.T) {
	rec := newRecorder(t, ok(connection("verifying", "null")))
	c := newTestClientAt(t, rec.srv.URL)
	arn := "arn:aws:iam::123456789012:role/PantechMigrateReadOnly"
	if _, err := c.Migrate.UpdateConnection(context.Background(), "cc_1", UpdateCloudConnectionRequest{RoleARN: &arn}); err != nil {
		t.Fatal(err)
	}
	if r := rec.request(0); r.Method != "PATCH" || r.URL.Path != "/public/v1/cloud-connections/cc_1" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	if rec.bodies[0] != `{"role_arn":"`+arn+`"}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
}

func TestMigrateGetScanIncludesTheInventory(t *testing.T) {
	rec := newRecorder(t, ok(`{"id":"cs_1","connection_id":"cc_1","status":"succeeded","failure":null,"created_at":"2026-10-27T10:00:00Z",
"started_at":null,"finished_at":null,"aws_account_id":"123456789012","regions":["us-east-1"],
"counts":{"instances":1,"volumes":0,"networks":0,"subnets":0,"security_groups":0,"public_ips":0,"load_balancers":0,"databases":0,"buckets":0,"users":0,"ssh_keys":0},
"cost":{"currency":"USD","average_monthly_minor":123456,"months":[{"month":"2026-09","amount_minor":123456}]},"warnings":[],"truncated":[],
"inventory":{"instances":[{"id":"i-1","name":"web","region":"us-east-1","zone":null,"instance_type":"t3.small","state":"running","vcpu":2,"memory_mib":2048,
"architecture":"x86_64","platform":"linux","volume_ids":[],"network_id":null,"subnet_id":null,"security_group_ids":[],"public_ip":null,"private_ip":null,"key_name":null,"launched_at":null}],
"volumes":[],"networks":[],"subnets":[],"security_groups":[],"public_ips":[],"load_balancers":[],"databases":[],"buckets":[],"users":[],"ssh_keys":[]}}`))
	c := newTestClientAt(t, rec.srv.URL)
	scan, err := c.Migrate.GetScan(context.Background(), "cs_1", &GetCloudScanParams{Include: "inventory"})
	if err != nil {
		t.Fatal(err)
	}
	if got := rec.request(0).URL.Query().Get("include"); got != "inventory" {
		t.Fatalf("include = %q", got)
	}
	if scan.Inventory == nil || len(scan.Inventory.Instances) != 1 || *scan.Inventory.Instances[0].VCPU != 2 {
		t.Fatalf("inventory %+v", scan.Inventory)
	}
	if scan.Cost == nil || scan.Cost.AverageMonthlyMinor != 123456 || scan.Counts.Instances != 1 {
		t.Fatalf("scan %+v", scan)
	}
}

func TestMigrateGetScanWithoutParamsSendsNoQuery(t *testing.T) {
	rec := newRecorder(t, ok(`{"id":"cs_1","status":"running","inventory":null,"counts":null,"cost":null}`))
	c := newTestClientAt(t, rec.srv.URL)
	scan, err := c.Migrate.GetScan(context.Background(), "cs_1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if rec.request(0).URL.RawQuery != "" || scan.Inventory != nil || scan.Counts != nil {
		t.Fatalf("query %q scan %+v", rec.request(0).URL.RawQuery, scan)
	}
}

func TestMigrateUpdateItemAndEstimate(t *testing.T) {
	rec := newRecorder(t, ok(`{"id":"mig_1","name":"Move","region":"af-abj","connection_id":"cc_1","scan_id":"cs_1","stage":"planning","version":2,
"summary":{"planned":1,"excluded":1,"unsupported":0,"by_kind":{"instance":1}},
"estimate":{"currency":"NGN","pantech_priced":true,"pantech_monthly_minor":5000000,"covers":["servers"],"not_covered":[],"source_known":true,
"source_monthly_minor":123456,"source_currency":"USD","source_months":3,"source_scope":"whole_account"},
"items":[{"key":"instance:us-east-1:i-1","kind":"instance","source":{"id":"i-1","name":"web","region":"us-east-1","summary":"t3.small","vcpu":2,"memory_mib":2048,"disk_gb":20},
"status":"planned","flags":[],"target":{"plan_slug":"s-2vcpu-4gb","disk_gb":20,"cidr":null,"network_key":null,"subnet_key":null,"instance_key":null,"rules":[],
"engine":null,"engine_version":null,"storage_gb":null,"role":null,"email":null,"public_key":null,"listeners":[]}}],
"created_at":"2026-10-27T10:00:00Z","updated_at":"2026-10-27T10:05:00Z"}`))
	c := newTestClientAt(t, rec.srv.URL)
	mig, err := c.Migrate.UpdateItem(context.Background(), "mig_1", UpdateMigrationItemRequest{Key: "volume:us-east-1:vol-1", Excluded: Ptr(true)})
	if err != nil {
		t.Fatal(err)
	}
	r := rec.request(0)
	if r.Method != "PATCH" || r.URL.Path != "/public/v1/migrations/mig_1/items" || r.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(rec.bodies[0]), &body); err != nil || len(body) != 2 || body["excluded"] != true {
		t.Fatalf("body %s", rec.bodies[0])
	}
	if mig.Estimate == nil || *mig.Estimate.PantechMonthlyMinor != 5000000 || mig.Summary.ByKind["instance"] != 1 {
		t.Fatalf("migration %+v", mig)
	}
	if len(mig.Items) != 1 || *mig.Items[0].Target.PlanSlug != "s-2vcpu-4gb" {
		t.Fatalf("items %+v", mig.Items)
	}
}

func TestMigrateListMigrationsHasNoItems(t *testing.T) {
	rec := newRecorder(t, ok(`{"data":[{"id":"mig_1","name":"Move","stage":"planning","estimate":null,"items":null,"summary":{"planned":0,"excluded":0,"unsupported":0,"by_kind":{}}}],"next_cursor":null}`))
	c := newTestClientAt(t, rec.srv.URL)
	list, err := c.Migrate.ListMigrations(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Items != nil || list[0].Estimate != nil {
		t.Fatalf("list %+v", list)
	}
}

func TestMigrateExportTerraform(t *testing.T) {
	rec := newRecorder(t, ok(`{"files":[{"name":"main.tf","content":"terraform {}\n"},{"name":"variables.tf","content":""}]}`))
	c := newTestClientAt(t, rec.srv.URL)
	tf, err := c.Migrate.ExportTerraform(context.Background(), "mig_1")
	if err != nil {
		t.Fatal(err)
	}
	if rec.request(0).URL.Path != "/public/v1/migrations/mig_1/terraform" || len(tf.Files) != 2 || tf.Files[0].Name != "main.tf" {
		t.Fatalf("files %+v", tf.Files)
	}
}

func TestMigrateWaitForConnection(t *testing.T) {
	rec := newRecorder(t,
		ok(connection("verifying", "null")),
		ok(connection("connected", "null")),
	)
	c := newTestClientAt(t, rec.srv.URL)
	conn, err := c.Migrate.WaitForConnection(context.Background(), "cc_1", nil)
	if err != nil || conn.Status != CloudConnectionConnected || rec.count() != 2 {
		t.Fatalf("conn %+v err %v after %d polls", conn, err, rec.count())
	}
}

func TestMigrateWaitForConnectionFails(t *testing.T) {
	rec := newRecorder(t, ok(connection("failed", `{"code":"assume_role_denied","message":"Check the role's trust policy."}`)))
	c := newTestClientAt(t, rec.srv.URL)
	conn, err := c.Migrate.WaitForConnection(context.Background(), "cc_1", nil)
	var failed *CloudFailureError
	if !errors.As(err, &failed) || failed.Failure.Code != "assume_role_denied" || conn == nil {
		t.Fatalf("err = %v", err)
	}
	if !strings.Contains(err.Error(), "assume_role_denied") || !strings.Contains(err.Error(), "cc_1") {
		t.Fatalf("message = %q", err)
	}
}

func TestMigrateWaitForScan(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"id":"cs_1","status":"queued"}`),
		ok(`{"id":"cs_1","status":"running"}`),
		ok(`{"id":"cs_1","status":"succeeded"}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	scan, err := c.Migrate.WaitForScan(context.Background(), "cs_1", nil)
	if err != nil || scan.Status != CloudScanSucceeded || rec.count() != 3 {
		t.Fatalf("scan %+v err %v", scan, err)
	}
	if rec.request(0).URL.Path != "/public/v1/cloud-scans/cs_1" {
		t.Fatalf("polled %s", rec.request(0).URL.Path)
	}
}

func TestMigrateWaitForScanFails(t *testing.T) {
	rec := newRecorder(t, ok(`{"id":"cs_2","status":"failed","failure":{"code":"scan_too_large","message":"Too much to read."}}`))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.Migrate.WaitForScan(context.Background(), "cs_2", nil)
	var failed *CloudFailureError
	if !errors.As(err, &failed) || failed.Kind != "scan" || failed.Failure.Code != "scan_too_large" {
		t.Fatalf("err = %v", err)
	}
}
