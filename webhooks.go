package pantech

import (
	"context"
	"iter"
	"net/url"
)

// WebhooksService manages webhooks and their deliveries.
type WebhooksService struct{ c *Client }

// ListEventTypes lists webhook event types: GET /webhook-event-types.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *WebhooksService) ListEventTypes(ctx context.Context, opts ...RequestOption) ([]WebhookEventType, error) {
	return listAll[WebhookEventType](ctx, s.c, "/webhook-event-types", nil, opts...)
}

// List lists webhooks: GET /webhooks.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *WebhooksService) List(ctx context.Context, opts ...RequestOption) ([]Webhook, error) {
	return listAll[Webhook](ctx, s.c, "/webhooks", nil, opts...)
}

// Create creates a webhook: POST /webhooks.
// Scope: write.
func (s *WebhooksService) Create(ctx context.Context, req CreateWebhookRequest, opts ...RequestOption) (*WebhookWithSecret, error) {
	var out WebhookWithSecret
	if err := s.c.Do(ctx, "POST", "/webhooks", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a webhook: GET /webhooks/{webhook_id}.
// Scope: read.
func (s *WebhooksService) Get(ctx context.Context, webhookID string, opts ...RequestOption) (*Webhook, error) {
	var out Webhook
	if err := s.c.Do(ctx, "GET", "/webhooks/"+url.PathEscape(webhookID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update changes a webhook: PATCH /webhooks/{webhook_id}.
// Scope: write.
func (s *WebhooksService) Update(ctx context.Context, webhookID string, req UpdateWebhookRequest, opts ...RequestOption) (*Webhook, error) {
	var out Webhook
	if err := s.c.Do(ctx, "PATCH", "/webhooks/"+url.PathEscape(webhookID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a webhook: DELETE /webhooks/{webhook_id}.
// Scope: write.
func (s *WebhooksService) Delete(ctx context.Context, webhookID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", "/webhooks/"+url.PathEscape(webhookID), nil, nil, nil, opts...)
}

// Test sends a test event: POST /webhooks/{webhook_id}/test.
// Scope: write.
func (s *WebhooksService) Test(ctx context.Context, webhookID string, opts ...RequestOption) (*WebhookDeliveryDetail, error) {
	var out WebhookDeliveryDetail
	if err := s.c.Do(ctx, "POST", "/webhooks/"+url.PathEscape(webhookID)+"/test", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// RotateSecret rotates a webhook's secret: POST /webhooks/{webhook_id}/rotate-secret.
// Scope: write.
func (s *WebhooksService) RotateSecret(ctx context.Context, webhookID string, opts ...RequestOption) (*WebhookWithSecret, error) {
	var out WebhookWithSecret
	if err := s.c.Do(ctx, "POST", "/webhooks/"+url.PathEscape(webhookID)+"/rotate-secret", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListDeliveries lists a webhook's deliveries: GET /webhooks/{webhook_id}/deliveries.
// Scope: read.
func (s *WebhooksService) ListDeliveries(ctx context.Context, webhookID string, params *PageParams, opts ...RequestOption) (*Page[WebhookDelivery], error) {
	var out Page[WebhookDelivery]
	if err := s.c.Do(ctx, "GET", "/webhooks/"+url.PathEscape(webhookID)+"/deliveries", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllDeliveries iterates over every item of ListDeliveries, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *WebhooksService) AllDeliveries(ctx context.Context, webhookID string, params *PageParams, opts ...RequestOption) iter.Seq2[WebhookDelivery, error] {
	var p PageParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[WebhookDelivery], error) {
		p.Cursor = cursor
		return s.ListDeliveries(ctx, webhookID, &p, opts...)
	})
}

// GetDelivery gets a delivery: GET /webhooks/{webhook_id}/deliveries/{delivery_id}.
// Scope: read.
func (s *WebhooksService) GetDelivery(ctx context.Context, webhookID string, deliveryID string, opts ...RequestOption) (*WebhookDeliveryDetail, error) {
	var out WebhookDeliveryDetail
	if err := s.c.Do(ctx, "GET", "/webhooks/"+url.PathEscape(webhookID)+"/deliveries/"+url.PathEscape(deliveryID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Redeliver sends a delivery again: POST /webhooks/{webhook_id}/deliveries/{delivery_id}/redeliver.
// Scope: write.
func (s *WebhooksService) Redeliver(ctx context.Context, webhookID string, deliveryID string, opts ...RequestOption) (*WebhookDelivery, error) {
	var out WebhookDelivery
	if err := s.c.Do(ctx, "POST", "/webhooks/"+url.PathEscape(webhookID)+"/deliveries/"+url.PathEscape(deliveryID)+"/redeliver", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
