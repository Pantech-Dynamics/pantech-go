package pantech

import (
	"context"
	"iter"
	"net/url"
)

// BillingService is credit, receipts and usage.
type BillingService struct{ c *Client }

// GetCredit gets credit: GET /billing/credit.
// Scope: read.
func (s *BillingService) GetCredit(ctx context.Context, opts ...RequestOption) (*CreditOutlook, error) {
	var out CreditOutlook
	if err := s.c.Do(ctx, "GET", "/billing/credit", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCreditPurchase buys credit: POST /billing/credit/purchases.
// Scope: write.
func (s *BillingService) CreateCreditPurchase(ctx context.Context, req CreateCreditPurchaseRequest, opts ...RequestOption) (*CreditPurchase, error) {
	var out CreditPurchase
	if err := s.c.Do(ctx, "POST", "/billing/credit/purchases", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCreditPurchase gets a credit purchase: GET /billing/credit/purchases/{purchase_id}.
// Scope: read.
func (s *BillingService) GetCreditPurchase(ctx context.Context, purchaseID string, opts ...RequestOption) (*CreditPurchase, error) {
	var out CreditPurchase
	if err := s.c.Do(ctx, "GET", "/billing/credit/purchases/"+url.PathEscape(purchaseID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListReceipts lists receipts: GET /billing/receipts.
// Scope: read.
func (s *BillingService) ListReceipts(ctx context.Context, params *PageParams, opts ...RequestOption) (*Page[Receipt], error) {
	var out Page[Receipt]
	if err := s.c.Do(ctx, "GET", "/billing/receipts", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllReceipts iterates over every item of ListReceipts, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *BillingService) AllReceipts(ctx context.Context, params *PageParams, opts ...RequestOption) iter.Seq2[Receipt, error] {
	var p PageParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Receipt], error) {
		p.Cursor = cursor
		return s.ListReceipts(ctx, &p, opts...)
	})
}

// GetReceipt gets a receipt: GET /billing/receipts/{receipt_id}.
// Scope: read.
func (s *BillingService) GetReceipt(ctx context.Context, receiptID string, opts ...RequestOption) (*Receipt, error) {
	var out Receipt
	if err := s.c.Do(ctx, "GET", "/billing/receipts/"+url.PathEscape(receiptID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetReceiptPDF downloads a receipt: GET /billing/receipts/{receipt_id}/pdf.
// Scope: read.
func (s *BillingService) GetReceiptPDF(ctx context.Context, receiptID string, opts ...RequestOption) ([]byte, error) {
	var out []byte
	if err := s.c.Do(ctx, "GET", "/billing/receipts/"+url.PathEscape(receiptID)+"/pdf", nil, nil, &out, append([]RequestOption{WithRequestHeader("Accept", "application/pdf")}, opts...)...); err != nil {
		return nil, err
	}
	return out, nil
}

// GetUsage gets usage: GET /usage.
// Scope: read.
func (s *BillingService) GetUsage(ctx context.Context, opts ...RequestOption) (*UsageSummary, error) {
	var out UsageSummary
	if err := s.c.Do(ctx, "GET", "/usage", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUsageByResource gets cost by resource: GET /usage/by-resource.
// Scope: read.
func (s *BillingService) GetUsageByResource(ctx context.Context, params *UsageByResourceParams, opts ...RequestOption) (*UsageByResource, error) {
	var out UsageByResource
	if err := s.c.Do(ctx, "GET", "/usage/by-resource", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
