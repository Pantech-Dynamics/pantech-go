package pantech

import (
	"context"
	"fmt"
	"iter"
	"net/url"
)

// Certificate statuses. Renewing keeps the current certificate working.
const (
	CertificatePendingValidation = "pending_validation"
	CertificateIssuing           = "issuing"
	CertificateIssued            = "issued"
	CertificateRenewing          = "renewing"
	CertificateExpired           = "expired"
	CertificateRevoked           = "revoked"
	CertificateFailed            = "failed"
)

// Certificate kinds.
const (
	CertificatePublic   = "public"
	CertificateImported = "imported"
)

// Certificate validation methods.
const (
	ValidationPantechDNS = "pantech_dns"
	ValidationCNAME      = "cname"
)

// Key types a certificate can be requested with.
const (
	KeyTypeECDSAP256 = "ecdsa-p256"
	KeyTypeRSA2048   = "rsa-2048"
)

// CertificatesService manages TLS certificates: issued free by a publicly
// trusted CA and renewed automatically, or imported. A request is validated
// through Pantech DNS when the name is in one of your active zones, else
// through a CNAME you add once at your DNS host (Certificate.Validations).
type CertificatesService struct{ c *Client }

// List lists certificates, newest first: GET /certificates.
// Scope: read.
func (s *CertificatesService) List(ctx context.Context, params *ListCertificatesParams, opts ...RequestOption) (*Page[Certificate], error) {
	var out Page[Certificate]
	if err := s.c.Do(ctx, "GET", "/certificates", params.values(), nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// All iterates over every item of List, fetching pages as needed.
// It stops at the first error, which it yields.
func (s *CertificatesService) All(ctx context.Context, params *ListCertificatesParams, opts ...RequestOption) iter.Seq2[Certificate, error] {
	var p ListCertificatesParams
	if params != nil {
		p = *params
	}
	return paginate(p.Cursor, func(cursor string) (*Page[Certificate], error) {
		p.Cursor = cursor
		return s.List(ctx, &p, opts...)
	})
}

// Request asks for a certificate for up to 100 names: POST /certificates. It
// is returned pending_validation; add the CNAMEs its Validations list for
// names not in a Pantech DNS zone, then WaitForCertificate.
// Scope: write.
func (s *CertificatesService) Request(ctx context.Context, req RequestCertificateRequest, opts ...RequestOption) (*Certificate, error) {
	var out Certificate
	if err := s.c.Do(ctx, "POST", "/certificates", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Import imports a certificate bought elsewhere, with its private key:
// POST /certificates/import. Imported certificates are never renewed;
// replace one in place with Reimport.
// Scope: write.
func (s *CertificatesService) Import(ctx context.Context, req ImportCertificateRequest, opts ...RequestOption) (*Certificate, error) {
	var out Certificate
	if err := s.c.Do(ctx, "POST", "/certificates/import", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a certificate: GET /certificates/{certificate_id}.
// Scope: read.
func (s *CertificatesService) Get(ctx context.Context, certificateID string, opts ...RequestOption) (*Certificate, error) {
	var out Certificate
	if err := s.c.Do(ctx, "GET", certificatePath(certificateID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update changes a certificate's description: PATCH /certificates/{certificate_id}.
// Scope: write.
func (s *CertificatesService) Update(ctx context.Context, certificateID string, req UpdateCertificateRequest, opts ...RequestOption) (*Certificate, error) {
	var out Certificate
	if err := s.c.Do(ctx, "PATCH", certificatePath(certificateID), nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Verify checks the validation CNAMEs now instead of at the next 5-minute
// check: POST /certificates/{certificate_id}/verify.
// Scope: write.
func (s *CertificatesService) Verify(ctx context.Context, certificateID string, opts ...RequestOption) (*Certificate, error) {
	var out Certificate
	if err := s.c.Do(ctx, "POST", certificatePath(certificateID)+"/verify", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Reimport replaces an imported certificate in place, for the same names:
// POST /certificates/{certificate_id}/reimport.
// Scope: write.
func (s *CertificatesService) Reimport(ctx context.Context, certificateID string, req ReimportCertificateRequest, opts ...RequestOption) (*Certificate, error) {
	var out Certificate
	if err := s.c.Do(ctx, "POST", certificatePath(certificateID)+"/reimport", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Export returns the certificate, its chain and its private key:
// POST /certificates/{certificate_id}/export. Only certificates created
// with Exportable can be exported. Every export is recorded and every owner
// and admin is emailed.
// Scope: delete.
func (s *CertificatesService) Export(ctx context.Context, certificateID string, opts ...RequestOption) (*ExportedCertificate, error) {
	var out ExportedCertificate
	if err := s.c.Do(ctx, "POST", certificatePath(certificateID)+"/export", nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Revoke asks the CA to revoke the certificate: POST /certificates/{certificate_id}/revoke.
// It cannot be undone. RevocationPending stays true until the CA confirms.
// Scope: delete.
func (s *CertificatesService) Revoke(ctx context.Context, certificateID string, req RevokeCertificateRequest, opts ...RequestOption) (*Certificate, error) {
	var out Certificate
	if err := s.c.Do(ctx, "POST", certificatePath(certificateID)+"/revoke", nil, &req, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a certificate and its key: DELETE /certificates/{certificate_id}.
// A certificate in use answers CERTIFICATE_IN_USE.
// Scope: delete.
func (s *CertificatesService) Delete(ctx context.Context, certificateID string, opts ...RequestOption) error {
	return s.c.Do(ctx, "DELETE", certificatePath(certificateID), nil, nil, nil, opts...)
}

// CertificateFailureError is returned by WaitForCertificate when the
// certificate failed, expired or was revoked. FailureCode is stable when set.
type CertificateFailureError struct {
	ID          string
	Status      string
	FailureCode string
}

func (e *CertificateFailureError) Error() string {
	msg := fmt.Sprintf("pantech: certificate %s is %s", e.ID, e.Status)
	if e.FailureCode != "" {
		msg += " (" + e.FailureCode + ")"
	}
	return msg
}

// WaitForCertificate polls GET /certificates/{id} until the certificate is
// issued, and returns it. Failed, expired and revoked return it with a
// *CertificateFailureError. A request waiting for CNAMEs can take up to 72
// hours, so bound the wait with opts.Timeout or ctx; the returned
// certificate's Validations show what is still pending. opts may be nil.
func (s *CertificatesService) WaitForCertificate(ctx context.Context, certificateID string, opts *WaitOptions) (*Certificate, error) {
	var last *Certificate
	err := s.c.poll(ctx, opts, func(ctx context.Context) (bool, error) {
		cert, err := s.Get(ctx, certificateID)
		if err != nil {
			return false, err
		}
		last = cert
		switch cert.Status {
		case CertificateIssued:
			return true, nil
		case CertificateFailed, CertificateExpired, CertificateRevoked:
			return true, &CertificateFailureError{ID: cert.ID, Status: cert.Status, FailureCode: deref(cert.FailureCode)}
		}
		return false, nil
	})
	return last, err
}

func certificatePath(certificateID string) string {
	return "/certificates/" + url.PathEscape(certificateID)
}
