package pantech

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Error is an error answer from the API: an application/problem+json document
// with a stable Code to branch on (never branch on Detail, which is prose).
type Error struct {
	// Status is the HTTP status code.
	Status int
	// Code is the stable problem code, such as RESOURCE_NOT_FOUND or
	// VALIDATION_FAILED. It is empty when the answer was not a problem
	// document (for example a proxy's error page).
	Code string
	// Title is the short, fixed summary of the problem type.
	Title string
	// Detail says what went wrong and what to do, in plain language. It is
	// safe to show to users.
	Detail string
	// RequestID identifies the request; quote it to support.
	RequestID string
	// Type is the problem type URI.
	Type string
	// Errors names each invalid field of a 422 VALIDATION_FAILED.
	Errors []FieldError
	// RetryAfter is the wait the API asked for (on 429 and some 503s), or 0.
	RetryAfter time.Duration
}

// FieldError is one invalid field of a validation failure.
type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error reads like "Instance not found (RESOURCE_NOT_FOUND, status 404,
// request req_…)", followed by any field errors.
func (e *Error) Error() string {
	msg := e.Detail
	if msg == "" || (len(e.Errors) > 0 && msg == fieldErrorsDetail(e.Errors)) {
		// A 422's detail repeats its field errors, listed below: show
		// the title instead.
		msg = e.Title
	}
	if msg == "" {
		msg = http.StatusText(e.Status)
	}
	var b strings.Builder
	b.WriteString("pantech: ")
	b.WriteString(msg)
	parts := []string{}
	if e.Code != "" {
		parts = append(parts, e.Code)
	}
	parts = append(parts, fmt.Sprintf("status %d", e.Status))
	if e.RequestID != "" {
		parts = append(parts, "request "+e.RequestID)
	}
	b.WriteString(" (" + strings.Join(parts, ", ") + ")")
	for _, fe := range e.Errors {
		text := fe.Message
		if text == "" {
			text = fe.Code
		}
		if fe.Field != "" {
			fmt.Fprintf(&b, "; %s: %s", fe.Field, text)
		} else {
			fmt.Fprintf(&b, "; %s", text)
		}
	}
	return b.String()
}

// HasField reports whether a validation failure names field with code.
func (e *Error) HasField(field, code string) bool {
	for _, fe := range e.Errors {
		if fe.Field == field && fe.Code == code {
			return true
		}
	}
	return false
}

// IsCode reports whether err is, or wraps, an *Error with this problem code,
// for example IsCode(err, "RESOURCE_NOT_FOUND").
func IsCode(err error, code string) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == code
}

// AsError returns the *Error in err's chain, if there is one.
func AsError(err error) (*Error, bool) {
	var e *Error
	ok := errors.As(err, &e)
	return e, ok
}

// fieldErrorsDetail is the detail the API builds from field errors alone:
// "field: message; field2: message2."
func fieldErrorsDetail(fields []FieldError) string {
	parts := make([]string, 0, len(fields))
	for _, fe := range fields {
		m := strings.TrimSuffix(strings.TrimSpace(fe.Message), ".")
		if fe.Field != "" {
			m = fe.Field + ": " + m
		}
		parts = append(parts, m)
	}
	return strings.Join(parts, "; ") + "."
}

