package pantech

import (
	"context"
	"net/url"
)

// CatalogService is what you can buy and what it costs: regions, plans, images, disk offerings, database engines, Kubernetes versions and price estimates.
type CatalogService struct{ c *Client }

// ListRegions lists regions: GET /regions.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *CatalogService) ListRegions(ctx context.Context, opts ...RequestOption) ([]Region, error) {
	return listAll[Region](ctx, s.c, "/regions", nil, opts...)
}

// ListPlans lists plans: GET /plans.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *CatalogService) ListPlans(ctx context.Context, params *ListPlansParams, opts ...RequestOption) ([]Plan, error) {
	return listAll[Plan](ctx, s.c, "/plans", params.values(), opts...)
}

// ListImages lists images: GET /images.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *CatalogService) ListImages(ctx context.Context, opts ...RequestOption) ([]Image, error) {
	return listAll[Image](ctx, s.c, "/images", nil, opts...)
}

// ListDiskOfferings lists disk offerings: GET /disk-offerings.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *CatalogService) ListDiskOfferings(ctx context.Context, opts ...RequestOption) ([]DiskOffering, error) {
	return listAll[DiskOffering](ctx, s.c, "/disk-offerings", nil, opts...)
}

// QuoteDiskOffering prices a volume: GET /disk-offerings/{slug}/quote.
// Scope: read.
func (s *CatalogService) QuoteDiskOffering(ctx context.Context, slug string, params *DiskOfferingQuoteParams, opts ...RequestOption) (*DiskOfferingQuote, error) {
	var out DiskOfferingQuote
	if err := s.c.Do(ctx, "GET", "/disk-offerings/"+url.PathEscape(slug)+"/quote", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListDatabaseEngines lists database engines: GET /database-engines.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *CatalogService) ListDatabaseEngines(ctx context.Context, opts ...RequestOption) ([]DatabaseEngine, error) {
	return listAll[DatabaseEngine](ctx, s.c, "/database-engines", nil, opts...)
}

// ListKubernetesVersions lists Kubernetes versions: GET /kubernetes-versions.
// Scope: read.
func (s *CatalogService) ListKubernetesVersions(ctx context.Context, params *ListKubernetesVersionsParams, opts ...RequestOption) (*KubernetesVersionList, error) {
	var out KubernetesVersionList
	if err := s.c.Do(ctx, "GET", "/kubernetes-versions", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// EstimatePrice estimates a monthly bill: POST /pricing/estimate.
// Scope: read.
func (s *CatalogService) EstimatePrice(ctx context.Context, req PricingEstimateRequest, opts ...RequestOption) (*PricingEstimate, error) {
	var out PricingEstimate
	if err := s.c.Do(ctx, "POST", "/pricing/estimate", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
