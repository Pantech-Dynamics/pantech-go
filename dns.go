package pantech

import (
	"context"
	"fmt"
	"iter"
	"net/url"
)

// DNS zone statuses. A zone in awaiting_delegation becomes active once the
// domain's registrar delegates to exactly its NameServers; after 14 days
// without that it is removed (failed, delegation_timeout).
const (
	DNSZonePending            = "pending"
	DNSZoneAwaitingDelegation = "awaiting_delegation"
	DNSZoneActive             = "active"
	DNSZoneSuspended          = "suspended"
	DNSZoneDeleting           = "deleting"
	DNSZoneDeleted            = "deleted"
	DNSZoneFailed             = "failed"
)

// DNS record statuses. A write is pending until Pantech DNS serves it,
// usually within a minute.
const (
	DNSRecordPending  = "pending"
	DNSRecordActive   = "active"
	DNSRecordDeleting = "deleting"
	DNSRecordFailed   = "failed"
)

// DNS record types. SOA is only ever listed: the apex SOA and NS are managed
// by Pantech.
const (
	DNSRecordA     = "A"
	DNSRecordAAAA  = "AAAA"
	DNSRecordCNAME = "CNAME"
	DNSRecordMX    = "MX"
	DNSRecordTXT   = "TXT"
	DNSRecordSRV   = "SRV"
	DNSRecordCAA   = "CAA"
	DNSRecordNS    = "NS"
	DNSRecordPTR   = "PTR"
	DNSRecordSOA   = "SOA"
)

// Actions of a DNSRecordChange.
const (
	DNSChangeCreate = "create"
	DNSChangeUpdate = "update"
	DNSChangeDelete = "delete"
)

// DNSService manages Pantech DNS, on Pantech's own name servers: hosted zones
// and their record sets. Writes answer at once with the resource pending; in
// an active zone it is served within about a minute (WaitForRecord). Records
// of a zone still awaiting delegation stay pending until the zone is active.
type DNSService struct{ c *Client }

// ListZones lists DNS zones: GET /dns/zones.
// Scope: read.
func (s *DNSService) ListZones(ctx context.Context, params *ListDNSZonesParams, opts ...RequestOption) (*Page[DNSZone], error) {
	var out Page[DNSZone]
	if err := s.c.Do(ctx, "GET", "/dns/zones", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllZones iterates over every item of ListZones, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *DNSService) AllZones(ctx context.Context, params *ListDNSZonesParams, opts ...RequestOption) iter.Seq2[DNSZone, error] {
	var p ListDNSZonesParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[DNSZone], error) {
		p.Cursor = cursor
		return s.ListZones(ctx, &p, opts...)
	})
}

