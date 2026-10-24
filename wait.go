package pantech

import (
	"context"
	"fmt"
	"time"
)

// Operation statuses. Succeeded and failed are final.
const (
	OperationSubmitting = "submitting"
	OperationSubmitted  = "submitted"
	OperationSucceeded  = "succeeded"
	OperationFailed     = "failed"
)

// Order statuses, shared by instance and database orders. Provisioned is
// success; payment_failed and failed are final failures, and any payment
// taken is returned to credit.
const (
	OrderAwaitingPayment = "awaiting_payment"
	OrderPaid            = "paid"
	OrderProvisioning    = "provisioning"
	OrderProvisioned     = "provisioned"
	OrderPaymentFailed   = "payment_failed"
	OrderFailed          = "failed"
)

// WaitOptions tune a wait. The zero value polls every 2 seconds with no time
// limit beyond ctx's.
type WaitOptions struct {
	// Interval is the pause between polls. Default 2s.
	Interval time.Duration
	// Timeout bounds the whole wait. 0 means only ctx bounds it.
	Timeout time.Duration
}

// OperationError is returned by a wait when the operation failed. Its
// Failure code is stable, and its reason is safe to show to users.
type OperationError struct {
	Operation *Operation
}

func (e *OperationError) Error() string {
	msg := fmt.Sprintf("pantech: operation %s (%s) failed", e.Operation.ID, e.Operation.Kind)
	if f := e.Operation.Failure; f != nil {
		msg += ": " + f.Reason + " (" + f.Code + ")"
	}
	return msg
}

// OrderError is returned by an order wait when the order ended in
// payment_failed or failed. Nothing was provisioned.
type OrderError struct {
	// Kind is "instance" or "database".
	Kind        string
	OrderID     string
	Status      string
	FailureCode string
}

func (e *OrderError) Error() string {
	msg := fmt.Sprintf("pantech: %s order %s ended as %s", e.Kind, e.OrderID, e.Status)
	if e.FailureCode != "" {
		msg += " (" + e.FailureCode + ")"
	}
	if e.Status == OrderPaymentFailed {
		msg += ": the payment did not go through and nothing was provisioned"
	} else {
		msg += ": nothing was provisioned and any payment taken is returned to your credit"
	}
	return msg
}

// Wait polls GET /operations/{id} until the operation succeeds, and returns
// it. A failed operation returns it with an *OperationError. ctx (and
// opts.Timeout) bound the wait; opts may be nil.
func (s *OperationsService) Wait(ctx context.Context, operationID string, opts *WaitOptions) (*Operation, error) {
	var last *Operation
	err := s.c.poll(ctx, opts, func(ctx context.Context) (bool, error) {
		op, err := s.Get(ctx, operationID)
		if err != nil {
			return false, err
		}
		last = op
		switch op.Status {
		case OperationSucceeded:
			return true, nil
		case OperationFailed:
			return true, &OperationError{Operation: op}
		}
		return false, nil
	})
	return last, err
}

// WaitForOrder polls GET /instance-orders/{id} until the order is
// provisioned, and returns it: the instance then exists (read it with Get,
// and wait for its OperationID, if any, to see it running). A failed or
// unpaid order returns an *OrderError.
func (s *InstancesService) WaitForOrder(ctx context.Context, orderID string, opts *WaitOptions) (*InstanceOrder, error) {
	var last *InstanceOrder
	err := s.c.poll(ctx, opts, func(ctx context.Context) (bool, error) {
		order, err := s.GetOrder(ctx, orderID)
		if err != nil {
			return false, err
		}
		last = order
		return orderDone("instance", order.ID, order.Status, order.FailureCode)
	})
	return last, err
}

// WaitForOrder polls GET /database-orders/{id} until the order is
// provisioned, and returns it. The database is then being created: wait for
// its OperationID with Operations.Wait to see it running. A failed or unpaid
// order returns an *OrderError.
func (s *DatabasesService) WaitForOrder(ctx context.Context, orderID string, opts *WaitOptions) (*DatabaseOrder, error) {
	var last *DatabaseOrder
	err := s.c.poll(ctx, opts, func(ctx context.Context) (bool, error) {
		order, err := s.GetOrder(ctx, orderID)
		if err != nil {
			return false, err
		}
		last = order
		return orderDone("database", order.ID, order.Status, order.FailureCode)
	})
	return last, err
}

func orderDone(kind, id, status string, failureCode *string) (bool, error) {
	switch status {
	case OrderProvisioned:
		return true, nil
	case OrderFailed, OrderPaymentFailed:
		e := &OrderError{Kind: kind, OrderID: id, Status: status}
		if failureCode != nil {
			e.FailureCode = *failureCode
		}
		return true, e
	}
	return false, nil
}

// poll runs check until it reports done or fails, sleeping between checks.
func (c *Client) poll(ctx context.Context, opts *WaitOptions, check func(context.Context) (bool, error)) error {
	interval := c.pollInterval
	if opts != nil && opts.Interval > 0 {
		interval = opts.Interval
	}
	if opts != nil && opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}
	for {
		done, err := check(ctx)
		if done || err != nil {
			return err
		}
		if err := c.sleep(ctx, interval); err != nil {
			return fmt.Errorf("pantech: waiting: %w", err)
		}
	}
}
