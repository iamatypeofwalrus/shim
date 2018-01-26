package shim

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// NewHTTPRequest creates an *http.Request from a context.Context and an events.APIGatewayProxyRequest
func NewHTTPRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*http.Request, error) {
	u := url.URL{}
	u.Path = request.Path

	// Query parameters may or may not present, but if they are pull them out
	// and encode them into the URL
	if len(request.QueryStringParameters) > 0 {
		queryParams := url.Values{}
		for k, v := range request.QueryStringParameters {
			queryParams.Add(k, v)
		}
		u.RawQuery = queryParams.Encode()
	}

	req, err := http.NewRequest(
		request.HTTPMethod,
		u.String(),
		strings.NewReader(request.Body),
	)

	if err != nil {
		return nil, err
	}

	return req.WithContext(ctx), nil
}
