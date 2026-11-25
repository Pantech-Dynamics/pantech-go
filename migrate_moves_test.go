package pantech

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestMigrateListProviders(t *testing.T) {
	rec := newRecorder(t, ok(`{"data":[
{"slug":"aws","name":"AWS","access":"federated","cost_readable":true,"people_readable":true,"available":true},
{"slug":"hetzner","name":"Hetzner","access":"agent","cost_readable":false,"people_readable":false,"available":true}],"next_cursor":null}`))
	c := newTestClientAt(t, rec.srv.URL)
	list, err := c.Migrate.ListProviders(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if r := rec.request(0); r.Method != "GET" || r.URL.Path != "/public/v1/cloud-providers" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	if len(list) != 2 || list[1].Slug != CloudProviderHetzner || list[1].Access != CloudAccessAgent || list[1].CostReadable {
		t.Fatalf("providers %+v", list)
	}
}

func TestMigrateAgentConnection(t *testing.T) {
	rec := newRecorder(t, cannedResponse{status: 201, body: `{"id":"ccn_1","provider":"hetzner","name":"Hetzner","status":"awaiting_agent","failure":null,
"role_arn":null,"aws_account_id":null,"external_account_id":null,"gcp":null,"azure":null,"huawei":null,"agent":null,"verified_at":null,
"created_at":"2026-11-20T10:00:00Z","updated_at":"2026-11-20T10:00:00Z",
"setup":{"access":"agent","issuer_url":null,"audience":null,"subject":"ccn_1","quick_create_url":null,"template":"Make a read-only token.",
"template_format":"agent","files":[],"permissions":["read"],"agent":{"image":"img","download_url":null,"api_url":"https://api.example.test",
"token_variables":["HCLOUD_TOKEN"],"token_how_to":"Security > API tokens","enroll_command":"agent enroll <ENROLLMENT_CODE>","run_command":"agent run","container_command":"run img"}}}`})
	c := newTestClientAt(t, rec.srv.URL)
	conn, err := c.Migrate.CreateConnection(context.Background(), CreateCloudConnectionRequest{Provider: Ptr(CloudProviderHetzner), Name: "Hetzner"})
	if err != nil {
		t.Fatal(err)
	}
	if rec.bodies[0] != `{"provider":"hetzner","name":"Hetzner"}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
	if conn.Status != CloudConnectionAwaitingAgent || conn.Setup.Access != CloudAccessAgent || conn.Setup.IssuerURL != nil ||
		conn.Setup.Agent == nil || conn.Setup.Agent.TokenVariables[0] != "HCLOUD_TOKEN" || conn.Agent != nil {
		t.Fatalf("connection %+v", conn)
	}
}

func TestMigrateUpdateConnectionSendsAFederation(t *testing.T) {
	rec := newRecorder(t, ok(`{"id":"ccn_2","provider":"gcp","status":"verifying","gcp":{"project_number":"123","pool_id":"p","provider_id":"pr",
"service_account_email":"sa@example.test","project_ids":["web"]}}`))
	c := newTestClientAt(t, rec.srv.URL)
	conn, err := c.Migrate.UpdateConnection(context.Background(), "ccn_2", UpdateCloudConnectionRequest{GCP: &CloudGCPFederation{
		ProjectNumber: "123", PoolID: "p", ProviderID: "pr", ServiceAccountEmail: "sa@example.test", ProjectIDs: []string{"web"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	want := `{"gcp":{"project_number":"123","pool_id":"p","provider_id":"pr","service_account_email":"sa@example.test","project_ids":["web"]}}`
	if rec.bodies[0] != want {
		t.Fatalf("body %s", rec.bodies[0])
	}
	if conn.GCP == nil || conn.GCP.ProjectIDs[0] != "web" {
		t.Fatalf("connection %+v", conn)
	}
}

func TestMigrateAgentEnrollmentAndRevoke(t *testing.T) {
	rec := newRecorder(t,
		cannedResponse{status: 201, body: `{"code":"ENR-1","expires_at":"2026-11-20T11:00:00Z","enroll_command":"agent enroll ENR-1"}`},
		ok(`{"id":"ccn_1","status":"awaiting_agent","agent":null}`))
	c := newTestClientAt(t, rec.srv.URL)
	enr, err := c.Migrate.CreateAgentEnrollment(context.Background(), "ccn_1")
	if err != nil || enr.Code != "ENR-1" || enr.ExpiresAt.IsZero() {
		t.Fatalf("enrollment %+v err %v", enr, err)
	}
	conn, err := c.Migrate.RevokeAgent(context.Background(), "ccn_1")
	if err != nil || conn.Status != CloudConnectionAwaitingAgent {
		t.Fatalf("connection %+v err %v", conn, err)
	}
	if r := rec.request(0); r.Method != "POST" || r.URL.Path != "/public/v1/cloud-connections/ccn_1/agent-enrollments" || r.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	if r := rec.request(1); r.Method != "DELETE" || r.URL.Path != "/public/v1/cloud-connections/ccn_1/agent" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
}

func TestMigrateSendInvites(t *testing.T) {
	rec := newRecorder(t, ok(`{"results":[{"key":"user:global:ada","email":"ada@example.test","role":"member","outcome":"invited","reason":null,"invite_id":"inv_1"},
{"key":"user:global:bob","email":null,"role":null,"outcome":"skipped","reason":"email_missing","invite_id":null}],
"migration":{"id":"mig_1","items":[{"key":"user:global:ada","kind":"user","target":{"invite_id":"inv_1","invited_at":"2026-11-20T10:00:00Z"}}]}}`))
	c := newTestClientAt(t, rec.srv.URL)
	res, err := c.Migrate.SendInvites(context.Background(), "mig_1", SendMigrationInvitesRequest{Keys: []string{"user:global:ada", "user:global:bob"}})
	if err != nil {
		t.Fatal(err)
	}
	if r := rec.request(0); r.Method != "POST" || r.URL.Path != "/public/v1/migrations/mig_1/invites" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	if rec.bodies[0] != `{"keys":["user:global:ada","user:global:bob"]}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
	if len(res.Results) != 2 || res.Results[1].Outcome != "skipped" || *res.Results[1].Reason != "email_missing" ||
		res.Migration.Items[0].Target.InvitedAt == nil {
		t.Fatalf("result %+v", res)
	}
}

func TestMigrateCopyAccessAndBucketImports(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"connection_id":"ccn_1","status":"awaiting_role","failure":null,"role_arn":null,"buckets":[{"name":"logs","region":"us-east-1",
"versioning":"off","ready":false,"failure":null}],"verified_at":null,"setup":{"issuer_url":"https://id.example.test","audience":"aud","subject":"ccn_1:copy","template":"{}"}}`),
		cannedResponse{status: 202, body: `{"id":"bim_1","connection_id":"ccn_1","migration_id":null,"source":{"provider":"aws","bucket":"logs","region":"us-east-1"},
"destination":{"bucket_id":"bkt_1","bucket_name":"logs","public_read":false},"status":"queued","failure":null,"cancel_requested":false,"pass":1,
"progress":{"objects_copied":0,"objects_skipped":0,"objects_failed":0,"bytes_copied":0,"bytes_skipped":0},"total_bytes_copied":0,
"mirror_deletes":false,"objects_deleted":0,"next_attempt_at":null,"created_at":"2026-11-20T10:00:00Z","started_at":null,"finished_at":null,
"updated_at":"2026-11-20T10:00:00Z","failures":null,"events":null}`},
		cannedResponse{status: 202, body: `{"id":"bim_1","status":"queued","pass":2,"mirror_deletes":true}`},
		cannedResponse{status: 202, body: `{"id":"bim_1","status":"running","cancel_requested":true}`},
	)
	c := newTestClientAt(t, rec.srv.URL)
	ctx := context.Background()
	access, err := c.Migrate.SetCopyAccess(ctx, "ccn_1", CopyAccessRequest{Buckets: []string{"logs"}})
	if err != nil || access.Setup.Template == nil || access.Buckets[0].Ready || *access.Buckets[0].Versioning != "off" {
		t.Fatalf("access %+v err %v", access, err)
	}
	if r := rec.request(0); r.Method != "PUT" || r.URL.Path != "/public/v1/cloud-connections/ccn_1/copy-access" || rec.bodies[0] != `{"buckets":["logs"]}` {
		t.Fatalf("sent %s %s %s", r.Method, r.URL.Path, rec.bodies[0])
	}
	imp, err := c.Migrate.StartBucketImport(ctx, StartBucketImportRequest{ConnectionID: "ccn_1", SourceBucket: "logs", DestinationBucketName: Ptr("logs")})
	if err != nil || imp.Status != BucketImportQueued || imp.Destination.BucketID != "bkt_1" || imp.Failures != nil {
		t.Fatalf("import %+v err %v", imp, err)
	}
	if rec.bodies[1] != `{"connection_id":"ccn_1","source_bucket":"logs","destination_bucket_name":"logs"}` {
		t.Fatalf("body %s", rec.bodies[1])
	}
	if imp, err = c.Migrate.ResyncBucketImport(ctx, "bim_1", ResyncBucketImportRequest{MirrorDeletes: Ptr(true)}); err != nil || imp.Pass != 2 {
		t.Fatalf("resync %+v err %v", imp, err)
	}
	if r := rec.request(2); r.URL.Path != "/public/v1/bucket-imports/bim_1/resync" || rec.bodies[2] != `{"mirror_deletes":true}` {
		t.Fatalf("sent %s %s", r.URL.Path, rec.bodies[2])
	}
	if imp, err = c.Migrate.CancelBucketImport(ctx, "bim_1"); err != nil || !imp.CancelRequested {
		t.Fatalf("cancel %+v err %v", imp, err)
	}
	if r := rec.request(3); r.Method != "POST" || r.URL.Path != "/public/v1/bucket-imports/bim_1/cancel" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
}

func TestMigrateServerCopy(t *testing.T) {
	copyJSON := `{"id":"sco_1","migration_id":"mig_1","item_key":"instance:us-east-1:i-1","connection_id":"ccn_1","source":"aws_snapshot",
"status":"%s","failure":null,"os_family":"linux","boot_mode":null,"instance":{"name":"web","plan_slug":"s-2","region":"af-abj","subnet_id":null,
"security_group_id":"sg_1","ssh_key_id":null,"id":null,"order_id":null},"disks":[{"index":0,"role":"root","source":"vol-1","size_bytes":21474836480,
"snapshot_id":"snap-1","volume_id":null}],"flags":[],"pass":{"number":1,"kind":"full","started_at":null},"pending_pass":null,
"progress":{"chunks_total":10,"chunks_done":4,"bytes_total":21474836480,"bytes_uploaded":8589934592,"wire_bytes":1024,"eta_seconds":60},
"synced_at":null,"expires_at":"2026-12-20T10:00:00Z","completed_at":null,"created_at":"2026-11-20T10:00:00Z","updated_at":"2026-11-20T10:00:00Z"}`
	rec := newRecorder(t,
		ok(`{"template":"{}","snapshot_tag":{"pantech-migrate-connection":"ccn_1"}}`),
		cannedResponse{status: 201, body: strings.Replace(copyJSON, "%s", "waiting", 1)},
		cannedResponse{status: 202, body: strings.Replace(copyJSON, "%s", "synced", 1)},
		cannedResponse{status: 202, body: strings.Replace(copyJSON, "%s", "importing", 1)},
	)
	c := newTestClientAt(t, rec.srv.URL)
	ctx := context.Background()
	tpl, err := c.Migrate.GetServerCopyTemplate(ctx, "mig_1", &GetServerCopyTemplateParams{ItemKey: "instance:us-east-1:i-1", CreateSnapshots: true})
	if err != nil || tpl.SnapshotTag["pantech-migrate-connection"] != "ccn_1" {
		t.Fatalf("template %+v err %v", tpl, err)
	}
	if q := rec.request(0).URL.Query(); q.Get("item_key") != "instance:us-east-1:i-1" || q.Get("create_snapshots") != "true" {
		t.Fatalf("query %v", q)
	}
	cp, err := c.Migrate.StartServerCopy(ctx, "mig_1", StartServerCopyRequest{
		ItemKey: "instance:us-east-1:i-1", Source: Ptr(ServerCopySourceAWSSnapshot), SecurityGroupID: Ptr("sg_1"),
		AWS: &ServerCopyAWSSource{CopyRoleARN: "arn:aws:iam::1:role/Copy", Snapshots: map[string]string{"vol-1": "snap-1"}},
	})
	if err != nil || cp.Status != ServerCopyWaiting || cp.Disks[0].SizeBytes != 21474836480 || *cp.Progress.ETASeconds != 60 {
		t.Fatalf("copy %+v err %v", cp, err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(rec.bodies[1]), &body); err != nil || body["source"] != "aws_snapshot" || body["aws"].(map[string]any)["copy_role_arn"] == nil {
		t.Fatalf("body %s", rec.bodies[1])
	}
	if cp, err = c.Migrate.SyncServerCopy(ctx, "sco_1", ServerCopyActRequest{Snapshots: map[string]string{"vol-1": "snap-2"}}); err != nil || cp.Status != ServerCopySynced {
		t.Fatalf("sync %+v err %v", cp, err)
	}
	if r := rec.request(2); r.URL.Path != "/public/v1/server-copies/sco_1/sync" || rec.bodies[2] != `{"snapshots":{"vol-1":"snap-2"}}` {
		t.Fatalf("sent %s %s", r.URL.Path, rec.bodies[2])
	}
	if cp, err = c.Migrate.CutOverServerCopy(ctx, "sco_1", ServerCopyActRequest{}); err != nil || cp.Status != ServerCopyImporting {
		t.Fatalf("cutover %+v err %v", cp, err)
	}
	if r := rec.request(3); r.URL.Path != "/public/v1/server-copies/sco_1/cutover" || rec.bodies[3] != `{}` {
		t.Fatalf("sent %s %s", r.URL.Path, rec.bodies[3])
	}
}

func TestMigrateDatabaseImport(t *testing.T) {
	rec := newRecorder(t,
		cannedResponse{status: 201, body: `{"id":"dbi_1","migration_id":"mig_1","item_key":"database:us-east-1:orders","destination_database_id":"db_1",
"target_database":"orders","mode":"replicate","source_ref":"prod-db","engine":"postgres","status":"awaiting_agent","failure":null,
"cutover_requested":false,"cancel_requested":false,"readiness_requested":false,"access_revoked":false,"agent_connected":false,
"lag_bytes":null,"position":null,"readiness":null,"destination":null,"tables":null,"charged":false,"version":1,
"created_at":"2026-11-20T10:00:00Z","updated_at":"2026-11-20T10:00:00Z","finished_at":null}`},
		ok(`{"id":"dbi_1","status":"blocked","readiness":{"checked_at":"2026-11-20T10:01:00Z","blocked":true,
"issues":[{"code":"extension_unavailable","blocking":true,"message":"postgis is not offered.","fix":null,"subject":"postgis"}],
"source":{"engine":"postgres","version":"16.4","provider":"aws_rds","database":"orders","size_bytes":1048576,"extensions":["postgis"],"table_count":3}},
"tables":[{"schema":"public","name":"orders","has_primary_key":true,"status":"verified","pages_copied":1,"rows_copied":10,"bytes_copied":8192,
"verification":{"source_rows":10,"destination_rows":10,"samples_checked":2,"mismatched_samples":0,"passed":true},"updated_at":"2026-11-20T10:02:00Z"}]}`),
		cannedResponse{status: 202, body: `{"id":"dbi_1","status":"preflight","readiness_requested":true}`},
		ok(`{"id":"dbi_1","status":"replicating","access_revoked":true}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	ctx := context.Background()
	imp, err := c.Migrate.StartDatabaseImport(ctx, "mig_1", StartDatabaseImportRequest{
		ItemKey: "database:us-east-1:orders", DestinationDatabaseID: "db_1", Mode: Ptr(DatabaseImportReplicate), SourceRef: "prod-db",
	})
	if err != nil || imp.Status != "awaiting_agent" || imp.Readiness != nil || imp.Charged {
		t.Fatalf("import %+v err %v", imp, err)
	}
	if rec.bodies[0] != `{"item_key":"database:us-east-1:orders","destination_database_id":"db_1","mode":"replicate","source_ref":"prod-db"}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
	if imp, err = c.Migrate.GetDatabaseImport(ctx, "dbi_1"); err != nil || !imp.Readiness.Blocked || *imp.Readiness.Source.SizeBytes != 1048576 ||
		!imp.Tables[0].Verification.Passed {
		t.Fatalf("import %+v err %v", imp, err)
	}
	if imp, err = c.Migrate.CheckDatabaseImportReadiness(ctx, "dbi_1"); err != nil || !imp.ReadinessRequested {
		t.Fatalf("readiness %+v err %v", imp, err)
	}
	if r := rec.request(2); r.Method != "POST" || r.URL.Path != "/public/v1/database-imports/dbi_1/readiness" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	if imp, err = c.Migrate.RevokeDatabaseImportAccess(ctx, "dbi_1"); err != nil || !imp.AccessRevoked {
		t.Fatalf("revoke %+v err %v", imp, err)
	}
	if r := rec.request(3); r.URL.Path != "/public/v1/database-imports/dbi_1/revoke-access" {
		t.Fatalf("sent %s", r.URL.Path)
	}
}

func TestMigrateKubernetesMove(t *testing.T) {
	rec := newRecorder(t,
		cannedResponse{status: 202, body: `{"id":"kmv_1","migration_id":"mig_1","item_key":"kubernetes_cluster:us-east-1:prod","connection_id":"ccn_1",
"cluster_id":"k8s_1","source_name":"prod","status":"exporting","failure":null,"options":{"include_secrets":true,"copy_volumes":false,
"publish_load_balancers":true,"namespaces":["web"]},"summary":{"objects":0,"refused":0,"skipped":0,"flagged":0,"volumes":0,"warnings":[],"addons":[]},
"allow_privileged":[],"previewed_at":null,"preview_ok":false,"applied_at":null,"completed_at":null,"sync_generation":0,
"created_at":"2026-11-20T10:00:00Z","updated_at":"2026-11-20T10:00:00Z"}`},
		ok(`{"id":"kmv_1","status":"ready","resources":[{"group":"apps","kind":"Deployment","namespace":"web","name":"api","api_version":"apps/v1",
"tier":6,"action":"create","flags":[],"policy_refused":true,"preview":{"result":"refused","code":"privileged","message":"runs privileged"},"apply":null}],
"volumes":[{"id":"v1","namespace":"web","claim":"data","size_bytes":1073741824,"source_storage_class":"gp3","target_storage_class":"standard",
"access_modes":["ReadWriteOnce"],"status":"pending","flags":[],"generation":0,"synced_generation":0,"bytes_total":0,"bytes_done":0,"files_total":0,
"checksum_verified":false,"synced_at":null,"failure":null}],"images":[],"holds":[{"group":"apps","kind":"Deployment","namespace":"web","name":"api","replicas":3,"suspend":false}]}`),
		cannedResponse{status: 202, body: `{"id":"kmv_1","status":"previewing","allow_privileged":[{"group":"apps","kind":"Deployment","namespace":"web","name":"api"}]}`},
		cannedResponse{status: 202, body: `{"id":"kmv_1","status":"cutting_over"}`},
		cannedResponse{status: 200, body: "apiVersion: v1\nkind: Namespace\n", header: map[string]string{"Content-Type": "application/yaml"}},
	)
	c := newTestClientAt(t, rec.srv.URL)
	ctx := context.Background()
	mv, err := c.Migrate.StartKubernetesMove(ctx, "mig_1", StartKubernetesMoveRequest{
		ItemKey: "kubernetes_cluster:us-east-1:prod", ClusterID: "k8s_1", IncludeSecrets: Ptr(true), Namespaces: []string{"web"},
	})
	if err != nil || mv.Status != KubernetesMoveExporting || !mv.Options.IncludeSecrets || mv.Options.Namespaces[0] != "web" {
		t.Fatalf("move %+v err %v", mv, err)
	}
	if rec.bodies[0] != `{"item_key":"kubernetes_cluster:us-east-1:prod","cluster_id":"k8s_1","include_secrets":true,"namespaces":["web"]}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
	detail, err := c.Migrate.GetKubernetesMove(ctx, "kmv_1")
	if err != nil || detail.ID != "kmv_1" || detail.Status != KubernetesMoveReady || !detail.Resources[0].PolicyRefused ||
		detail.Resources[0].Preview.Result != "refused" || detail.Volumes[0].SizeBytes != 1073741824 || detail.Holds[0].Replicas != 3 {
		t.Fatalf("detail %+v err %v", detail, err)
	}
	ref := KubernetesMoveRef{Group: "apps", Kind: "Deployment", Namespace: "web", Name: "api"}
	if mv, err = c.Migrate.PreviewKubernetesMove(ctx, "kmv_1", PreviewKubernetesMoveRequest{AllowPrivileged: []KubernetesMoveRef{ref}}); err != nil ||
		mv.AllowPrivileged[0] != ref {
		t.Fatalf("preview %+v err %v", mv, err)
	}
	if rec.bodies[2] != `{"allow_privileged":[{"group":"apps","kind":"Deployment","namespace":"web","name":"api"}]}` {
		t.Fatalf("body %s", rec.bodies[2])
	}
	if mv, err = c.Migrate.CutOverKubernetesMove(ctx, "kmv_1", CutoverKubernetesMoveRequest{SourceWritesStopped: true}); err != nil || mv.Status != KubernetesMoveCuttingOver {
		t.Fatalf("cutover %+v err %v", mv, err)
	}
	if r := rec.request(3); r.URL.Path != "/public/v1/kubernetes-moves/kmv_1/cutover" || rec.bodies[3] != `{"source_writes_stopped":true}` {
		t.Fatalf("sent %s %s", r.URL.Path, rec.bodies[3])
	}
	yaml, err := c.Migrate.DownloadKubernetesMoveManifests(ctx, "kmv_1")
	if err != nil || !strings.HasPrefix(string(yaml), "apiVersion: v1") {
		t.Fatalf("manifests %q err %v", yaml, err)
	}
	if r := rec.request(4); r.URL.Path != "/public/v1/kubernetes-moves/kmv_1/manifests" || r.Header.Get("Accept") != "application/yaml" {
		t.Fatalf("sent %s accept %q", r.URL.Path, r.Header.Get("Accept"))
	}
}

