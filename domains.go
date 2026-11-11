package pantech

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"net/url"
	"time"
)

// Domain statuses.
const (
	DomainPending        = "pending"
	DomainRegistering    = "registering"
	DomainTransferring   = "transferring"
	DomainActive         = "active"
	DomainExpired        = "expired"
	DomainFailed         = "failed"
	DomainTransferredOut = "transferred_out"
	DomainSuspended      = "suspended"
)

// Domain order statuses: pending, then paid, then submitted, then succeeded
// or failed. A pending order whose credit falls short becomes
// awaiting_payment for up to an hour: in US dollars the saved card is
// charged the difference, in naira the customer adds credit, and
// RetryOrderPayment tries again at once. Succeeded and failed are final; a
// failed order's amount is returned to credit.
const (
	DomainOrderPending         = "pending"
	DomainOrderAwaitingPayment = "awaiting_payment"
	DomainOrderPaid            = "paid"
	DomainOrderSubmitted       = "submitted"
	DomainOrderSucceeded       = "succeeded"
	DomainOrderFailed          = "failed"
)

// Domain order kinds.
const (
	DomainOrderRegister = "register"
	DomainOrderRenew    = "renew"
	DomainOrderTransfer = "transfer"
)

// Domain search result statuses. Unknown means availability could not be
// checked now; try again.
const (
	DomainAvailable   = "available"
	DomainUnavailable = "unavailable"
	DomainUnknown     = "unknown"
)

// DomainsService searches for, registers, transfers in, renews and
// configures domains. Registrations, transfers and renewals are paid from the
// organization's credit before the registrar is asked.
type DomainsService struct{ c *Client }

// ListTLDs lists the TLDs Pantech sells to you, priced per year in your
// currency: GET /domains/tlds.
// Scope: read.
func (s *DomainsService) ListTLDs(ctx context.Context, opts ...RequestOption) ([]TLDOffer, error) {
	return listAll[TLDOffer](ctx, s.c, "/domains/tlds", nil, opts...)
}

// Search answers whether a name is available, across the TLDs sold (or the
// exact name, when query has a TLD that is), plus suggestions, each with its
// price: GET /domains/search.
// Scope: read.
func (s *DomainsService) Search(ctx context.Context, query string, opts ...RequestOption) ([]DomainSearchResult, error) {
	var out struct {
		Data []DomainSearchResult `json:"data"`
	}
	if err := s.c.Do(ctx, "GET", "/domains/search", url.Values{"query": {query}}, nil, &out, opts...); err != nil {
		return nil, err
	}
	return nonNil(out.Data), nil
}

// CheckTransfer answers whether a domain can be transferred to Pantech now,
// and at what price: GET /domains/transfer-check.
// Scope: read.
func (s *DomainsService) CheckTransfer(ctx context.Context, name string, opts ...RequestOption) (*DomainTransferCheck, error) {
	var out DomainTransferCheck
	if err := s.c.Do(ctx, "GET", "/domains/transfer-check", url.Values{"name": {name}}, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// List lists domains: GET /domains.
// Scope: read.
func (s *DomainsService) List(ctx context.Context, params *ListDomainsParams, opts ...RequestOption) (*Page[Domain], error) {
	var out Page[Domain]
	if err := s.c.Do(ctx, "GET", "/domains", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// All iterates over every item of List, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *DomainsService) All(ctx context.Context, params *ListDomainsParams, opts ...RequestOption) iter.Seq2[Domain, error] {
	var p ListDomainsParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Domain], error) {
		p.Cursor = cursor
		return s.List(ctx, &p, opts...)
	})
}

