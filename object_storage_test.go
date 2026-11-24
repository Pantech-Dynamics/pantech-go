package pantech

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestObjectStorageStatusUsageAndKeys(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"enabled":true,"endpoint":"https://s3.af-abj.pantechdynamics.com","beta":true,"direct_upload_enabled":true,"limits":{"max_buckets":25,"max_storage_bytes":107374182400,"max_objects":1000000}}`),
		ok(`{"range":"7d","total_bytes_sent":4096,"total_bytes_received":8192,"total_read_ops":12,"total_write_ops":5,"hours":[{"hour":"2026-11-24T10:00:00Z","bytes_sent":4096,"bytes_received":8192,"read_ops":12,"write_ops":5}]}`),
		ok(`{"data":[{"id":"osk_1","name":"backups","region":"af-abj","access_key_id":"AKIA1","read_only":true,"buckets":["assets"],"created_at":"2026-11-24T10:00:00Z"}]}`),
		status(201, `{"id":"osk_2","name":"ci","region":"af-abj","access_key_id":"AKIA2","read_only":false,"buckets":["assets"],"created_at":"2026-11-24T11:00:00Z","secret_access_key":"super-secret-value","endpoint":"https://s3.af-abj.pantechdynamics.com"}`),
		status(204, ``),
	)
	c := newTestClientAt(t, rec.srv.URL)
	ctx := context.Background()

	st, err := c.ObjectStorage.GetStatus(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !st.Enabled || !st.DirectUploadEnabled || st.Limits == nil || st.Limits.MaxBuckets != 25 {
		t.Fatalf("unexpected status: %+v", st)
	}

	usage, err := c.ObjectStorage.GetUsage(ctx, &ObjectStorageUsageParams{Range: "7d"})
	if err != nil {
		t.Fatal(err)
	}
	if rec.request(1).URL.RawQuery != "range=7d" {
		t.Fatalf("query = %q", rec.request(1).URL.RawQuery)
	}
	if usage.TotalBytesSent != 4096 || usage.TotalWriteOps != 5 || len(usage.Hours) != 1 {
		t.Fatalf("unexpected usage: %+v", usage)
	}

	keys, err := c.ObjectStorage.ListKeys(ctx, &ListObjectStorageKeysParams{Region: "af-abj"})
	if err != nil {
		t.Fatal(err)
	}
	if rec.request(2).URL.RawQuery != "region=af-abj" || len(keys) != 1 || keys[0].ID != "osk_1" {
		t.Fatalf("keys = %+v query = %q", keys, rec.request(2).URL.RawQuery)
	}

	created, err := c.ObjectStorage.CreateKey(ctx, CreateObjectStorageKeyRequest{
		Name:    "ci",
		Region:  Ptr("af-abj"),
		Buckets: Ptr([]string{"assets"}),
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.SecretAccessKey != "super-secret-value" || created.Endpoint != "https://s3.af-abj.pantechdynamics.com" {
		t.Fatalf("unexpected created key")
	}
	if strings.Contains(fmt.Sprint(created), "super-secret-value") || strings.Contains(fmt.Sprintf("%#v", created), "super-secret-value") {
		t.Fatalf("CreatedObjectStorageKey leaked secret in fmt output: %s", fmt.Sprint(created))
	}

	if err := c.ObjectStorage.RevokeKey(ctx, "osk_2"); err != nil {
		t.Fatal(err)
	}
	if req := rec.request(4); req.Method != "DELETE" || req.URL.Path != "/public/v1/object-storage/keys/osk_2" {
		t.Fatalf("revoke sent %s %s", req.Method, req.URL.Path)
	}
}

func TestBucketsCRUDStatsAndEmpty(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"data":[{"id":"bkt_1","name":"assets","region":"af-abj","endpoint":"https://s3.af-abj.pantechdynamics.com","quota_bytes":0,"size_bytes":1024,"object_count":2,"versioning":true,"object_lock":false,"public_prefixes":["public/"],"state":"active","created_at":"2026-11-24T10:00:00Z","retired_at":null,"purge_due_at":null}]}`),
		status(201, `{"id":"bkt_2","name":"backups","region":"af-abj","endpoint":"https://s3.af-abj.pantechdynamics.com","quota_bytes":1073741824,"size_bytes":0,"object_count":0,"versioning":true,"object_lock":true,"public_prefixes":[],"state":"provisioning","created_at":"2026-11-24T10:05:00Z","retired_at":null,"purge_due_at":null,"operation_id":"op_1"}`),
		ok(`{"id":"bkt_2","name":"backups","region":"af-abj","endpoint":"https://s3.af-abj.pantechdynamics.com","quota_bytes":1073741824,"size_bytes":0,"object_count":0,"versioning":true,"object_lock":true,"public_prefixes":[],"state":"active","created_at":"2026-11-24T10:05:00Z","retired_at":null,"purge_due_at":null}`),
		ok(`{"id":"bkt_2","name":"backups","region":"af-abj","endpoint":"https://s3.af-abj.pantechdynamics.com","quota_bytes":2147483648,"size_bytes":0,"object_count":0,"versioning":true,"object_lock":true,"public_prefixes":["pub/"],"state":"updating","created_at":"2026-11-24T10:05:00Z","retired_at":null,"purge_due_at":null,"operation_id":"op_2"}`),
		ok(`{"size_bytes":4096,"object_count":7}`),
		ok(`{"deleted_count":5,"aborted_uploads":1,"locked_count":2,"sample_locked_keys":["legal/a.pdf"]}`),
		status(202, `{"operation_id":"op_3","resource_id":"bkt_2","status":"submitting"}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	ctx := context.Background()

	list, err := c.Buckets.List(ctx)
	if err != nil || len(list) != 1 || list[0].ID != "bkt_1" {
		t.Fatalf("list = %+v err = %v", list, err)
	}

	created, err := c.Buckets.Create(ctx, CreateBucketRequest{
		Name:       "backups",
		QuotaBytes: Ptr(int64(1073741824)),
		Versioning: Ptr(true),
		ObjectLock: Ptr(true),
	})
	if err != nil || created.ID != "bkt_2" || created.OperationID != "op_1" {
		t.Fatalf("created = %+v err = %v", created, err)
	}

	got, err := c.Buckets.Get(ctx, "bkt_2")
	if err != nil || got.State != "active" {
		t.Fatalf("got = %+v err = %v", got, err)
	}

	updated, err := c.Buckets.Update(ctx, "bkt_2", UpdateBucketRequest{
		QuotaBytes:     Ptr(int64(2147483648)),
		PublicPrefixes: Ptr([]string{"pub/"}),
	})
	if err != nil || updated.OperationID == nil || *updated.OperationID != "op_2" {
		t.Fatalf("updated = %+v err = %v", updated, err)
	}

	stats, err := c.Buckets.GetStats(ctx, "bkt_2")
	if err != nil || stats.SizeBytes != 4096 || stats.ObjectCount != 7 {
		t.Fatalf("stats = %+v err = %v", stats, err)
	}

	emptied, err := c.Buckets.Empty(ctx, "bkt_2")
	if err != nil || emptied.DeletedCount != 5 || emptied.AbortedUploads != 1 || emptied.LockedCount != 2 || len(emptied.SampleLockedKeys) != 1 {
		t.Fatalf("emptied = %+v err = %v", emptied, err)
	}

	acc, err := c.Buckets.Delete(ctx, "bkt_2", &DeleteBucketParams{Force: true})
	if err != nil || acc.OperationID != "op_3" {
		t.Fatalf("delete = %+v err = %v", acc, err)
	}
	if rec.request(6).URL.RawQuery != "force=true" {
		t.Fatalf("delete query = %q", rec.request(6).URL.RawQuery)
	}
}

func TestBucketRulesCORSAndLifecycleAndTagsAndLock(t *testing.T) {
	corsJSON := `{"rules":[{"id":"web","allowed_origins":["https://app.example.com"],"allowed_methods":["GET","PUT"],"allowed_headers":["*"],"expose_headers":["ETag"],"max_age_seconds":3600}]}`
	lcJSON := `{"rules":[{"id":"expire-tmp","prefix":"tmp/","enabled":true,"expiration_days":30,"noncurrent_expiration_days":7,"abort_incomplete_multipart_days":3,"expired_object_delete_marker":false}]}`
	tagsJSON := `{"tags":[{"key":"env","value":"prod"}]}`
	lockJSON := `{"enabled":true,"mode":"GOVERNANCE","days":30,"years":0}`

	rec := newRecorder(t,
		ok(corsJSON), ok(corsJSON),
		ok(lcJSON), ok(lcJSON),
		ok(tagsJSON), ok(tagsJSON),
		ok(lockJSON), ok(lockJSON),
	)
	c := newTestClientAt(t, rec.srv.URL)
	ctx := context.Background()

	cors, err := c.Buckets.GetCORS(ctx, "bkt_1")
	if err != nil || len(cors.Rules) != 1 || cors.Rules[0].MaxAgeSeconds != 3600 {
		t.Fatalf("get cors = %+v err = %v", cors, err)
	}
	cors, err = c.Buckets.PutCORS(ctx, "bkt_1", *cors)
	if err != nil || len(cors.Rules) != 1 {
		t.Fatalf("put cors = %+v err = %v", cors, err)
	}

	lc, err := c.Buckets.GetLifecycle(ctx, "bkt_1")
	if err != nil || len(lc.Rules) != 1 || lc.Rules[0].ExpirationDays != 30 {
		t.Fatalf("get lifecycle = %+v err = %v", lc, err)
	}
	lc, err = c.Buckets.PutLifecycle(ctx, "bkt_1", *lc)
	if err != nil || len(lc.Rules) != 1 {
		t.Fatalf("put lifecycle = %+v err = %v", lc, err)
	}

	tags, err := c.Buckets.GetTags(ctx, "bkt_1")
	if err != nil || len(tags.Tags) != 1 || tags.Tags[0].Value != "prod" {
		t.Fatalf("get tags = %+v err = %v", tags, err)
	}
	tags, err = c.Buckets.PutTags(ctx, "bkt_1", *tags)
	if err != nil || len(tags.Tags) != 1 {
		t.Fatalf("put tags = %+v err = %v", tags, err)
	}

	lock, err := c.Buckets.GetObjectLock(ctx, "bkt_1")
	if err != nil || !lock.Enabled || lock.Mode != "GOVERNANCE" || lock.Days != 30 {
		t.Fatalf("get lock = %+v err = %v", lock, err)
	}
	lock, err = c.Buckets.PutObjectLockConfig(ctx, "bkt_1", *lock)
	if err != nil || !lock.Enabled || lock.Days != 30 {
		t.Fatalf("put lock = %+v err = %v", lock, err)
	}
}

func TestBucketObjectsOperations(t *testing.T) {
	rawContent := []byte{0x89, 0x50, 0x4e, 0x47, 0x00, 0xff}
	rec := newRecorder(t,
		ok(`{"prefix":"img/","folders":["img/thumbs/"],"objects":[{"key":"img/logo.png","size_bytes":6,"etag":"e1","content_type":"image/png","last_modified":"2026-11-24T12:00:00Z","storage_class":"STANDARD","version_id":"v1","is_latest":true,"public":true,"object_public":true,"public_source":"object","public_url":"https://api.pantechdynamics.com/storage/public/bkt_1/img/logo.png"}],"next_cursor":"c2","is_truncated":true}`),
		ok(`{"versions":[{"key":"img/logo.png","version_id":"v1","is_latest":true,"is_delete_marker":false,"size_bytes":6,"etag":"e1","last_modified":"2026-11-24T12:00:00Z","storage_class":"STANDARD"}],"next_key_marker":"","next_version_id_marker":"","is_truncated":false}`),
		ok(`{"key":"img/logo.png","size_bytes":6,"etag":"e1","content_type":"image/png","cache_control":"max-age=60","content_disposition":"inline","last_modified":"2026-11-24T12:00:00Z","storage_class":"STANDARD","version_id":"v1","metadata":{"author":"ops"},"tags":[{"key":"env","value":"prod"}],"retention_mode":"GOVERNANCE","retain_until":"2026-12-24T12:00:00Z","legal_hold":true,"public":true,"object_public":true,"public_source":"object","public_url":"https://api.pantechdynamics.com/storage/public/bkt_1/img/logo.png"}`),
		cannedResponse{status: 200, body: string(rawContent), header: map[string]string{"Content-Type": "image/png"}},
		ok(`{"direct":true,"url":"https://s3.af-abj.pantechdynamics.com/signed-get","expires_in_seconds":300}`),
		ok(`{"key":"img/logo.png","size_bytes":6,"etag":"e1","content_type":"image/png","last_modified":"2026-11-24T12:00:00Z","version_id":"v1","public":true,"object_public":true,"public_source":"object","public_url":"https://api.pantechdynamics.com/storage/public/bkt_1/img/logo.png"}`),
		ok(`{"key":"raw.bin","size_bytes":6,"etag":"e2","content_type":"application/octet-stream","last_modified":"2026-11-24T12:00:00Z","version_id":"v2","public":false,"object_public":false,"public_source":"private","public_url":""}`),
		ok(`{"direct":true,"url":"https://s3.af-abj.pantechdynamics.com/signed-put","expires_in_seconds":300}`),
		ok(`{"prefix":"docs/"}`),
		ok(`{"key":"img/logo-copy.png","etag":"e1","version_id":"v2","last_modified":"2026-11-24T12:01:00Z","public":false,"object_public":false,"public_source":"private","public_url":""}`),
		ok(`{"deleted":1,"failed":[{"key":"locked.txt","version_id":"v1","code":"AccessDenied","locked":true}]}`),
		ok(`{"tags":[{"key":"tier","value":"hot"}]}`),
		ok(`{"ok":true}`),
		ok(`{"key":"img/logo.png","public":false,"object_public":false,"public_source":"private","public_url":""}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	ctx := context.Background()

	listed, err := c.Buckets.ListObjects(ctx, "bkt_1", &ListBucketObjectsParams{
		Prefix:    "img/",
		Delimiter: Ptr(""),
		Limit:     50,
	})
	if err != nil || len(listed.Objects) != 1 || listed.NextCursor != "c2" || !listed.IsTruncated {
		t.Fatalf("list objects = %+v err = %v", listed, err)
	}
	if q := rec.request(0).URL.Query(); !q.Has("delimiter") || q.Get("delimiter") != "" || q.Get("prefix") != "img/" || q.Get("limit") != "50" {
		t.Fatalf("unexpected list query: %v", q)
	}

	vers, err := c.Buckets.ListObjects(ctx, "bkt_1", &ListBucketObjectsParams{Versions: true, Prefix: "img/"})
	if err != nil || len(vers.Versions) != 1 || vers.Versions[0].VersionID != "v1" {
		t.Fatalf("list versions = %+v err = %v", vers, err)
	}

	meta, err := c.Buckets.GetObjectMeta(ctx, "bkt_1", &BucketObjectMetaParams{Key: "img/logo.png", VersionID: "v1"})
	if err != nil || meta.Metadata["author"] != "ops" || !meta.LegalHold || meta.RetainUntil == nil {
		t.Fatalf("meta = %+v err = %v", meta, err)
	}

	dl, err := c.Buckets.DownloadObject(ctx, "bkt_1", &DownloadBucketObjectParams{Key: "img/logo.png", Inline: true})
	if err != nil || !bytes.Equal(dl, rawContent) {
		t.Fatalf("download len=%d err=%v", len(dl), err)
	}
	if got := rec.request(3).Header.Get("Accept"); got != "application/octet-stream" {
		t.Fatalf("Accept = %q, want application/octet-stream", got)
	}

	dlURL, err := c.Buckets.PresignDownload(ctx, "bkt_1", PresignObjectDownloadRequest{Key: "img/logo.png", ExpiresInSeconds: Ptr(300)})
	if err != nil || !dlURL.Direct || dlURL.ExpiresInSeconds != 300 {
		t.Fatalf("presign download = %+v err = %v", dlURL, err)
	}

	up, err := c.Buckets.UploadObject(ctx, "bkt_1", &UploadBucketObjectParams{
		Key:                "img/logo.png",
		ContentType:        "image/png",
		CacheControl:       "max-age=60",
		ContentDisposition: "inline",
		Public:             true,
		Metadata:           map[string]string{"author": "ops"},
	}, rawContent)
	if err != nil || up.Key != "img/logo.png" || up.SizeBytes != 6 {
		t.Fatalf("upload = %+v err = %v", up, err)
	}
	upReq := rec.request(5)
	if upReq.Header.Get("Content-Type") != "image/png" ||
		upReq.Header.Get("X-Object-Cache-Control") != "max-age=60" ||
		upReq.Header.Get("X-Object-Content-Disposition") != "inline" ||
		upReq.Header.Get("X-Object-Public") != "true" ||
		upReq.Header.Get("X-Object-Meta-Author") != "ops" {
		t.Fatalf("unexpected upload headers: %v", upReq.Header)
	}
	if !bytes.Equal([]byte(rec.bodies[5]), rawContent) {
		t.Fatalf("upload body corrupted: got %x want %x", []byte(rec.bodies[5]), rawContent)
	}

	// Default Content-Type for raw []byte upload is application/octet-stream.
	if _, err := c.Buckets.UploadObject(ctx, "bkt_1", &UploadBucketObjectParams{Key: "raw.bin"}, rawContent); err != nil {
		t.Fatal(err)
	}
	if got := rec.request(6).Header.Get("Content-Type"); got != "application/octet-stream" {
		t.Fatalf("default binary Content-Type = %q, want application/octet-stream", got)
	}

	upURL, err := c.Buckets.PresignUpload(ctx, "bkt_1", PresignObjectUploadRequest{Key: "img/logo.png", ContentType: Ptr("image/png")})
	if err != nil || !upURL.Direct {
		t.Fatalf("presign upload = %+v err = %v", upURL, err)
	}

	folder, err := c.Buckets.CreateFolder(ctx, "bkt_1", CreateBucketFolderRequest{Prefix: "docs/"})
	if err != nil || folder.Prefix != "docs/" {
		t.Fatalf("folder = %+v err = %v", folder, err)
	}

	copied, err := c.Buckets.CopyObject(ctx, "bkt_1", CopyBucketObjectRequest{
		SourceKey:      "img/logo.png",
		DestinationKey: "img/logo-copy.png",
		DeleteSource:   Ptr(false),
	})
	if err != nil || copied.Key != "img/logo-copy.png" {
		t.Fatalf("copied = %+v err = %v", copied, err)
	}

	del, err := c.Buckets.DeleteObjects(ctx, "bkt_1", DeleteBucketObjectsRequest{
		Items: []DeleteBucketObjectItem{{Key: "img/logo-copy.png"}, {Key: "locked.txt", VersionID: Ptr("v1")}},
	})
	if err != nil || del.Deleted != 1 || len(del.Failed) != 1 || !del.Failed[0].Locked {
		t.Fatalf("delete objects = %+v err = %v", del, err)
	}

	otags, err := c.Buckets.PutObjectTags(ctx, "bkt_1", PutBucketObjectTagsRequest{
		Key:  "img/logo.png",
		Tags: []StorageTag{{Key: "tier", Value: "hot"}},
	})
	if err != nil || len(otags.Tags) != 1 || otags.Tags[0].Value != "hot" {
		t.Fatalf("put object tags = %+v err = %v", otags, err)
	}

	until := time.Date(2026, 12, 24, 12, 0, 0, 0, time.UTC)
	olock, err := c.Buckets.PutObjectLock(ctx, "bkt_1", PutBucketObjectLockRequest{
		Key:           "img/logo.png",
		RetentionMode: Ptr("GOVERNANCE"),
		RetainUntil:   &until,
		LegalHold:     Ptr(true),
	})
	if err != nil || !olock.OK {
		t.Fatalf("put object lock = %+v err = %v", olock, err)
	}

	vis, err := c.Buckets.PutObjectVisibility(ctx, "bkt_1", PutBucketObjectVisibilityRequest{
		Key:    "img/logo.png",
		Public: false,
	})
	if err != nil || vis.Public || vis.PublicSource != "private" {
		t.Fatalf("put object visibility = %+v err = %v", vis, err)
	}
}

