package pantech

import (
	"context"
	"net/url"
)

// ObjectStorageService reports object storage status and usage and manages S3
// access keys. Buckets and their objects are managed through Client.Buckets.
type ObjectStorageService struct{ c *Client }

// GetStatus gets object storage status and the S3 endpoint: GET /object-storage.
// Scope: read.
func (s *ObjectStorageService) GetStatus(ctx context.Context, opts ...RequestOption) (*ObjectStorageStatus, error) {
	var out ObjectStorageStatus
	if err := s.c.Do(ctx, "GET", "/object-storage", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUsage gets hourly object storage transfer and request usage: GET /object-storage/usage.
// Scope: read.
func (s *ObjectStorageService) GetUsage(ctx context.Context, params *ObjectStorageUsageParams, opts ...RequestOption) (*ObjectStorageUsage, error) {
	var out ObjectStorageUsage
	if err := s.c.Do(ctx, "GET", "/object-storage/usage", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSettings changes organization-wide object storage settings such as
// Block All Public Access: PATCH /object-storage/settings.
// Scope: write.
func (s *ObjectStorageService) UpdateSettings(ctx context.Context, req UpdateObjectStorageSettingsRequest, opts ...RequestOption) (*ObjectStorageSettings, error) {
	var out ObjectStorageSettings
	if err := s.c.Do(ctx, "PATCH", "/object-storage/settings", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListKeys lists active object storage access keys without their secrets: GET /object-storage/keys.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *ObjectStorageService) ListKeys(ctx context.Context, params *ListObjectStorageKeysParams, opts ...RequestOption) ([]ObjectStorageKey, error) {
	return listAll[ObjectStorageKey](ctx, s.c, "/object-storage/keys", params.values(), opts...)
}

// CreateKey creates an S3 access key: POST /object-storage/keys.
// Its SecretAccessKey is returned in this response only.
// Scope: write.
func (s *ObjectStorageService) CreateKey(ctx context.Context, req CreateObjectStorageKeyRequest, opts ...RequestOption) (*CreatedObjectStorageKey, error) {
	var out CreatedObjectStorageKey
	if err := s.c.Do(ctx, "POST", "/object-storage/keys", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// RevokeKey revokes an S3 access key: DELETE /object-storage/keys/{storage_key_id}.
// Scope: write.
func (s *ObjectStorageService) RevokeKey(ctx context.Context, storageKeyID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", "/object-storage/keys/"+url.PathEscape(storageKeyID), nil, nil, nil, opts...)
}