// QuoteZones says what hosted zones cost: GET /dns/zones/quote. The first
// FreeZones live zones are free; each further one is billed by the hour.
// Scope: read.
func (s *DNSService) QuoteZones(ctx context.Context, opts ...RequestOption) (*DNSZoneQuote, error) {
	var out DNSZoneQuote
	if err := s.c.Do(ctx, "GET", "/dns/zones/quote", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateZone creates a hosted zone: POST /dns/zones. It is pending, then
// awaiting_delegation with the NameServers to set at the domain's registrar.
// Scope: write.
func (s *DNSService) CreateZone(ctx context.Context, req CreateDNSZoneRequest, opts ...RequestOption) (*DNSZone, error) {
	var out DNSZone
	if err := s.c.Do(ctx, "POST", "/dns/zones", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetZone gets a zone: GET /dns/zones/{zone_id}.
// Scope: read.
func (s *DNSService) GetZone(ctx context.Context, zoneID string, opts ...RequestOption) (*DNSZone, error) {
	var out DNSZone
	if err := s.c.Do(ctx, "GET", zonePath(zoneID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateZone changes a zone's description: PATCH /dns/zones/{zone_id}.
// Scope: write.
func (s *DNSService) UpdateZone(ctx context.Context, zoneID string, req UpdateDNSZoneRequest, opts ...RequestOption) (*DNSZone, error) {
	var out DNSZone
	if err := s.c.Do(ctx, "PATCH", zonePath(zoneID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteZone deletes a zone and its records: DELETE /dns/zones/{zone_id}.
// It stops answering at once and is returned deleting. A zone a Pantech
// domain delegates to answers DNS_ZONE_IN_USE_BY_DOMAIN.
// Scope: delete.
func (s *DNSService) DeleteZone(ctx context.Context, zoneID string, opts ...RequestOption) (*DNSZone, error) {
	var out DNSZone
	if err := s.c.Do(ctx, "DELETE", zonePath(zoneID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// VerifyZone checks the zone's delegation now instead of at the next
// scheduled check: POST /dns/zones/{zone_id}/verify. The result shows on the
// zone within seconds; follow it with WaitForZone.
// Scope: write.
func (s *DNSService) VerifyZone(ctx context.Context, zoneID string, opts ...RequestOption) (*DNSZone, error) {
	var out DNSZone
	if err := s.c.Do(ctx, "POST", zonePath(zoneID)+"/verify", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListRecords lists a zone's record sets, by name then type: GET /dns/zones/{zone_id}/records.
// Scope: read.
func (s *DNSService) ListRecords(ctx context.Context, zoneID string, params *ListDNSRecordsParams, opts ...RequestOption) (*Page[DNSRecord], error) {
	var out Page[DNSRecord]
	if err := s.c.Do(ctx, "GET", zonePath(zoneID)+"/records", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// AllRecords iterates over every item of ListRecords, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *DNSService) AllRecords(ctx context.Context, zoneID string, params *ListDNSRecordsParams, opts ...RequestOption) iter.Seq2[DNSRecord, error] {
	var p ListDNSRecordsParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[DNSRecord], error) {
		p.Cursor = cursor
		return s.ListRecords(ctx, zoneID, &p, opts...)
	})
}

// CreateRecord adds a record set: POST /dns/zones/{zone_id}/records. It is
// returned pending; follow it with WaitForRecord.
// Scope: write.
func (s *DNSService) CreateRecord(ctx context.Context, zoneID string, req CreateDNSRecordRequest, opts ...RequestOption) (*DNSRecord, error) {
	var out DNSRecord
	if err := s.c.Do(ctx, "POST", zonePath(zoneID)+"/records", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRecord gets a record set: GET /dns/zones/{zone_id}/records/{record_id}.
// Scope: read.
func (s *DNSService) GetRecord(ctx context.Context, zoneID, recordID string, opts ...RequestOption) (*DNSRecord, error) {
	var out DNSRecord
	if err := s.c.Do(ctx, "GET", recordPath(zoneID, recordID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateRecord changes a record set's TTL, values, allow_private or comment:
// PATCH /dns/zones/{zone_id}/records/{record_id}. Pass WithIfMatch(record.Version)
// to change it only if nobody else has since (412 VERSION_CONFLICT otherwise).
// Scope: write.
func (s *DNSService) UpdateRecord(ctx context.Context, zoneID, recordID string, req UpdateDNSRecordRequest, opts ...RequestOption) (*DNSRecord, error) {
	var out DNSRecord
	if err := s.c.Do(ctx, "PATCH", recordPath(zoneID, recordID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRecord deletes a record set: DELETE /dns/zones/{zone_id}/records/{record_id}.
// It is returned deleting. WithIfMatch applies as for UpdateRecord.
// Scope: delete.
func (s *DNSService) DeleteRecord(ctx context.Context, zoneID, recordID string, opts ...RequestOption) (*DNSRecord, error) {
	var out DNSRecord
	if err := s.c.Do(ctx, "DELETE", recordPath(zoneID, recordID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// BatchRecords applies up to 1000 creates, updates and deletes all or
// nothing: POST /dns/zones/{zone_id}/records/batch. It returns the records
// changed, in the order of the changes.
// Scope: write, and delete when a change deletes.
func (s *DNSService) BatchRecords(ctx context.Context, zoneID string, req DNSRecordBatchRequest, opts ...RequestOption) ([]DNSRecord, error) {
	var out struct {
		Data []DNSRecord `json:"data"`
	}
	if err := s.c.Do(ctx, "POST", zonePath(zoneID)+"/records/batch", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return nonNil(out.Data), nil
}

// ExportZone returns the zone as a BIND zone file: GET /dns/zones/{zone_id}/export.
// Scope: read.
func (s *DNSService) ExportZone(ctx context.Context, zoneID string, opts ...RequestOption) ([]byte, error) {
	var out []byte
	if err := s.c.Do(ctx, "GET", zonePath(zoneID)+"/export", nil, nil, &out, append([]RequestOption{WithRequestHeader("Accept", "text/dns")}, opts...)...); err != nil {
		return nil, err
	}
	return out, nil
}

// ImportZone adds or changes the records of a BIND zone file: POST /dns/zones/{zone_id}/import.
// With Replace, records the file does not have are deleted too. It returns
// the records changed. Up to 1000 changes, all or nothing.
// Scope: write, and delete with Replace.
func (s *DNSService) ImportZone(ctx context.Context, zoneID string, req ImportDNSZoneRequest, opts ...RequestOption) ([]DNSRecord, error) {
	var out struct {
		Data []DNSRecord `json:"data"`
	}
	if err := s.c.Do(ctx, "POST", zonePath(zoneID)+"/import", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return nonNil(out.Data), nil
}

// DNSFailureError is returned by WaitForZone and WaitForRecord when the zone
// or record ended as failed. FailureCode is stable.
type DNSFailureError struct {
	// Kind is "zone" or "record".
	Kind        string
	ID          string
	FailureCode string
}

func (e *DNSFailureError) Error() string {
	msg := fmt.Sprintf("pantech: DNS %s %s failed", e.Kind, e.ID)
	if e.FailureCode != "" {
		msg += " (" + e.FailureCode + ")"
	}
	return msg
}

// WaitForZone polls GET /dns/zones/{id} until the zone is active, and returns
// it. A failed zone returns it with a *DNSFailureError. A zone awaiting
// delegation becomes active only once the registrar points at its
// NameServers, which can take hours: bound the wait with opts.Timeout or ctx,
// and use VerifyZone to check at once after changing the registrar. opts may
// be nil.
func (s *DNSService) WaitForZone(ctx context.Context, zoneID string, opts *WaitOptions) (*DNSZone, error) {
	var last *DNSZone
	err := s.c.poll(ctx, opts, func(ctx context.Context) (bool, error) {
		zone, err := s.GetZone(ctx, zoneID)
		if err != nil {
			return false, err
		}
		last = zone
		switch zone.Status {
		case DNSZoneActive:
			return true, nil
		case DNSZoneFailed:
			return true, &DNSFailureError{Kind: "zone", ID: zone.ID, FailureCode: deref(zone.FailureCode)}
		}
		return false, nil
	})
	return last, err
}

// WaitForRecord polls GET /dns/zones/{zone_id}/records/{id} until the record
// is active, that is served by Pantech DNS, and returns it. In a zone still
// awaiting delegation that happens only once the zone is active. A failed record
// returns it with a *DNSFailureError; it keeps its last good values. opts may
// be nil.
func (s *DNSService) WaitForRecord(ctx context.Context, zoneID, recordID string, opts *WaitOptions) (*DNSRecord, error) {
	var last *DNSRecord
	err := s.c.poll(ctx, opts, func(ctx context.Context) (bool, error) {
		record, err := s.GetRecord(ctx, zoneID, recordID)
		if err != nil {
			return false, err
		}
		last = record
		switch record.Status {
		case DNSRecordActive:
			return true, nil
		case DNSRecordFailed:
			return true, &DNSFailureError{Kind: "record", ID: record.ID, FailureCode: deref(record.FailureCode)}
		}
		return false, nil
	})
	return last, err
}

func zonePath(zoneID string) string { return "/dns/zones/" + url.PathEscape(zoneID) }

func recordPath(zoneID, recordID string) string {
	return zonePath(zoneID) + "/records/" + url.PathEscape(recordID)
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// nonNil turns a missing list into an empty one, so callers can range and
// len without a nil check, as listAll does.
func nonNil[T any](items []T) []T {
	if items == nil {
		return []T{}
	}
	return items
}
