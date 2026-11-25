package pantech

import (
	"context"
	"net/url"
)

// Database import modes.
const (
	DatabaseImportDump      = "dump"
	DatabaseImportReplicate = "replicate"
)

// StartDatabaseImport starts moving one of a plan's databases into one of
// your Pantech databases, through the Migrate agent running in your network:
// POST /migrations/{migration_id}/database-imports. The agent reads the
// source with credentials that stay on it, checks it against the destination
// first and reports what to fix in Readiness. The move is free; the
// destination bills as usual.
// Scope: write.
func (s *MigrateService) StartDatabaseImport(ctx context.Context, migrationID string, req StartDatabaseImportRequest, opts ...RequestOption) (*DatabaseImport, error) {
	var out DatabaseImport
	if err := s.c.Do(ctx, "POST", "/migrations/"+url.PathEscape(migrationID)+"/database-imports", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListDatabaseImports lists the organization's database imports, newest
// first, without their tables: GET /database-imports.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *MigrateService) ListDatabaseImports(ctx context.Context, opts ...RequestOption) ([]DatabaseImport, error) {
	return listAll[DatabaseImport](ctx, s.c, "/database-imports", nil, opts...)
}

// GetDatabaseImport gets an import with its readiness report, per-table
// progress and verification, and what the destination reported about itself:
// GET /database-imports/{import_id}.
// Scope: read.
func (s *MigrateService) GetDatabaseImport(ctx context.Context, importID string, opts ...RequestOption) (*DatabaseImport, error) {
	var out DatabaseImport
	if err := s.c.Do(ctx, "GET", "/database-imports/"+url.PathEscape(importID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// CheckDatabaseImportReadiness asks the agent to run the pre-flight again,
// for example after fixing what it reported: POST /database-imports/{import_id}/readiness.
// Allowed before the copy starts.
// Scope: write.
func (s *MigrateService) CheckDatabaseImportReadiness(ctx context.Context, importID string, opts ...RequestOption) (*DatabaseImport, error) {
	return s.actOnDatabaseImport(ctx, importID, "readiness", opts)
}

// CutOverDatabaseImport finishes an import: POST /database-imports/{import_id}/cutover.
// For dump, once every table is copied; for replicate, stop writes to the
// source first.
// Scope: write.
func (s *MigrateService) CutOverDatabaseImport(ctx context.Context, importID string, opts ...RequestOption) (*DatabaseImport, error) {
	return s.actOnDatabaseImport(ctx, importID, "cutover", opts)
}

// CancelDatabaseImport stops an import and removes its access to your
// Pantech database: POST /database-imports/{import_id}/cancel. What was
// already copied stays.
// Scope: write.
func (s *MigrateService) CancelDatabaseImport(ctx context.Context, importID string, opts ...RequestOption) (*DatabaseImport, error) {
	return s.actOnDatabaseImport(ctx, importID, "cancel", opts)
}

// RevokeDatabaseImportAccess removes the import login and its loopback-only
// access from your Pantech database now and keeps the data:
// POST /database-imports/{import_id}/revoke-access. An import still copying stops.
// Scope: write.
func (s *MigrateService) RevokeDatabaseImportAccess(ctx context.Context, importID string, opts ...RequestOption) (*DatabaseImport, error) {
	return s.actOnDatabaseImport(ctx, importID, "revoke-access", opts)
}

func (s *MigrateService) actOnDatabaseImport(ctx context.Context, importID, action string, opts []RequestOption) (*DatabaseImport, error) {
	var out DatabaseImport
	if err := s.c.Do(ctx, "POST", "/database-imports/"+url.PathEscape(importID)+"/"+action, nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
