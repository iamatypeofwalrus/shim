package shim

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

const (
	contentLength = "Content-Length"
)

var (
	errDecodingBody              = errors.New("encountered an error while base64 decoding request body")
	errCouldNotParsePath         = errors.New("could not parse path from event")
	errCouldNotCreateHTTPRequest = errors.New("encountered error while create http request")
)

// NewHTTPRequest creates an *http.Request from a context.Context and an events.APIGatewayProxyRequest
func NewHTTPRequest(ctx context.Context, event events.APIGatewayProxyRequest) (*http.Request, error) {
	u, err := url.Parse(event.Path)
	if err != nil {
		return nil, errCouldNotParsePath
	}

	// Query parameters may or may not present, but if they are pull them out
	// and encode them into the URL
	if len(event.QueryStringParameters) > 0 {
		queryParams := url.Values{}
		for k, v := range event.QueryStringParameters {
			queryParams.Add(k, v)
		}
		u.RawQuery = queryParams.Encode()
	}

	// Handle base64 encoding
	body := event.Body
	if event.IsBase64Encoded {
		d, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, errDecodingBody
		}

		body = string(d)
	}

	req, err := http.NewRequest(
		event.HTTPMethod,
		u.String(),
		strings.NewReader(body),
	)

	if err != nil {
		return nil, errCouldNotCreateHTTPRequest
	}

	// Pass along headers
	for h, v := range event.Headers {
		req.Header.Set(h, v)
	}

	// Set Host
	req.URL.Host = req.Header.Get("Host")
	req.Host = req.Header.Get("Host")

	// Pass along remote IP
	req.RemoteAddr = event.RequestContext.Identity.SourceIP

	// Ensure Content-Length is set correctly
	if req.Header.Get(contentLength) == "" && body != "" {
		req.Header.Set(contentLength, strconv.Itoa(len(body)))
	}

	// Pass along context to http.Handler
	req = req.WithContext(ctx)

	return req, nil
}