// newError builds an *Error from a non-2xx answer. A body that is not a
// problem document still gives a usable error from the status; an HTML page
// or long text from a proxy is never copied into it.
func newError(status int, header http.Header, body []byte) *Error {
	e := &Error{Status: status}
	var doc struct {
		Type      string       `json:"type"`
		Title     string       `json:"title"`
		Code      string       `json:"code"`
		Detail    string       `json:"detail"`
		RequestID string       `json:"request_id"`
		Errors    []FieldError `json:"errors"`
	}
	if err := json.Unmarshal(body, &doc); err == nil {
		e.Type, e.Title, e.Code, e.Detail, e.RequestID, e.Errors = doc.Type, doc.Title, doc.Code, doc.Detail, doc.RequestID, doc.Errors
	} else if text := strings.TrimSpace(string(body)); text != "" && len(text) <= 300 && !strings.HasPrefix(text, "<") {
		e.Detail = text
	}
	if e.Title == "" {
		e.Title = http.StatusText(status)
	}
	if e.Detail == "" && e.Code == "" && status >= 500 {
		e.Detail = fmt.Sprintf("The API is temporarily unavailable (%d %s). Try again in a moment.", status, http.StatusText(status))
	}
	if header != nil {
		if e.RequestID == "" {
			e.RequestID = header.Get("X-Request-Id")
		}
		e.RetryAfter = parseRetryAfter(header.Get("Retry-After"), time.Now())
	}
	return e
}

// Problem codes of Pantech DNS and domains, in addition to the platform's
// (VALIDATION_FAILED, FORBIDDEN, NOT_FOUND, IDEMPOTENCY_CONFLICT, ...).
// Branch on them with IsCode.
const (
	// DNS is not enabled on this deployment (503).
	CodeDNSUnavailable = "DNS_UNAVAILABLE"
	// Domains are not sold on this deployment (503).
	CodeDomainsUnavailable = "DOMAINS_UNAVAILABLE"

	// 404s.
	CodeDNSZoneNotFound     = "DNS_ZONE_NOT_FOUND"
	CodeDNSRecordNotFound   = "DNS_RECORD_NOT_FOUND"
	CodeDomainNotFound      = "DOMAIN_NOT_FOUND"
	CodeDomainOrderNotFound = "DOMAIN_ORDER_NOT_FOUND"

	// The organization has 25 zones (403).
	CodeDNSZoneLimitReached = "DNS_ZONE_LIMIT_REACHED"
	// 10 zones were created in the last 24 hours (403).
	CodeDNSZoneCreateLimitReached = "DNS_ZONE_CREATE_LIMIT_REACHED"
	// The zone has 10,000 record sets (403).
	CodeDNSRecordLimitReached = "DNS_RECORD_LIMIT_REACHED"
	// The organization already has this zone (409).
	CodeDNSZoneNameTaken = "DNS_ZONE_NAME_TAKEN"
	// Another organization serves this name (409).
	CodeDNSZoneActiveElsewhere = "DNS_ZONE_ACTIVE_ELSEWHERE"
	// A record set with this name and type exists: change it instead (409).
	CodeDNSRecordExists = "DNS_RECORD_EXISTS"
	// A CNAME must stand alone at its name (409).
	CodeDNSCNAMEConflict = "DNS_CNAME_CONFLICT"
	// The apex SOA and NS cannot be changed (409).
	CodeDNSRecordReadOnly = "DNS_RECORD_READ_ONLY"
	// An operator suspended the zone (409).
	CodeDNSZoneSuspended = "DNS_ZONE_SUSPENDED"
	// The zone is being deleted or failed (409).
	CodeDNSZoneNotWritable = "DNS_ZONE_NOT_WRITABLE"
	// A Pantech domain delegates to this zone: point it elsewhere first (409).
	CodeDNSZoneInUseByDomain = "DNS_ZONE_IN_USE_BY_DOMAIN"
	// If-Match or a batch change's version is not the current version (412).
	CodeVersionConflict = "VERSION_CONFLICT"
	// More than 1000 changes in one batch or import (422).
	CodeTooManyChanges = "TOO_MANY_CHANGES"
	// The zone file does not parse (422).
	CodeInvalidZoneFile = "INVALID_ZONE_FILE"

	// Availability cannot be checked now: try again (503).
	CodeDomainLookupUnavailable = "DOMAIN_LOOKUP_UNAVAILABLE"
	// Pantech does not sell this TLD (422).
	CodeTLDNotSold = "TLD_NOT_SOLD"
	// The name is taken, premium or reserved (409).
	CodeDomainNotAvailable = "DOMAIN_NOT_AVAILABLE"
	// The registrar refuses a transfer now (409).
	CodeDomainNotTransferable = "DOMAIN_NOT_TRANSFERABLE"
	// The domain is already on Pantech (409).
	CodeDomainAlreadyHeld = "DOMAIN_ALREADY_HELD"
	// Another order for this domain is running (409).
	CodeDomainOrderInFlight = "DOMAIN_ORDER_IN_FLIGHT"
	// The order is not waiting for a payment, or its window has ended (409).
	CodeDomainOrderNotAwaitingPayment = "DOMAIN_ORDER_NOT_AWAITING_PAYMENT"
	// 20 domain orders in the last 24 hours (403).
	CodeDomainOrderLimitReached = "DOMAIN_ORDER_LIMIT_REACHED"
	// The domain is not active (409).
	CodeDomainNotConfigurable = "DOMAIN_NOT_CONFIGURABLE"
	// Turn the transfer lock off first (409).
	CodeDomainLocked = "DOMAIN_LOCKED"
	// The zone's name is not the domain's (422).
	CodeDomainZoneMismatch = "DOMAIN_ZONE_MISMATCH"
)

