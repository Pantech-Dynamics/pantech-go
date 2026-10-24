package pantech

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/url"
	"strings"
)

// Stream follows an operation over Server-Sent Events: GET
// /operations/{operation_id}/events. It yields the operation each time it
// changes, and ends after it succeeds or fails. If the stream stops early the
// API sends an error event, yielded as an *Error (OPERATION_STREAM_TIMEOUT,
// OPERATION_STREAM_UNAVAILABLE or UNAUTHENTICATED): reconnect, or use Wait.
// The stream is opened once and never retried. Stop ranging to close it.
// Scope: read.
func (s *OperationsService) Stream(ctx context.Context, operationID string) iter.Seq2[*Operation, error] {
	return func(yield func(*Operation, error) bool) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		path := "/operations/" + url.PathEscape(operationID) + "/events"
		req, err := s.c.newRequest(ctx, http.MethodGet, s.c.baseURL+path, nil, "", http.Header{"Accept": {"text/event-stream"}})
		if err != nil {
			yield(nil, err)
			return
		}
		// The stream may stay open for up to 15 minutes, longer than the
		// client's per-request timeout, so it uses a copy without one.
		hc := *s.c.http
		hc.Timeout = 0
		res, err := hc.Do(req)
		if err != nil {
			yield(nil, fmt.Errorf("pantech: GET %s: %w", path, err))
			return
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode > 299 {
			body, _ := readCapped(res.Body)
			yield(nil, newError(res.StatusCode, res.Header, body))
			return
		}
		for ev, err := range readEvents(res.Body) {
			if err != nil {
				if ctx.Err() == nil {
					yield(nil, fmt.Errorf("pantech: reading the operation stream: %w", err))
				}
				return
			}
			switch ev.name {
			case "operation":
				var op Operation
				if err := json.Unmarshal([]byte(ev.data), &op); err != nil {
					yield(nil, fmt.Errorf("pantech: decoding an operation event: %w", err))
					return
				}
				if !yield(&op, nil) {
					return
				}
			case "error":
				var p struct {
					Code      string `json:"code"`
					Detail    string `json:"detail"`
					RequestID string `json:"request_id"`
				}
				_ = json.Unmarshal([]byte(ev.data), &p)
				yield(nil, &Error{Status: res.StatusCode, Code: p.Code, Detail: p.Detail, RequestID: p.RequestID, Title: "Operation stream ended"})
				return
			}
		}
	}
}

type sseEvent struct {
	name string
	data string
}

// readEvents parses a text/event-stream body into events. Comments (such as
// ": keepalive") are skipped.
func readEvents(r io.Reader) iter.Seq2[sseEvent, error] {
	return func(yield func(sseEvent, error) bool) {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 64<<10), 4<<20)
		var ev sseEvent
		var data []string
		for sc.Scan() {
			line := sc.Text()
			switch {
			case line == "":
				if len(data) > 0 || ev.name != "" {
					ev.data = strings.Join(data, "\n")
					if ev.name == "" {
						ev.name = "message"
					}
					if !yield(ev, nil) {
						return
					}
				}
				ev, data = sseEvent{}, nil
			case strings.HasPrefix(line, ":"):
			default:
				field, value, _ := strings.Cut(line, ":")
				value = strings.TrimPrefix(value, " ")
				switch field {
				case "event":
					ev.name = value
				case "data":
					data = append(data, value)
				}
			}
		}
		if err := sc.Err(); err != nil {
			yield(sseEvent{}, err)
		}
	}
}