func TestBucketMultipartFlow(t *testing.T) {
	partBytes := []byte("chunk-1-binary-payload")
	rec := newRecorder(t,
		ok(`{"upload_id":"up_1","mode":"proxy","part_size_bytes":8388608}`),
		ok(`{"direct":true,"url":"https://s3.af-abj.pantechdynamics.com/part-1","expires_in_seconds":900}`),
		ok(`{"part_number":1,"etag":"etag_part_1"}`),
		ok(`{"key":"archive.tar.zst","etag":"etag_final","version_id":"v1","last_modified":"2026-11-24T13:00:00Z","public":false,"object_public":false,"public_source":"private","public_url":""}`),
		status(204, ``),
	)
	c := newTestClientAt(t, rec.srv.URL)
	ctx := context.Background()

	sess, err := c.Buckets.StartMultipartUpload(ctx, "bkt_1", StartMultipartUploadRequest{
		Key:         "archive.tar.zst",
		ContentType: Ptr("application/zstd"),
		TotalSize:   Ptr(int64(16777216)),
	})
	if err != nil || sess.UploadID != "up_1" || sess.PartSizeBytes != 8388608 {
		t.Fatalf("start = %+v err = %v", sess, err)
	}

	signed, err := c.Buckets.SignMultipartPart(ctx, "bkt_1", SignMultipartPartRequest{
		Key:        "archive.tar.zst",
		UploadID:   sess.UploadID,
		PartNumber: 1,
	})
	if err != nil || !signed.Direct || signed.ExpiresInSeconds != 900 {
		t.Fatalf("sign part = %+v err = %v", signed, err)
	}

	part, err := c.Buckets.UploadMultipartPart(ctx, "bkt_1", &UploadMultipartPartParams{
		Key:        "archive.tar.zst",
		UploadID:   sess.UploadID,
		PartNumber: 1,
	}, partBytes)
	if err != nil || part.PartNumber != 1 || part.ETag != "etag_part_1" {
		t.Fatalf("upload part = %+v err = %v", part, err)
	}
	partReq := rec.request(2)
	if partReq.Header.Get("Content-Type") != "application/octet-stream" {
		t.Fatalf("part Content-Type = %q", partReq.Header.Get("Content-Type"))
	}
	if q := partReq.URL.Query(); q.Get("key") != "archive.tar.zst" || q.Get("upload_id") != "up_1" || q.Get("part_number") != "1" {
		t.Fatalf("part query = %v", q)
	}
	if !bytes.Equal([]byte(rec.bodies[2]), partBytes) {
		t.Fatalf("part body = %q, want %q", rec.bodies[2], partBytes)
	}

	done, err := c.Buckets.CompleteMultipartUpload(ctx, "bkt_1", CompleteMultipartUploadRequest{
		Key:      "archive.tar.zst",
		UploadID: sess.UploadID,
		Parts:    []MultipartPart{*part},
	})
	if err != nil || done.Key != "archive.tar.zst" || done.ETag != "etag_final" {
		t.Fatalf("complete = %+v err = %v", done, err)
	}

	if err := c.Buckets.AbortMultipartUpload(ctx, "bkt_1", AbortMultipartUploadRequest{
		Key:      "archive.tar.zst",
		UploadID: "up_2",
	}); err != nil {
		t.Fatalf("abort: %v", err)
	}
}
