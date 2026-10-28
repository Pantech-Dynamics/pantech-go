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
