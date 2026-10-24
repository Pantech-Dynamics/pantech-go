package pantech

import (
	"context"
	"iter"
	"net/url"
	"strconv"
)

// SnapshotsService manages snapshots of root disks, volumes and databases.
type SnapshotsService struct{ c *Client }

// List lists snapshots: GET /snapshots.
// Scope: read.
func (s *SnapshotsService) List(ctx context.Context, params *ListParams, opts ...RequestOption) (*Page[Snapshot], error) {
	var out Page[Snapshot]
	if err := s.c.Do(ctx, "GET", "/snapshots", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// All iterates over every item of List, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *SnapshotsService) All(ctx context.Context, params *ListParams, opts ...RequestOption) iter.Seq2[Snapshot, error] {
	var p ListParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Snapshot], error) {
		p.Cursor = cursor
		return s.List(ctx, &p, opts...)
	})
}

// Quote prices a snapshot: GET /snapshots/quote.
// Scope: read.
func (s *SnapshotsService) Quote(ctx context.Context, sizeGB int, opts ...RequestOption) (*OptionalQuote, error) {
	var out OptionalQuote
	if err := s.c.Do(ctx, "GET", "/snapshots/quote", url.Values{"size_gb": {strconv.Itoa(sizeGB)}}, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a snapshot: GET /snapshots/{snapshot_id}.
// Scope: read.
func (s *SnapshotsService) Get(ctx context.Context, snapshotID string, opts ...RequestOption) (*Snapshot, error) {
	var out Snapshot
	if err := s.c.Do(ctx, "GET", "/snapshots/"+url.PathEscape(snapshotID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a snapshot: DELETE /snapshots/{snapshot_id}.
// Scope: write.
func (s *SnapshotsService) Delete(ctx context.Context, snapshotID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/snapshots/"+url.PathEscape(snapshotID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Restore restores a snapshot to a new volume: POST /snapshots/{snapshot_id}/restore.
// Scope: write.
func (s *SnapshotsService) Restore(ctx context.Context, snapshotID string, req RestoreSnapshotRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/snapshots/"+url.PathEscape(snapshotID)+"/restore", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUsage gets a snapshot's usage: GET /snapshots/{snapshot_id}/usage.
// Scope: read.
func (s *SnapshotsService) GetUsage(ctx context.Context, snapshotID string, opts ...RequestOption) (*UsageSummary, error) {
	var out UsageSummary
	if err := s.c.Do(ctx, "GET", "/snapshots/"+url.PathEscape(snapshotID)+"/usage", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
