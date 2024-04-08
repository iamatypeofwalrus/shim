package shim

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// NewHttpRequestFromAPIGatewayV2HTTPRequest creates an *http.Request from the context passed from the Lambda library, and the event itself.
func NewHttpRequestFromAPIGatewayV2HTTPRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*http.Request, error) {
	u, err := url.Parse(event.RawPath)
	if err != nil {
		return nil, fmt.Errorf("shim could not parse path from event: %w", err)
	}

	if event.RawQueryString != "" {
		u.RawQuery = event.RawQueryString
	}

	body := event.Body
	if event.IsBase64Encoded {
		d, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, fmt.Errorf("shim encountered an error while base64 decoding request body: %w", err)
		}

		body = string(d)
	}

	req, err := http.NewRequest(
		event.RequestContext.HTTP.Method,
		u.String(),
		strings.NewReader(body),
	)

	if err != nil {
		return nil, fmt.Errorf("shim could not create http request from event: %w", err)
	}

	// Xray tracing is passed in via x-amzn-trace-id header that is on the Lambda Event
	for h, v := range event.Headers {
		req.Header.Set(h, v)
	}

	requestID := event.RequestContext.RequestID
	if requestID != "" {
		req.Header.Set("x-request-id", requestID)
	}

	req.URL.Host = req.Header.Get("Host")
	req.Host = req.Header.Get("Host")

	req.RemoteAddr = event.RequestContext.HTTP.SourceIP

	if req.Header.Get(contentLength) == "" && body != "" {
		req.Header.Set(contentLength, strconv.Itoa(len(body)))
	}

	req = req.WithContext(ctx)

	return req, nil
}