// Problem codes of the container registry. Branch on them with IsCode.
const (
	// The registry is not enabled on this deployment (503).
	CodeRegistryUnavailable = "REGISTRY_UNAVAILABLE"
	// The organization has not claimed a namespace yet: SetNamespace first (404).
	CodeRegistryNamespaceNotSet = "REGISTRY_NAMESPACE_NOT_SET"
	// Another organization has this namespace (409).
	CodeRegistryNamespaceTaken = "REGISTRY_NAMESPACE_TAKEN"
	// The namespace cannot change once a repository exists (409).
	CodeRegistryNamespaceLocked = "REGISTRY_NAMESPACE_LOCKED"
	// The namespace name is not valid (422, field name).
	CodeInvalidRegistryNamespace = "INVALID_REGISTRY_NAMESPACE"
	// The repository name is not valid (422, field name).
	CodeInvalidRepositoryName = "INVALID_REPOSITORY_NAME"

	// 404s.
	CodeRepositoryNotFound    = "REPOSITORY_NOT_FOUND"
	CodeTagNotFound           = "TAG_NOT_FOUND"
	CodeRegistryTokenNotFound = "REGISTRY_TOKEN_NOT_FOUND"

	// A repository with this name exists (409).
	CodeRepositoryExists = "REPOSITORY_EXISTS"
	// The organization has 100 repositories (403).
	CodeRepositoryLimitReached = "REPOSITORY_LIMIT_REACHED"
	// The organization uses its 50 GiB; pushes are refused (403).
	CodeRegistryStorageLimitReached = "REGISTRY_STORAGE_LIMIT_REACHED"
	// The organization has 50 registry tokens (403).
	CodeRegistryTokenLimitReached = "REGISTRY_TOKEN_LIMIT_REACHED"
	// The repository has immutable tags (409).
	CodeTagImmutable = "TAG_IMMUTABLE"
	// An operator suspended the repository (409).
	CodeRepositorySuspended = "REPOSITORY_SUSPENDED"
)

// Problem codes of teams and access policies. Branch on them with IsCode.
const (
	// A policy or the member baseline refuses the action; Detail names the action (403).
	CodeAccessDenied = "ACCESS_DENIED"

	// 404s.
	CodeIAMTeamNotFound       = "IAM_TEAM_NOT_FOUND"
	CodeIAMPolicyNotFound     = "IAM_POLICY_NOT_FOUND"
	CodeIAMAttachmentNotFound = "IAM_ATTACHMENT_NOT_FOUND"

	// 409s: the name is taken in the organization.
	CodeIAMTeamNameTaken   = "IAM_TEAM_NAME_TAKEN"
	CodeIAMPolicyNameTaken = "IAM_POLICY_NAME_TAKEN"

	// 403s: 50 teams, 100 policies, 10 attachments per member or team.
	CodeIAMTeamLimitReached       = "IAM_TEAM_LIMIT_REACHED"
	CodeIAMPolicyLimitReached     = "IAM_POLICY_LIMIT_REACHED"
	CodeIAMAttachmentLimitReached = "IAM_ATTACHMENT_LIMIT_REACHED"

	// A Pantech managed policy cannot be changed or deleted (409).
	CodeIAMManagedPolicy = "IAM_MANAGED_POLICY"
	// Only owners act on admins, or on teams containing admins (403).
	CodeIAMAdminTarget = "IAM_ADMIN_TARGET"
	// A statement is not valid; Errors names the field, such as
	// statements[0].actions[1] (422).
	CodeInvalidPolicy = "INVALID_POLICY"
)

