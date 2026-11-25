package pantech

import (
	"context"
	"net/url"
)

// Server copy statuses. Completed, failed and cancelled are final; synced
// waits for SyncServerCopy or CutOverServerCopy.
const (
	ServerCopyWaiting          = "waiting"
	ServerCopyCopying          = "copying"
	ServerCopySynced           = "synced"
	ServerCopyImporting        = "importing"
	ServerCopyCreatingInstance = "creating_instance"
	ServerCopyCompleted        = "completed"
	ServerCopyFailed           = "failed"
	ServerCopyCancelled        = "cancelled"
)

// Server copy sources.
const (
	ServerCopySourceAgent       = "agent"
	ServerCopySourceAWSSnapshot = "aws_snapshot"
)

// StartServerCopy starts copying one planned server's disks into a new
// Pantech instance: POST /migrations/{migration_id}/server-copies. Nothing is
// created yet: the first pass copies every disk, SyncServerCopy copies what
// changed since, and CutOverServerCopy runs the last pass and builds the
// instance through the normal create path, so it is billed like any new
// instance. The copy itself is free.
// Scope: write.
func (s *MigrateService) StartServerCopy(ctx context.Context, migrationID string, req StartServerCopyRequest, opts ...RequestOption) (*ServerCopy, error) {
	var out ServerCopy
	if err := s.c.Do(ctx, "POST", "/migrations/"+url.PathEscape(migrationID)+"/server-copies", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetServerCopyTemplate returns the template of the AWS copy role for one
// planned server, for a copy with source aws_snapshot:
// GET /migrations/{migration_id}/server-copy-template. params.ItemKey is required.
// Scope: read.
func (s *MigrateService) GetServerCopyTemplate(ctx context.Context, migrationID string, params *GetServerCopyTemplateParams, opts ...RequestOption) (*ServerCopyTemplate, error) {
	var out ServerCopyTemplate
	if err := s.c.Do(ctx, "GET", "/migrations/"+url.PathEscape(migrationID)+"/server-copy-template", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListServerCopies lists every server copy in the organization, newest
// first, with its progress: GET /server-copies.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *MigrateService) ListServerCopies(ctx context.Context, opts ...RequestOption) ([]ServerCopy, error) {
	return listAll[ServerCopy](ctx, s.c, "/server-copies", nil, opts...)
}

// GetServerCopy gets a server copy: GET /server-copies/{copy_id}.
// Scope: read.
func (s *MigrateService) GetServerCopy(ctx context.Context, copyID string, opts ...RequestOption) (*ServerCopy, error) {
	var out ServerCopy
	if err := s.c.Do(ctx, "GET", "/server-copies/"+url.PathEscape(copyID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// SyncServerCopy asks for another pass that copies only what changed since
// the last one: POST /server-copies/{copy_id}/sync. For a copy from AWS
// snapshots that makes none itself, name each disk's new snapshot in
// req.Snapshots.
// Scope: write.
func (s *MigrateService) SyncServerCopy(ctx context.Context, copyID string, req ServerCopyActRequest, opts ...RequestOption) (*ServerCopy, error) {
	return s.actOnServerCopy(ctx, copyID, "sync", req, opts)
}

// CutOverServerCopy asks for the last pass, then Pantech builds the instance:
// POST /server-copies/{copy_id}/cutover. Stop the application first.
// Scope: write.
func (s *MigrateService) CutOverServerCopy(ctx context.Context, copyID string, req ServerCopyActRequest, opts ...RequestOption) (*ServerCopy, error) {
	return s.actOnServerCopy(ctx, copyID, "cutover", req, opts)
}

// CancelServerCopy stops a copy and deletes everything Pantech staged for it:
// POST /server-copies/{copy_id}/cancel. Refused once the instance is being
// created: it is yours then, delete it as any instance.
// Scope: write.
func (s *MigrateService) CancelServerCopy(ctx context.Context, copyID string, req ServerCopyActRequest, opts ...RequestOption) (*ServerCopy, error) {
	return s.actOnServerCopy(ctx, copyID, "cancel", req, opts)
}

func (s *MigrateService) actOnServerCopy(ctx context.Context, copyID, action string, req ServerCopyActRequest, opts []RequestOption) (*ServerCopy, error) {
	var out ServerCopy
	if err := s.c.Do(ctx, "POST", "/server-copies/"+url.PathEscape(copyID)+"/"+action, nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
