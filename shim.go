package shim

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

const (
	contentType = "Content-Type"
)

// New returns an initialized Shim
func New(h http.Handler) *Shim {
	return &Shim{
		Handler: h,
	}
}

// Shim provides
type Shim struct {
	Handler http.Handler
}

// Handle converts an APIGatewayProxyRequest into a http.Request, creates a new ResponseWriter,
// and passes them to its http.Handler
func (s *Shim) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO: verify that the path here is the raw path with query and path parameters
	req, err := http.NewRequest(
		request.HTTPMethod,
		request.Path,
		strings.NewReader(request.Body),
	)

	req = req.WithContext(ctx)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	rw := &ResponseWriter{}
	s.Handler.ServeHTTP(rw, req)

	headers := FormatHeaders(rw.Headers)
	SetDefaultContentType(headers, rw.Body.Bytes())

	return events.APIGatewayProxyResponse{
		StatusCode: rw.Code,
		Body:       rw.Body.String(),
		Headers:    headers,
	}, nil
}

// SetDefaultContentType attempts to detect and set the Content-Type header for a given response body
func SetDefaultContentType(lambdaHeaders map[string]string, body []byte) {
	if _, ok := lambdaHeaders[contentType]; !ok {
		lambdaHeaders[contentType] = http.DetectContentType(body)
	}
}

// FormatHeaders is a convinience function that converts an http.Header into a map[string]string
func FormatHeaders(h http.Header) map[string]string {
	headers := make(map[string]string)

	for k, v := range h {
		// No great options here. Rather than cat-ing the string array together
		// and on average getting it right, I'd rather just take the first item
		// in the array.
		var str string
		if len(v) > 0 {
			str = v[0]
		}

		headers[k] = str
	}

	return headers
}