func TestMigrateKubernetesMoveActionsWithoutBodies(t *testing.T) {
	rec := newRecorder(t, cannedResponse{status: 202, body: `{"id":"kmv_1","status":"applying"}`})
	c := newTestClientAt(t, rec.srv.URL)
	ctx := context.Background()
	for i, call := range []func() (*KubernetesMove, error){
		func() (*KubernetesMove, error) { return c.Migrate.ApplyKubernetesMove(ctx, "kmv_1") },
		func() (*KubernetesMove, error) { return c.Migrate.SyncKubernetesMoveVolumes(ctx, "kmv_1") },
		func() (*KubernetesMove, error) { return c.Migrate.CancelKubernetesMove(ctx, "kmv_1") },
	} {
		if _, err := call(); err != nil {
			t.Fatal(err)
		}
		if rec.bodies[i] != "" {
			t.Fatalf("call %d sent a body %q", i, rec.bodies[i])
		}
	}
	want := []string{"apply", "sync-volumes", "cancel"}
	for i, action := range want {
		if r := rec.request(i); r.Method != "POST" || r.URL.Path != "/public/v1/kubernetes-moves/kmv_1/"+action {
			t.Fatalf("sent %s %s", r.Method, r.URL.Path)
		}
	}
}

func TestMigrateNewInventoryAndTargets(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"id":"csn_1","runner":"agent","provider":"hetzner","external_account_id":"proj-1","status":"succeeded",
"counts":{"kubernetes_clusters":1},"inventory":{"instances":[{"id":"1","os_name":"Windows Server 2022"}],
"users":[{"id":"u1","name":"ada","email":null,"is_admin":false,"source":"team","access_level":"write","roles":["Member"]}],
"kubernetes_clusters":[{"id":"c1","name":"prod","region":"fsn1","provider":"self_managed","version":"1.30","endpoint_public":true,
"node_pools":[{"name":"pool","instance_type":"cx22","vcpu":2,"memory_mib":4096,"architecture":"x86_64","gpu":false,"windows":false,
"desired":3,"autoscaling":true,"min":1,"max":5}],"addons":["ingress-nginx"],"created_at":null}]}}`),
		ok(`{"id":"mig_1","items":[{"key":"bucket:us-east-1:logs","kind":"bucket","target":{"bucket":{"name":"logs","region":"af-abj",
"versioning":false,"public_read":false,"size_bytes":5368709120,"object_count":42,"size_as_of":"2026-11-19"}}},
{"key":"kubernetes_cluster:fsn1:c1","kind":"kubernetes_cluster","target":{"kubernetes":{"version_id":"kv_1","version":"1.30",
"node_pools":[{"name":"pool","plan_slug":"k-2","count":3,"autoscaling":true,"min":1,"max":5,"source_name":"pool"}],"addons":[],"allow_downgrade":false}}},
{"key":"instance:fsn1:1","kind":"instance","target":{"image_slug":"windows-2022"}}]}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	ctx := context.Background()
	scan, err := c.Migrate.GetScan(ctx, "csn_1", &GetCloudScanParams{Include: "inventory"})
	if err != nil {
		t.Fatal(err)
	}
	inv := scan.Inventory
	if scan.Runner != "agent" || *scan.Provider != "hetzner" || scan.Counts.KubernetesClusters != 1 || *inv.Instances[0].OSName != "Windows Server 2022" ||
		*inv.Users[0].AccessLevel != "write" || inv.Users[0].Roles[0] != "Member" || inv.KubernetesClusters[0].NodePools[0].Max != 5 {
		t.Fatalf("scan %+v", scan)
	}
	mig, err := c.Migrate.UpdateItem(ctx, "mig_1", UpdateMigrationItemRequest{Key: "kubernetes_cluster:fsn1:c1", AllowVersionDowngrade: Ptr(true)})
	if err != nil {
		t.Fatal(err)
	}
	if rec.bodies[1] != `{"key":"kubernetes_cluster:fsn1:c1","allow_version_downgrade":true}` {
		t.Fatalf("body %s", rec.bodies[1])
	}
	b, k, w := mig.Items[0].Target.Bucket, mig.Items[1].Target.Kubernetes, mig.Items[2].Target.ImageSlug
	if b == nil || *b.SizeBytes != 5368709120 || *b.SizeAsOf != "2026-11-19" || k == nil || *k.NodePools[0].PlanSlug != "k-2" || w == nil || *w != "windows-2022" {
		t.Fatalf("items %+v", mig.Items)
	}
}