// Register buys a domain: POST /domains. The order's amount
// (TLDOffer.RegisterPrice, plus VAT) is taken from credit before the
// registrar is asked; when credit falls short the order waits in
// awaiting_payment for an hour (see DomainOrder). Follow the order with WaitForOrder. The domain gets a Pantech DNS
// zone of its name. AcceptTerms must be true. Retrying with the same
// Idempotency-Key (WithIdempotencyKey) answers the original order.
// Scope: write.
func (s *DomainsService) Register(ctx context.Context, req RegisterDomainRequest, opts ...RequestOption) (*DomainOrderAccepted, error) {
	var out DomainOrderAccepted
	if err := s.c.Do(ctx, "POST", "/domains", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Transfer moves a domain to Pantech from another registrar with its auth
// code: POST /domains/transfers. It is paid like Register
// (TLDOffer.TransferPrice) and adds a year; follow the order with
// WaitForOrder. A transfer can take several days.
// Scope: write.
func (s *DomainsService) Transfer(ctx context.Context, req TransferDomainRequest, opts ...RequestOption) (*DomainOrderAccepted, error) {
	var out DomainOrderAccepted
	if err := s.c.Do(ctx, "POST", "/domains/transfers", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a domain: GET /domains/{domain_id}.
// Scope: read.
func (s *DomainsService) Get(ctx context.Context, domainID string, opts ...RequestOption) (*Domain, error) {
	var out Domain
	if err := s.c.Do(ctx, "GET", domainPath(domainID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update changes auto-renew, the transfer lock or privacy: PATCH /domains/{domain_id}.
// Lock and privacy are applied at the registrar while SettingsPending is true.
// Scope: write.
func (s *DomainsService) Update(ctx context.Context, domainID string, req UpdateDomainRequest, opts ...RequestOption) (*Domain, error) {
	var out Domain
	if err := s.c.Do(ctx, "PATCH", domainPath(domainID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetNameServers points the domain at Pantech DNS or at custom name servers:
// PUT /domains/{domain_id}/name-servers.
// Scope: write.
func (s *DomainsService) SetNameServers(ctx context.Context, domainID string, req SetDomainNameServersRequest, opts ...RequestOption) (*Domain, error) {
	var out Domain
	if err := s.c.Do(ctx, "PUT", domainPath(domainID)+"/name-servers", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetContact gets the registrant contact: GET /domains/{domain_id}/contact.
// It is personal data, for owners and admins only.
// Scope: write.
func (s *DomainsService) GetContact(ctx context.Context, domainID string, opts ...RequestOption) (*DomainContact, error) {
	var out DomainContact
	if err := s.c.Do(ctx, "GET", domainPath(domainID)+"/contact", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateContact replaces the registrant contact: PUT /domains/{domain_id}/contact.
// Some registries email the registrant to confirm a change.
// Scope: write.
func (s *DomainsService) UpdateContact(ctx context.Context, domainID string, req UpdateDomainContactRequest, opts ...RequestOption) (*Domain, error) {
	var out Domain
	if err := s.c.Do(ctx, "PUT", domainPath(domainID)+"/contact", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Renew extends a domain, paid from credit (TLDOffer.RenewPrice): POST /domains/{domain_id}/renew.
// Follow the order with WaitForOrder.
// Scope: write.
func (s *DomainsService) Renew(ctx context.Context, domainID string, req RenewDomainRequest, opts ...RequestOption) (*DomainOrderAccepted, error) {
	var out DomainOrderAccepted
	if err := s.c.Do(ctx, "POST", domainPath(domainID)+"/renew", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// RequestAuthCode asks for the code that moves the domain to another
// registrar: POST /domains/{domain_id}/auth-code. The transfer lock must be
// off. Read the code with WaitForAuthCode. Anyone holding the code can take
// the domain away from Pantech.
// Scope: delete.
func (s *DomainsService) RequestAuthCode(ctx context.Context, domainID string, opts ...RequestOption) (*Domain, error) {
	var out Domain
	if err := s.c.Do(ctx, "POST", domainPath(domainID)+"/auth-code", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAuthCode reads a requested auth code: GET /domains/{domain_id}/auth-code.
// While it is being fetched it answers 202 and Pending reports true.
// Scope: delete.
func (s *DomainsService) GetAuthCode(ctx context.Context, domainID string, opts ...RequestOption) (*DomainAuthCode, error) {
	var out DomainAuthCode
	if err := s.c.Do(ctx, "GET", domainPath(domainID)+"/auth-code", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListOrders lists a domain's orders: GET /domains/{domain_id}/orders.
// Scope: read.
// It follows next_cursor until the list ends, so the slice holds every item.
func (s *DomainsService) ListOrders(ctx context.Context, domainID string, opts ...RequestOption) ([]DomainOrder, error) {
	return listAll[DomainOrder](ctx, s.c, domainPath(domainID)+"/orders", nil, opts...)
}

// GetOrder gets a domain order: GET /domain-orders/{order_id}.
// Scope: read.
func (s *DomainsService) GetOrder(ctx context.Context, orderID string, opts ...RequestOption) (*DomainOrder, error) {
	var out DomainOrder
	if err := s.c.Do(ctx, "GET", "/domain-orders/"+url.PathEscape(orderID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// RetryOrderPayment makes an order that waits in awaiting_payment try its
// payment again now: POST /domain-orders/{order_id}/retry-payment. In US
// dollars the saved card is charged the difference once more (after a
// decline, or once a new card is saved); in naira the credit just added is
// taken. The order then waits another hour. The call itself charges nothing
// and an order is never charged twice, so repeating it (with or without the
// same Idempotency-Key) is safe. It answers 202 with the order; follow it
// with WaitForPayment. An order that is not waiting for a payment, or whose
// window has ended, returns 409 CodeDomainOrderNotAwaitingPayment.
// Scope: write.
func (s *DomainsService) RetryOrderPayment(ctx context.Context, orderID string, opts ...RequestOption) (*DomainOrder, error) {
	var out DomainOrder
	if err := s.c.Do(ctx, "POST", "/domain-orders/"+url.PathEscape(orderID)+"/retry-payment", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DomainOrderError is returned by WaitForOrder when the order failed.
// FailureCode is stable, such as insufficient_credit or domain_unavailable.
type DomainOrderError struct {
	OrderID     string
	Kind        string
	DomainName  string
	FailureCode string
	// Refunded reports whether the amount went back to credit.
	Refunded bool
}

func (e *DomainOrderError) Error() string {
	msg := fmt.Sprintf("pantech: domain %s order %s for %s failed", e.Kind, e.OrderID, e.DomainName)
	if e.FailureCode != "" {
		msg += " (" + e.FailureCode + ")"
	}
	if e.Refunded {
		msg += ": the amount was returned to your credit"
	}
	return msg
}

// ErrPaymentRequired is matched (with errors.Is) by the
// *DomainPaymentRequiredError WaitForOrder returns for an order that waits
// for a payment.
var ErrPaymentRequired = errors.New("pantech: the order waits for a payment")

// DomainPaymentRequiredError is returned by WaitForOrder when the order's
// credit fell short and it waits in awaiting_payment. In US dollars the
// saved card is charged the difference; PaymentURL is set only when the bank
// asks the user to confirm, so send them there when it is. In naira the user
// adds credit. Then follow the order with WaitForPayment, or call
// RetryOrderPayment to try again at once. The order fails
// insufficient_credit, with nothing charged, if it is not paid by
// PaymentExpiresAt.
type DomainPaymentRequiredError struct {
	Order            *DomainOrder
	PaymentURL       string
	PaymentExpiresAt *time.Time
}

func (e *DomainPaymentRequiredError) Error() string {
	msg := fmt.Sprintf("pantech: domain %s order %s for %s waits for a payment", e.Order.Kind, e.Order.ID, e.Order.DomainName)
	if e.PaymentURL != "" {
		msg += ": confirm it at " + e.PaymentURL
	}
	if e.PaymentExpiresAt != nil {
		msg += " before " + e.PaymentExpiresAt.Format(time.RFC3339)
	}
	return msg
}

// Is makes errors.Is(err, ErrPaymentRequired) true.
func (e *DomainPaymentRequiredError) Is(target error) bool { return target == ErrPaymentRequired }

// WaitForOrder polls GET /domain-orders/{id} until the order succeeds, and
// returns it. A failed order returns it with a *DomainOrderError. An order
// that waits for a payment returns at once with a
// *DomainPaymentRequiredError (carrying the confirmation link when the bank
// asks for one), since only the user can move it on; follow it with
// WaitForPayment.
// Registrations usually finish within minutes; a transfer can take days, so
// bound that wait with opts.Timeout or ctx. opts may be nil.
func (s *DomainsService) WaitForOrder(ctx context.Context, orderID string, opts *WaitOptions) (*DomainOrder, error) {
	return s.waitOrder(ctx, orderID, opts, false)
}

// WaitForPayment is WaitForOrder for an order the user is paying for: it
// keeps polling while the order is awaiting_payment, until it succeeds or
// fails (insufficient_credit when the payment window ends unpaid). Bound it
// with opts.Timeout or ctx. opts may be nil.
func (s *DomainsService) WaitForPayment(ctx context.Context, orderID string, opts *WaitOptions) (*DomainOrder, error) {
	return s.waitOrder(ctx, orderID, opts, true)
}

func (s *DomainsService) waitOrder(ctx context.Context, orderID string, opts *WaitOptions, throughPayment bool) (*DomainOrder, error) {
	var last *DomainOrder
	err := s.c.poll(ctx, opts, func(ctx context.Context) (bool, error) {
		order, err := s.GetOrder(ctx, orderID)
		if err != nil {
			return false, err
		}
		last = order
		switch order.Status {
		case DomainOrderSucceeded:
			return true, nil
		case DomainOrderFailed:
			return true, &DomainOrderError{OrderID: order.ID, Kind: order.Kind, DomainName: order.DomainName, FailureCode: deref(order.FailureCode), Refunded: order.Refunded}
		case DomainOrderAwaitingPayment:
			if !throughPayment {
				return true, &DomainPaymentRequiredError{Order: order, PaymentURL: deref(order.PaymentURL), PaymentExpiresAt: order.PaymentExpiresAt}
			}
		}
		return false, nil
	})
	return last, err
}

// WaitForAuthCode polls GET /domains/{id}/auth-code until the code requested
// with RequestAuthCode is ready, and returns it. opts may be nil.
func (s *DomainsService) WaitForAuthCode(ctx context.Context, domainID string, opts *WaitOptions) (*DomainAuthCode, error) {
	var last *DomainAuthCode
	err := s.c.poll(ctx, opts, func(ctx context.Context) (bool, error) {
		code, err := s.GetAuthCode(ctx, domainID)
		if err != nil {
			return false, err
		}
		last = code
		return !code.Pending(), nil
	})
	return last, err
}

func domainPath(domainID string) string { return "/domains/" + url.PathEscape(domainID) }
