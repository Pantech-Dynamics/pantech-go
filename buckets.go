package pantech

import (
	"context"
	"net/url"
)

// BucketsService manages object storage buckets, bucket rules and objects.
type BucketsService struct{ c *Client }

func bucketPath(bucketID string) string {
	return "/buckets/" + url.PathEscape(bucketID)
}

// List lists buckets: GET /buckets.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *BucketsService) List(ctx context.Context, opts ...RequestOption) ([]Bucket, error) {
	return listAll[Bucket](ctx, s.c, "/buckets", nil, opts...)
}

// Create creates a bucket: POST /buckets.
// Scope: write.
func (s *BucketsService) Create(ctx context.Context, req CreateBucketRequest, opts ...RequestOption) (*CreatedBucket, error) {
	var out CreatedBucket
	if err := s.c.Do(ctx, "POST", "/buckets", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a bucket: GET /buckets/{bucket_id}.
// Scope: read.
func (s *BucketsService) Get(ctx context.Context, bucketID string, opts ...RequestOption) (*Bucket, error) {
	var out Bucket
	if err := s.c.Do(ctx, "GET", bucketPath(bucketID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update changes a bucket's quota, versioning or public prefixes: PATCH /buckets/{bucket_id}.
// Scope: write.
func (s *BucketsService) Update(ctx context.Context, bucketID string, req UpdateBucketRequest, opts ...RequestOption) (*UpdatedBucket, error) {
	var out UpdatedBucket
	if err := s.c.Do(ctx, "PATCH", bucketPath(bucketID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a bucket: DELETE /buckets/{bucket_id}.
// Pass &DeleteBucketParams{Force: true} to delete a bucket that still holds objects.
// Scope: write.
func (s *BucketsService) Delete(ctx context.Context, bucketID string, params *DeleteBucketParams, opts ...RequestOption) (*Accepted, error) {
	var out Accepted
	if err := s.c.Do(ctx, "DELETE", bucketPath(bucketID), params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetStats gets live stored bytes and object count for a bucket: GET /buckets/{bucket_id}/stats.
// Scope: read.
func (s *BucketsService) GetStats(ctx context.Context, bucketID string, opts ...RequestOption) (*BucketStats, error) {
	var out BucketStats
	if err := s.c.Do(ctx, "GET", bucketPath(bucketID)+"/stats", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Empty deletes all objects, versions and delete markers and aborts incomplete multipart uploads:
// POST /buckets/{bucket_id}/empty.
// Scope: write.
func (s *BucketsService) Empty(ctx context.Context, bucketID string, opts ...RequestOption) (*EmptyBucketResult, error) {
	var out EmptyBucketResult
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/empty", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListObjects lists or searches objects and folders (or versions when Versions is true):
// GET /buckets/{bucket_id}/objects.
// Scope: read.
func (s *BucketsService) ListObjects(ctx context.Context, bucketID string, params *ListBucketObjectsParams, opts ...RequestOption) (*BucketObjectList, error) {
	var out BucketObjectList
	if err := s.c.Do(ctx, "GET", bucketPath(bucketID)+"/objects", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetObjectMeta gets an object's metadata, tags and object-lock status:
// GET /buckets/{bucket_id}/objects/meta.
// Scope: read.
func (s *BucketsService) GetObjectMeta(ctx context.Context, bucketID string, params *BucketObjectMetaParams, opts ...RequestOption) (*BucketObjectMeta, error) {
	var out BucketObjectMeta
	if err := s.c.Do(ctx, "GET", bucketPath(bucketID)+"/objects/meta", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadObject downloads an object's raw content: GET /buckets/{bucket_id}/objects/download.
// Scope: read.
func (s *BucketsService) DownloadObject(ctx context.Context, bucketID string, params *DownloadBucketObjectParams, opts ...RequestOption) ([]byte, error) {
	var out []byte
	allOpts := append([]RequestOption{WithRequestHeader("Accept", "application/octet-stream")}, opts...)
	if err := s.c.Do(ctx, "GET", bucketPath(bucketID)+"/objects/download", params.values(), nil, &out, allOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

// PresignDownload mints a short-lived direct download URL: POST /buckets/{bucket_id}/objects/download-url.
// Scope: read.
func (s *BucketsService) PresignDownload(ctx context.Context, bucketID string, req PresignObjectDownloadRequest, opts ...RequestOption) (*PresignedURL, error) {
	var out PresignedURL
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/objects/download-url", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UploadObject uploads raw object bytes into a bucket: PUT /buckets/{bucket_id}/objects/upload.
// Scope: write.
func (s *BucketsService) UploadObject(ctx context.Context, bucketID string, params *UploadBucketObjectParams, body []byte, opts ...RequestOption) (*UploadedBucketObject, error) {
	var out UploadedBucketObject
	allOpts := append(params.requestOptions(), opts...)
	if err := s.c.Do(ctx, "PUT", bucketPath(bucketID)+"/objects/upload", params.values(), body, &out, allOpts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// PresignUpload mints a short-lived direct upload URL: POST /buckets/{bucket_id}/objects/upload-url.
// Scope: write.
func (s *BucketsService) PresignUpload(ctx context.Context, bucketID string, req PresignObjectUploadRequest, opts ...RequestOption) (*PresignedURL, error) {
	var out PresignedURL
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/objects/upload-url", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateUploadForm signs an HTML form upload (S3 POST policy) for direct browser uploads:
// POST /buckets/{bucket_id}/objects/upload-form.
// Scope: write.
func (s *BucketsService) CreateUploadForm(ctx context.Context, bucketID string, req CreateBucketUploadFormRequest, opts ...RequestOption) (*BucketUploadForm, error) {
	var out BucketUploadForm
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/objects/upload-form", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateFolder creates a folder prefix marker: POST /buckets/{bucket_id}/objects/folder.
// Scope: write.
func (s *BucketsService) CreateFolder(ctx context.Context, bucketID string, req CreateBucketFolderRequest, opts ...RequestOption) (*BucketFolder, error) {
	var out BucketFolder
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/objects/folder", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// CopyObject copies, renames or replaces metadata on an object: POST /buckets/{bucket_id}/objects/copy.
// Scope: write.
func (s *BucketsService) CopyObject(ctx context.Context, bucketID string, req CopyBucketObjectRequest, opts ...RequestOption) (*BucketObjectSummary, error) {
	var out BucketObjectSummary
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/objects/copy", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteObjects deletes objects, versions or a folder prefix: POST /buckets/{bucket_id}/objects/delete.
// Scope: write.
func (s *BucketsService) DeleteObjects(ctx context.Context, bucketID string, req DeleteBucketObjectsRequest, opts ...RequestOption) (*DeleteBucketObjectsResult, error) {
	var out DeleteBucketObjectsResult
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/objects/delete", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// PutObjectTags replaces key-value tags on an object or version: PUT /buckets/{bucket_id}/objects/tags.
// Scope: write.
func (s *BucketsService) PutObjectTags(ctx context.Context, bucketID string, req PutBucketObjectTagsRequest, opts ...RequestOption) (*StorageTags, error) {
	var out StorageTags
	if err := s.c.Do(ctx, "PUT", bucketPath(bucketID)+"/objects/tags", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// PutObjectLock sets retention or legal hold on an object version: PUT /buckets/{bucket_id}/objects/lock.
// Scope: write.
func (s *BucketsService) PutObjectLock(ctx context.Context, bucketID string, req PutBucketObjectLockRequest, opts ...RequestOption) (*PutBucketObjectLockResult, error) {
	var out PutBucketObjectLockResult
	if err := s.c.Do(ctx, "PUT", bucketPath(bucketID)+"/objects/lock", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// PutObjectVisibility toggles public or private read access on an individual object:
// PUT /buckets/{bucket_id}/objects/visibility.
// Scope: write.
func (s *BucketsService) PutObjectVisibility(ctx context.Context, bucketID string, req PutBucketObjectVisibilityRequest, opts ...RequestOption) (*BucketObjectVisibility, error) {
	var out BucketObjectVisibility
	if err := s.c.Do(ctx, "PUT", bucketPath(bucketID)+"/objects/visibility", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// RestoreObjectVersion restores a previous object version or removes a delete marker:
// POST /buckets/{bucket_id}/objects/restore-version.
// Scope: write.
func (s *BucketsService) RestoreObjectVersion(ctx context.Context, bucketID string, req RestoreBucketObjectVersionRequest, opts ...RequestOption) (*RestoredBucketObjectVersion, error) {
	var out RestoredBucketObjectVersion
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/objects/restore-version", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListMultipartUploads lists active in-progress multipart uploads in a bucket:
// GET /buckets/{bucket_id}/objects/multipart.
// Scope: read.
func (s *BucketsService) ListMultipartUploads(ctx context.Context, bucketID string, opts ...RequestOption) ([]MultipartUpload, error) {
	var out struct {
		Uploads []MultipartUpload `json:"uploads"`
	}
	if err := s.c.Do(ctx, "GET", bucketPath(bucketID)+"/objects/multipart", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	if out.Uploads == nil {
		return []MultipartUpload{}, nil
	}
	return out.Uploads, nil
}

// StartMultipartUpload starts a multipart upload session: POST /buckets/{bucket_id}/objects/multipart/start.
// Scope: write.
func (s *BucketsService) StartMultipartUpload(ctx context.Context, bucketID string, req StartMultipartUploadRequest, opts ...RequestOption) (*MultipartUploadSession, error) {
	var out MultipartUploadSession
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/objects/multipart/start", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// SignMultipartPart signs a direct upload URL for one multipart part:
// POST /buckets/{bucket_id}/objects/multipart/sign-part.
// Scope: write.
func (s *BucketsService) SignMultipartPart(ctx context.Context, bucketID string, req SignMultipartPartRequest, opts ...RequestOption) (*PresignedURL, error) {
	var out PresignedURL
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/objects/multipart/sign-part", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// SignMultipartParts batch-signs direct upload URLs for up to 100 multipart parts:
// POST /buckets/{bucket_id}/objects/multipart/sign-parts.
// Scope: write.
func (s *BucketsService) SignMultipartParts(ctx context.Context, bucketID string, req SignMultipartPartsRequest, opts ...RequestOption) (*SignedMultipartParts, error) {
	var out SignedMultipartParts
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/objects/multipart/sign-parts", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UploadMultipartPart uploads one part of an active multipart upload:
// PUT /buckets/{bucket_id}/objects/multipart/part.
// Scope: write.
func (s *BucketsService) UploadMultipartPart(ctx context.Context, bucketID string, params *UploadMultipartPartParams, body []byte, opts ...RequestOption) (*MultipartPart, error) {
	var out MultipartPart
	if err := s.c.Do(ctx, "PUT", bucketPath(bucketID)+"/objects/multipart/part", params.values(), body, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// CompleteMultipartUpload assembles uploaded parts into the final object:
// POST /buckets/{bucket_id}/objects/multipart/complete.
// Scope: write.
func (s *BucketsService) CompleteMultipartUpload(ctx context.Context, bucketID string, req CompleteMultipartUploadRequest, opts ...RequestOption) (*BucketObjectSummary, error) {
	var out BucketObjectSummary
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/objects/multipart/complete", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AbortMultipartUpload cancels an in-progress multipart upload:
// POST /buckets/{bucket_id}/objects/multipart/abort.
// Scope: write.
func (s *BucketsService) AbortMultipartUpload(ctx context.Context, bucketID string, req AbortMultipartUploadRequest, opts ...RequestOption) error {
	return s.c.Do(ctx, "POST", bucketPath(bucketID)+"/objects/multipart/abort", nil, &req, nil, opts...)
}

// GetCORS gets a bucket's CORS rules: GET /buckets/{bucket_id}/cors.
// Scope: read.
func (s *BucketsService) GetCORS(ctx context.Context, bucketID string, opts ...RequestOption) (*BucketCORSConfig, error) {
	var out BucketCORSConfig
	if err := s.c.Do(ctx, "GET", bucketPath(bucketID)+"/cors", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// PutCORS replaces a bucket's CORS rules: PUT /buckets/{bucket_id}/cors.
// Scope: write.
func (s *BucketsService) PutCORS(ctx context.Context, bucketID string, req BucketCORSConfig, opts ...RequestOption) (*BucketCORSConfig, error) {
	var out BucketCORSConfig
	if err := s.c.Do(ctx, "PUT", bucketPath(bucketID)+"/cors", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLifecycle gets a bucket's lifecycle rules: GET /buckets/{bucket_id}/lifecycle.
// Scope: read.
func (s *BucketsService) GetLifecycle(ctx context.Context, bucketID string, opts ...RequestOption) (*BucketLifecycleConfig, error) {
	var out BucketLifecycleConfig
	if err := s.c.Do(ctx, "GET", bucketPath(bucketID)+"/lifecycle", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// PutLifecycle replaces a bucket's lifecycle rules: PUT /buckets/{bucket_id}/lifecycle.
// Scope: write.
func (s *BucketsService) PutLifecycle(ctx context.Context, bucketID string, req BucketLifecycleConfig, opts ...RequestOption) (*BucketLifecycleConfig, error) {
	var out BucketLifecycleConfig
	if err := s.c.Do(ctx, "PUT", bucketPath(bucketID)+"/lifecycle", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTags gets a bucket's tags: GET /buckets/{bucket_id}/tags.
// Scope: read.
func (s *BucketsService) GetTags(ctx context.Context, bucketID string, opts ...RequestOption) (*BucketTagsConfig, error) {
	var out BucketTagsConfig
	if err := s.c.Do(ctx, "GET", bucketPath(bucketID)+"/tags", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// PutTags replaces a bucket's tags: PUT /buckets/{bucket_id}/tags.
// Scope: write.
func (s *BucketsService) PutTags(ctx context.Context, bucketID string, req BucketTagsConfig, opts ...RequestOption) (*BucketTagsConfig, error) {
	var out BucketTagsConfig
	if err := s.c.Do(ctx, "PUT", bucketPath(bucketID)+"/tags", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetObjectLock gets default object-lock retention settings for a bucket:
// GET /buckets/{bucket_id}/object-lock.
// Scope: read.
func (s *BucketsService) GetObjectLock(ctx context.Context, bucketID string, opts ...RequestOption) (*BucketObjectLockConfig, error) {
	var out BucketObjectLockConfig
	if err := s.c.Do(ctx, "GET", bucketPath(bucketID)+"/object-lock", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// PutObjectLockConfig sets default object-lock retention settings for a bucket:
// PUT /buckets/{bucket_id}/object-lock.
// Scope: write.
func (s *BucketsService) PutObjectLockConfig(ctx context.Context, bucketID string, req BucketObjectLockConfig, opts ...RequestOption) (*BucketObjectLockConfig, error) {
	var out BucketObjectLockConfig
	if err := s.c.Do(ctx, "PUT", bucketPath(bucketID)+"/object-lock", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListShared lists buckets shared with your organization by other Pantech organizations:
// GET /buckets/shared.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *BucketsService) ListShared(ctx context.Context, opts ...RequestOption) ([]SharedBucket, error) {
	return listAll[SharedBucket](ctx, s.c, "/buckets/shared", nil, opts...)
}

// ListShares lists cross-organization shares configured on a bucket:
// GET /buckets/{bucket_id}/shares.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *BucketsService) ListShares(ctx context.Context, bucketID string, opts ...RequestOption) ([]BucketShare, error) {
	return listAll[BucketShare](ctx, s.c, bucketPath(bucketID)+"/shares", nil, opts...)
}

// CreateShare shares a bucket with another Pantech organization:
// POST /buckets/{bucket_id}/shares.
// Scope: write.
func (s *BucketsService) CreateShare(ctx context.Context, bucketID string, req CreateBucketShareRequest, opts ...RequestOption) (*BucketShare, error) {
	var out BucketShare
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/shares", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// RevokeShare revokes a cross-organization share on a bucket:
// DELETE /buckets/{bucket_id}/shares/{share_id}.
// Scope: delete.
func (s *BucketsService) RevokeShare(ctx context.Context, bucketID, shareID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", bucketPath(bucketID)+"/shares/"+url.PathEscape(shareID), nil, nil, nil, opts...)
}

// ListNotifications lists object event notification rules on a bucket:
// GET /buckets/{bucket_id}/notifications.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *BucketsService) ListNotifications(ctx context.Context, bucketID string, opts ...RequestOption) ([]BucketNotification, error) {
	return listAll[BucketNotification](ctx, s.c, bucketPath(bucketID)+"/notifications", nil, opts...)
}

// CreateNotification creates a bucket object event notification rule:
// POST /buckets/{bucket_id}/notifications.
// Scope: write.
func (s *BucketsService) CreateNotification(ctx context.Context, bucketID string, req CreateBucketNotificationRequest, opts ...RequestOption) (*BucketNotification, error) {
	var out BucketNotification
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/notifications", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateNotification updates a bucket object event notification rule:
// PATCH /buckets/{bucket_id}/notifications/{notification_id}.
// Scope: write.
func (s *BucketsService) UpdateNotification(ctx context.Context, bucketID, notificationID string, req UpdateBucketNotificationRequest, opts ...RequestOption) (*BucketNotification, error) {
	var out BucketNotification
	if err := s.c.Do(ctx, "PATCH", bucketPath(bucketID)+"/notifications/"+url.PathEscape(notificationID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteNotification deletes a bucket object event notification rule:
// DELETE /buckets/{bucket_id}/notifications/{notification_id}.
// Scope: delete.
func (s *BucketsService) DeleteNotification(ctx context.Context, bucketID, notificationID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", bucketPath(bucketID)+"/notifications/"+url.PathEscape(notificationID), nil, nil, nil, opts...)
}

// ListDomains lists custom domains attached to a bucket:
// GET /buckets/{bucket_id}/domains.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *BucketsService) ListDomains(ctx context.Context, bucketID string, opts ...RequestOption) ([]BucketDomain, error) {
	return listAll[BucketDomain](ctx, s.c, bucketPath(bucketID)+"/domains", nil, opts...)
}

// CreateDomain attaches a custom domain to a bucket:
// POST /buckets/{bucket_id}/domains.
// Scope: write.
func (s *BucketsService) CreateDomain(ctx context.Context, bucketID string, req CreateBucketDomainRequest, opts ...RequestOption) (*BucketDomain, error) {
	var out BucketDomain
	if err := s.c.Do(ctx, "POST", bucketPath(bucketID)+"/domains", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDomain updates a bucket's custom domain configuration:
// PATCH /buckets/{bucket_id}/domains/{domain_id}.
// Scope: write.
func (s *BucketsService) UpdateDomain(ctx context.Context, bucketID, domainID string, req UpdateBucketDomainRequest, opts ...RequestOption) (*BucketDomain, error) {
	var out BucketDomain
	if err := s.c.Do(ctx, "PATCH", bucketPath(bucketID)+"/domains/"+url.PathEscape(domainID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteDomain detaches a custom domain from a bucket:
// DELETE /buckets/{bucket_id}/domains/{domain_id}.
// Scope: delete.
func (s *BucketsService) DeleteDomain(ctx context.Context, bucketID, domainID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", bucketPath(bucketID)+"/domains/"+url.PathEscape(domainID), nil, nil, nil, opts...)
}
