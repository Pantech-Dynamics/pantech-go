package pantech

import (
	"net/url"
	"strconv"
)

// PageParams are the query parameters of the paging of lists that take no filters. Zero values are left out.
type PageParams struct {
	// Items per page, 1 to 100.
	Limit int
	// The previous page's `next_cursor`, to get the next page.
	Cursor string
}

func (p *PageParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit != 0 {
		v.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		v.Set("cursor", p.Cursor)
	}
	return v
}

// UsageByResourceParams are the query parameters of usage by resource. Zero values are left out.
type UsageByResourceParams struct {
	// One of: 30d, 7d, month, last_month.
	Period string
}

func (p *UsageByResourceParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Period != "" {
		v.Set("period", p.Period)
	}
	return v
}

// ListPlansParams are the query parameters of list plans. Zero values are left out.
type ListPlansParams struct {
	// Price for standard instances or VPC instances. One of: standard, vpc.
	Placement string
}

func (p *ListPlansParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Placement != "" {
		v.Set("placement", p.Placement)
	}
	return v
}

// DiskOfferingQuoteParams are the query parameters of disk offering quote. Zero values are left out.
type DiskOfferingQuoteParams struct {
	SizeGB int
}

func (p *DiskOfferingQuoteParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.SizeGB != 0 {
		v.Set("size_gb", strconv.Itoa(p.SizeGB))
	}
	return v
}

// ListKubernetesVersionsParams are the query parameters of list kubernetes versions. Zero values are left out.
type ListKubernetesVersionsParams struct {
	// Only versions in this zone.
	ZoneID string
}

func (p *ListKubernetesVersionsParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.ZoneID != "" {
		v.Set("zone_id", p.ZoneID)
	}
	return v
}

// ListParams are the query parameters of the paging and filters most lists take. Zero values are left out.
type ListParams struct {
	// Items per page, 1 to 100.
	Limit int
	// The previous page's `next_cursor`, to get the next page.
	Cursor string
	// Only items whose name contains this text (case-insensitive), at most 255 characters.
	Q string
	// Only items with this `desired_state`: `running`, `stopped` or `deleted`. One of: running, stopped, deleted.
	DesiredState string
	// Only items with this `observed_state`. One of: pending, provisioning, running, stopping, stopped, deleting, deleted, failed, unknown.
	ObservedState string
}

func (p *ListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit != 0 {
		v.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		v.Set("cursor", p.Cursor)
	}
	if p.Q != "" {
		v.Set("q", p.Q)
	}
	if p.DesiredState != "" {
		v.Set("desired_state", p.DesiredState)
	}
	if p.ObservedState != "" {
		v.Set("observed_state", p.ObservedState)
	}
	return v
}

// InstanceMetricsParams are the query parameters of instance metrics. Zero values are left out.
type InstanceMetricsParams struct {
	// One of: 1h, 24h, 7d, 30d.
	Interval string
}

func (p *InstanceMetricsParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Interval != "" {
		v.Set("interval", p.Interval)
	}
	return v
}

// ListPublicIPsParams are the query parameters of list public ips. Zero values are left out.
type ListPublicIPsParams struct {
	NetworkID string
}

func (p *ListPublicIPsParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.NetworkID != "" {
		v.Set("network_id", p.NetworkID)
	}
	return v
}

// ListLoadBalancersParams are the query parameters of list load balancers. Zero values are left out.
type ListLoadBalancersParams struct {
	// Only load balancers in this private network (`net_…`).
	NetworkID string
	// Only load balancers on this public IP (`pip_…`).
	PublicIPID string
}

func (p *ListLoadBalancersParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.NetworkID != "" {
		v.Set("network_id", p.NetworkID)
	}
	if p.PublicIPID != "" {
		v.Set("public_ip_id", p.PublicIPID)
	}
	return v
}

// GetCloudScanParams are the query parameters of get cloud scan. Zero values are left out.
type GetCloudScanParams struct {
	// `inventory` returns everything the scan found. One of: inventory.
	Include string
}

func (p *GetCloudScanParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Include != "" {
		v.Set("include", p.Include)
	}
	return v
}

// GetServerCopyTemplateParams are the query parameters of get server copy template. Zero values are left out.
type GetServerCopyTemplateParams struct {
	// The planned server's item key. Required.
	ItemKey string
	// The copy role may also snapshot that server's own volumes and delete the snapshots it made.
	CreateSnapshots bool
}

func (p *GetServerCopyTemplateParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.ItemKey != "" {
		v.Set("item_key", p.ItemKey)
	}
	if p.CreateSnapshots {
		v.Set("create_snapshots", "true")
	}
	return v
}

