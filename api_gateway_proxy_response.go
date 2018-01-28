package shim

import (
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

const (
	multipleValueSeperator = ","
)

// NewAPIGatewayProxyResponse converts a shim.ResponseWriter into an events.APIGatewayProxyResponse
func NewAPIGatewayProxyResponse(rw *ResponseWriter) events.APIGatewayProxyResponse {
	headers := formatHeaders(rw.Headers)
	setDefaultContentType(headers, rw.Body.Bytes())

	// TODO: if body type is not mime type convert body to base64

	return events.APIGatewayProxyResponse{
		StatusCode: rw.Code,
		Body:       rw.Body.String(),
		Headers:    headers,
	}
}

// formatHeaders converts an http.Headers map into a map[string]string which is expected by Lambda.
func formatHeaders(h http.Header) map[string]string {
	headers := make(map[string]string)

	for k, v := range h {
		// e.g. convert accept-encoding into Accept-Encoding
		canonicalKey := http.CanonicalHeaderKey(k)
		var str string
		if len(v) == 1 {
			str = v[0]
		} else if len(v) > 1 {
			// Per RFC 2616 combine headers with multiple values with ","
			// Source: https://www.w3.org/Protocols/rfc2616/rfc2616-sec4.html#sec4.2
			str = strings.Join(v, multipleValueSeperator)
		}

		headers[canonicalKey] = str
	}

	return headers
}

func setDefaultContentType(lambdaHeaders map[string]string, body []byte) {
	if _, ok := lambdaHeaders[contentType]; !ok {
		lambdaHeaders[contentType] = http.DetectContentType(body)
	}
}