// Problem codes of the certificate manager. Branch on them with IsCode.
const (
	// Certificates are not enabled on this deployment (503).
	CodeCertificatesUnavailable = "CERTIFICATES_UNAVAILABLE"
	CodeCertificateNotFound     = "CERTIFICATE_NOT_FOUND"
	// 100 certificates (403).
	CodeCertificateLimitReached = "CERTIFICATE_LIMIT_REACHED"
	// 20 requests and imports in 24 hours (403).
	CodeCertificateRequestLimitReached = "CERTIFICATE_REQUEST_LIMIT_REACHED"
	// The same names were requested 5 times in 7 days; RetryAfter says when (429).
	CodeCertificateRateLimited = "CERTIFICATE_RATE_LIMITED"
	// A Pantech service uses it: detach it first (409).
	CodeCertificateInUse = "CERTIFICATE_IN_USE"
	// Created without Exportable (409).
	CodeCertificateNotExportable = "CERTIFICATE_NOT_EXPORTABLE"
	// There is no issued certificate to export or revoke yet (409).
	CodeCertificateNotIssued = "CERTIFICATE_NOT_ISSUED"
	// Imported, already revoked, or a revocation is pending (409).
	CodeCertificateNotRevocable = "CERTIFICATE_NOT_REVOCABLE"
	// Reimport of a certificate Pantech issues (409).
	CodeCertificateNotImported = "CERTIFICATE_NOT_IMPORTED"
	// Verify of a certificate not waiting for validation (409).
	CodeCertificateNotPending = "CERTIFICATE_NOT_PENDING"
	// A reimport names other names than the certificate (422).
	CodeCertificateNamesMismatch = "CERTIFICATE_NAMES_MISMATCH"
	// Changed at the same time by another request (409).
	CodeConflict = "CONFLICT"
)

// Field codes of VALIDATION_FAILED for HTTPS load balancers and .ng domains.
// Check one with HasField, for example
// e.HasField("protocol", pantech.CodeLoadBalancerTLSNotSupported).
const (
	// `protocol`: `https` is offered on `appliance` load balancers only.
	CodeLoadBalancerTLSNotSupported = "LOAD_BALANCER_TLS_NOT_SUPPORTED"
	// `protocol`: the protocol must be `tcp` or `https`.
	CodeLoadBalancerProtocolNotSupported = "LOAD_BALANCER_PROTOCOL_NOT_SUPPORTED"
	// `certificate_id`: an `https` listener needs an issued certificate.
	CodeCertificateRequired = "CERTIFICATE_REQUIRED"
	// `certificate_id`: the certificate is not issued or renewing, or is being revoked.
	CodeCertificateNotUsable = "CERTIFICATE_NOT_USABLE"
	// `certificate_id`: certificate_id, tls_policy and redirect_http apply to `https` only.
	CodeLoadBalancerTLSFieldsNotAllowed = "LOAD_BALANCER_TLS_FIELDS_NOT_ALLOWED"
	// `tls_policy`: use `modern` or `tls13`.
	CodeInvalidTLSPolicy = "INVALID_TLS_POLICY"
	// `privacy`: the TLD's registry always hides contacts (.ng), so privacy cannot be off.
	CodePrivacyAlwaysOn = "PRIVACY_ALWAYS_ON"
)
