package pantech

import (
	"context"
	"iter"
	"net/url"
)

// VolumesService manages data disks.
type VolumesService struct{ c *Client }

// List lists volumes: GET /volumes.
// Scope: read.
func (s *VolumesService) List(ctx context.Context, params *ListParams, opts ...RequestOption) (*Page[Volume], error) {
	var out Page[Volume]
	if err := s.c.Do(ctx, "GET", "/volumes", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// All iterates over every item of List, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *VolumesService) All(ctx context.Context, params *ListParams, opts ...RequestOption) iter.Seq2[Volume, error] {
	var p ListParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Volume], error) {
		p.Cursor = cursor
		return s.List(ctx, &p, opts...)
	})
}

// Create creates a volume: POST /volumes.
// Scope: write.
func (s *VolumesService) Create(ctx context.Context, req CreateVolumeRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/volumes", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a volume: GET /volumes/{volume_id}.
// Scope: read.
func (s *VolumesService) Get(ctx context.Context, volumeID string, opts ...RequestOption) (*Volume, error) {
	var out Volume
	if err := s.c.Do(ctx, "GET", "/volumes/"+url.PathEscape(volumeID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a volume: DELETE /volumes/{volume_id}.
// Scope: write.
func (s *VolumesService) Delete(ctx context.Context, volumeID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", "/volumes/"+url.PathEscape(volumeID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Attach attaches a volume: POST /volumes/{volume_id}/attach.
// Scope: write.
func (s *VolumesService) Attach(ctx context.Context, volumeID string, req AttachVolumeRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/volumes/"+url.PathEscape(volumeID)+"/attach", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Detach detaches a volume: POST /volumes/{volume_id}/detach.
// Scope: write.
func (s *VolumesService) Detach(ctx context.Context, volumeID string, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/volumes/"+url.PathEscape(volumeID)+"/detach", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Resize grows a volume: POST /volumes/{volume_id}/resize.
// Scope: write.
func (s *VolumesService) Resize(ctx context.Context, volumeID string, req ResizeVolumeRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/volumes/"+url.PathEscape(volumeID)+"/resize", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSnapshotSchedule gets a volume's snapshot schedule: GET /volumes/{volume_id}/snapshot-schedule.
// Scope: read.
func (s *VolumesService) GetSnapshotSchedule(ctx context.Context, volumeID string, opts ...RequestOption) (*SnapshotSchedule, error) {
	var out SnapshotSchedule
	if err := s.c.Do(ctx, "GET", "/volumes/"+url.PathEscape(volumeID)+"/snapshot-schedule", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// PutSnapshotSchedule sets a volume's snapshot schedule: PUT /volumes/{volume_id}/snapshot-schedule.
// Scope: write.
func (s *VolumesService) PutSnapshotSchedule(ctx context.Context, volumeID string, req PutSnapshotScheduleRequest, opts ...RequestOption) (*SnapshotSchedule, error) {
	var out SnapshotSchedule
	if err := s.c.Do(ctx, "PUT", "/volumes/"+url.PathEscape(volumeID)+"/snapshot-schedule", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListSnapshots lists a volume's snapshots: GET /volumes/{volume_id}/snapshots.
// Scope: read.
func (s *VolumesService) ListSnapshots(ctx context.Context, volumeID string, params *ListParams, opts ...RequestOption) (*Page[Snapshot], error) {
	var out Page[Snapshot]
	if err := s.c.Do(ctx, "GET", "/volumes/"+url.PathEscape(volumeID)+"/snapshots", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllSnapshots iterates over every item of ListSnapshots, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *VolumesService) AllSnapshots(ctx context.Context, volumeID string, params *ListParams, opts ...RequestOption) iter.Seq2[Snapshot, error] {
	var p ListParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Snapshot], error) {
		p.Cursor = cursor
		return s.ListSnapshots(ctx, volumeID, &p, opts...)
	})
}

// CreateSnapshot snapshots a volume: POST /volumes/{volume_id}/snapshots.
// Scope: write.
func (s *VolumesService) CreateSnapshot(ctx context.Context, volumeID string, req CreateSnapshotRequest, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "POST", "/volumes/"+url.PathEscape(volumeID)+"/snapshots", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUsage gets a volume's usage: GET /volumes/{volume_id}/usage.
// Scope: read.
func (s *VolumesService) GetUsage(ctx context.Context, volumeID string, opts ...RequestOption) (*UsageSummary, error) {
	var out UsageSummary
	if err := s.c.Do(ctx, "GET", "/volumes/"+url.PathEscape(volumeID)+"/usage", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