// ListDNSZonesParams are the query parameters of list DNS zones. Zero values are left out.
type ListDNSZonesParams struct {
	// Items per page, 1 to 100.
	Limit int
	// The previous page's `next_cursor`, to get the next page.
	Cursor string
	// Only zones whose name contains this text.
	Q string
	// Only zones with this status. One of: pending, awaiting_delegation, active, suspended, deleting, failed.
	Status string
}

func (p *ListDNSZonesParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit != 0 {
		v.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		v.Set("cursor", p.Cursor)
	}
	if p.Q != "" {
		v.Set("q", p.Q)
	}
	if p.Status != "" {
		v.Set("status", p.Status)
	}
	return v
}

// ListDNSRecordsParams are the query parameters of list DNS records. Zero values are left out.
type ListDNSRecordsParams struct {
	// Items per page, 1 to 100.
	Limit int
	// The previous page's `next_cursor`, to get the next page.
	Cursor string
	// Only records with exactly this name, relative to the zone (`@`, `www`).
	Name string
	// Only records whose name or values contain this text.
	Q string
	// Only records of this type, such as `MX`.
	Type string
}

func (p *ListDNSRecordsParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit != 0 {
		v.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		v.Set("cursor", p.Cursor)
	}
	if p.Name != "" {
		v.Set("name", p.Name)
	}
	if p.Q != "" {
		v.Set("q", p.Q)
	}
	if p.Type != "" {
		v.Set("type", p.Type)
	}
	return v
}

// ListDomainsParams are the query parameters of list domains. Zero values are left out.
type ListDomainsParams struct {
	// Items per page, 1 to 100.
	Limit int
	// The previous page's `next_cursor`, to get the next page.
	Cursor string
	// Only domains whose name contains this text.
	Q string
	// Only domains with this status. One of: pending, registering, transferring, active, expired, failed, transferred_out, suspended.
	Status string
}

func (p *ListDomainsParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit != 0 {
		v.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		v.Set("cursor", p.Cursor)
	}
	if p.Q != "" {
		v.Set("q", p.Q)
	}
	if p.Status != "" {
		v.Set("status", p.Status)
	}
	return v
}

// ListRegistryParams are the query parameters of list registry repositories and list registry tags. Zero values are left out.
type ListRegistryParams struct {
	// Items per page, 1 to 100.
	Limit int
	// The previous page's `next_cursor`, to get the next page.
	Cursor string
	// Only items whose name contains this text.
	Q string
}

func (p *ListRegistryParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit != 0 {
		v.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		v.Set("cursor", p.Cursor)
	}
	if p.Q != "" {
		v.Set("q", p.Q)
	}
	return v
}

// ListIAMPoliciesParams are the query parameters of list policies. Zero values are left out.
type ListIAMPoliciesParams struct {
	// Items per page, 1 to 100.
	Limit int
	// The previous page's `next_cursor`, to get the next page.
	Cursor string
	// Only Pantech's managed policies (true) or only yours (false).
	Managed *bool
}

func (p *ListIAMPoliciesParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit != 0 {
		v.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		v.Set("cursor", p.Cursor)
	}
	if p.Managed != nil {
		v.Set("managed", strconv.FormatBool(*p.Managed))
	}
	return v
}

// ListIAMMembersParams are the query parameters of list members. Zero values are left out.
type ListIAMMembersParams struct {
	// Items per page, 1 to 100.
	Limit int
	// The previous page's `next_cursor`, to get the next page.
	Cursor string
	// Only members whose email or name contains this text.
	Q string
}

func (p *ListIAMMembersParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit != 0 {
		v.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		v.Set("cursor", p.Cursor)
	}
	if p.Q != "" {
		v.Set("q", p.Q)
	}
	return v
}

// ListCertificatesParams are the query parameters of list certificates. Zero values are left out.
type ListCertificatesParams struct {
	// Items per page, 1 to 100.
	Limit int
	// The previous page's `next_cursor`, to get the next page.
	Cursor string
	// Only certificates with a name containing this text.
	Q string
	// One of: pending_validation, issuing, issued, renewing, expired, revoked, failed.
	Status string
	// One of: public, imported.
	Kind string
}

func (p *ListCertificatesParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Limit != 0 {
		v.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		v.Set("cursor", p.Cursor)
	}
	if p.Q != "" {
		v.Set("q", p.Q)
	}
	if p.Status != "" {
		v.Set("status", p.Status)
	}
	if p.Kind != "" {
		v.Set("kind", p.Kind)
	}
	return v
}

// ObjectStorageUsageParams are the query parameters of object storage usage. Zero values are left out.
type ObjectStorageUsageParams struct {
	// One of: 24h, 7d, 30d.
	Range string
}

