package pantech

import (
	"context"
	"net/url"
)

// Bucket import statuses. Succeeded, failed and cancelled are final; a
// succeeded import can be resynced.
const (
	BucketImportQueued    = "queued"
	BucketImportRunning   = "running"
	BucketImportSucceeded = "succeeded"
	BucketImportFailed    = "failed"
	BucketImportCancelled = "cancelled"
)

// GetCopyAccess gets a connection's opt-in copy access: the buckets chosen
// for copying, the copy role and what its last verification found, and the
// template for that role: GET /cloud-connections/{connection_id}/copy-access.
// Status is not_requested until buckets are chosen.
// Scope: read.
func (s *MigrateService) GetCopyAccess(ctx context.Context, connectionID string, opts ...RequestOption) (*CopyAccess, error) {
	var out CopyAccess
	if err := s.c.Do(ctx, "GET", "/cloud-connections/"+url.PathEscape(connectionID)+"/copy-access", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetCopyAccess chooses the buckets Pantech may copy and returns the template
// of a copy role that may read only those: PUT /cloud-connections/{connection_id}/copy-access.
// Create a stack from it, then send its CopyRoleArn output as RoleARN;
// Pantech then verifies the role against each bucket without reading any
// object. Retry verifies again.
// Scope: write.
func (s *MigrateService) SetCopyAccess(ctx context.Context, connectionID string, req CopyAccessRequest, opts ...RequestOption) (*CopyAccess, error) {
	var out CopyAccess
	if err := s.c.Do(ctx, "PUT", "/cloud-connections/"+url.PathEscape(connectionID)+"/copy-access", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListBucketImports lists the organization's bucket imports, newest first,
// without their failures and events: GET /bucket-imports.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *MigrateService) ListBucketImports(ctx context.Context, opts ...RequestOption) ([]BucketImport, error) {
	return listAll[BucketImport](ctx, s.c, "/bucket-imports", nil, opts...)
}

// StartBucketImport copies one source bucket's objects into a new or existing
// Pantech bucket: POST /bucket-imports. The bucket must have verified copy
// access. Only current versions are copied and nothing is deleted on either
// side. Follow it with GetBucketImport. Moving is free; stored data bills as
// usual.
// Scope: write.
func (s *MigrateService) StartBucketImport(ctx context.Context, req StartBucketImportRequest, opts ...RequestOption) (*BucketImport, error) {
	var out BucketImport
	if err := s.c.Do(ctx, "POST", "/bucket-imports", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetBucketImport gets an import's status and progress, the first 100
// objects that failed its pass, and its audit trail: GET /bucket-imports/{import_id}.
// Scope: read.
func (s *MigrateService) GetBucketImport(ctx context.Context, importID string, opts ...RequestOption) (*BucketImport, error) {
	var out BucketImport
	if err := s.c.Do(ctx, "GET", "/bucket-imports/"+url.PathEscape(importID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelBucketImport stops an import: a queued one at once, a running one at
// its next checkpoint: POST /bucket-imports/{import_id}/cancel. Objects
// already copied stay.
// Scope: write.
func (s *MigrateService) CancelBucketImport(ctx context.Context, importID string, opts ...RequestOption) (*BucketImport, error) {
	var out BucketImport
	if err := s.c.Do(ctx, "POST", "/bucket-imports/"+url.PathEscape(importID)+"/cancel", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ResyncBucketImport queues another pass that copies only what changed:
// POST /bucket-imports/{import_id}/resync. Run a last resync after writes to
// the source stop, then switch over. MirrorDeletes also deletes the objects
// this import copied that the source no longer has, and needs the delete
// scope.
// Scope: write.
func (s *MigrateService) ResyncBucketImport(ctx context.Context, importID string, req ResyncBucketImportRequest, opts ...RequestOption) (*BucketImport, error) {
	var out BucketImport
	if err := s.c.Do(ctx, "POST", "/bucket-imports/"+url.PathEscape(importID)+"/resync", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
