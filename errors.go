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