func (p *ObjectStorageUsageParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Range != "" {
		v.Set("range", p.Range)
	}
	return v
}

// ListObjectStorageKeysParams are the query parameters of list object storage keys. Zero values are left out.
type ListObjectStorageKeysParams struct {
	// Region code, e.g. `af-abj`. Defaults to the region object storage is offered in.
	Region string
}

func (p *ListObjectStorageKeysParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Region != "" {
		v.Set("region", p.Region)
	}
	return v
}

// DeleteBucketParams are the query parameters of delete bucket. Zero values are left out.
type DeleteBucketParams struct {
	// Delete the bucket even though it holds objects, which are removed with it.
	Force bool
}

func (p *DeleteBucketParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Force {
		v.Set("force", "true")
	}
	return v
}

// ListBucketObjectsParams are the query parameters of list bucket objects. Zero values are left out.
type ListBucketObjectsParams struct {
	// Limit results to keys starting with this prefix.
	Prefix string
	// Folder delimiter (defaults to `/` when nil; pass Ptr("") for a flat listing).
	Delimiter *string
	// Search query across keys in the bucket.
	Q string
	// Pagination cursor returned as `next_cursor`.
	Cursor string
	// Maximum entries per page (1 to 1000).
	Limit int
	// When true, list object versions and delete markers instead of current objects.
	Versions bool
	// Key marker for versioned pagination (`next_key_marker`).
	KeyMarker string
	// Version ID marker for versioned pagination (`next_version_id_marker`).
	VersionIDMarker string
}

func (p *ListBucketObjectsParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Prefix != "" {
		v.Set("prefix", p.Prefix)
	}
	if p.Delimiter != nil {
		v.Set("delimiter", *p.Delimiter)
	}
	if p.Q != "" {
		v.Set("q", p.Q)
	}
	if p.Cursor != "" {
		v.Set("cursor", p.Cursor)
	}
	if p.Limit != 0 {
		v.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Versions {
		v.Set("versions", "true")
	}
	if p.KeyMarker != "" {
		v.Set("key_marker", p.KeyMarker)
	}
	if p.VersionIDMarker != "" {
		v.Set("version_id_marker", p.VersionIDMarker)
	}
	return v
}

// BucketObjectMetaParams are the query parameters of get bucket object metadata.
type BucketObjectMetaParams struct {
	Key       string
	VersionID string
}

func (p *BucketObjectMetaParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Key != "" {
		v.Set("key", p.Key)
	}
	if p.VersionID != "" {
		v.Set("version_id", p.VersionID)
	}
	return v
}

// DownloadBucketObjectParams are the query parameters of download bucket object.
type DownloadBucketObjectParams struct {
	Key       string
	VersionID string
	Inline    bool
}

func (p *DownloadBucketObjectParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Key != "" {
		v.Set("key", p.Key)
	}
	if p.VersionID != "" {
		v.Set("version_id", p.VersionID)
	}
	if p.Inline {
		v.Set("inline", "true")
	}
	return v
}

// UploadBucketObjectParams are the query parameters and optional object headers of upload bucket object.
type UploadBucketObjectParams struct {
	Key                string
	ContentType        string
	CacheControl       string
	ContentDisposition string
	Public             bool
	Metadata           map[string]string
}

func (p *UploadBucketObjectParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Key != "" {
		v.Set("key", p.Key)
	}
	return v
}

func (p *UploadBucketObjectParams) requestOptions() []RequestOption {
	if p == nil {
		return nil
	}
	var opts []RequestOption
	if p.ContentType != "" {
		opts = append(opts, WithRequestHeader("Content-Type", p.ContentType))
	}
	if p.CacheControl != "" {
		opts = append(opts, WithRequestHeader("X-Object-Cache-Control", p.CacheControl))
	}
	if p.ContentDisposition != "" {
		opts = append(opts, WithRequestHeader("X-Object-Content-Disposition", p.ContentDisposition))
	}
	if p.Public {
		opts = append(opts, WithRequestHeader("X-Object-Public", "true"))
	}
	for k, v := range p.Metadata {
		opts = append(opts, WithRequestHeader("X-Object-Meta-"+k, v))
	}
	return opts
}

// UploadMultipartPartParams are the query parameters of upload multipart part.
type UploadMultipartPartParams struct {
	Key        string
	UploadID   string
	PartNumber int
}

func (p *UploadMultipartPartParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Key != "" {
		v.Set("key", p.Key)
	}
	if p.UploadID != "" {
		v.Set("upload_id", p.UploadID)
	}
	if p.PartNumber != 0 {
		v.Set("part_number", strconv.Itoa(p.PartNumber))
	}
	return v
}
