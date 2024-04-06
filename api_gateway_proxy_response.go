package shim

import (
	"encoding/base64"
	"mime"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

const (
	httpHeaderContentType  = "Content-Type"
	multipleValueSeperator = ","
	prefixText             = "text/"
)

var (
	// textFormats contains Content-Types that we know to be text, but do not fall under the text/* category
	textFormats = []string{
		"application/json",
		"application/xml",
		"application/javascript",
	}
)

// NewAPIGatewayProxyResponse converts a shim.ResponseWriter into an events.APIGatewayProxyResponse
func NewAPIGatewayProxyResponse(rw *ResponseWriter) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		StatusCode: rw.Code,
	}

	headers := formatHeaders(rw.Headers)
	headers[httpHeaderContentType] = http.DetectContentType(rw.Body.Bytes())
	resp.Headers = headers

	if shouldConvertToBase64(resp.Headers[httpHeaderContentType]) {
		resp.Body = base64.StdEncoding.EncodeToString(rw.Body.Bytes())
		resp.IsBase64Encoded = true
	} else {
		resp.Body = string(rw.Body.String())
	}

	return resp
}

// formatHeaders converts an http.Headers map into a map[string]string. If there are multiple values for a key
// in the http.Headers they are combined together with "," per RFC 2616
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

func shouldConvertToBase64(ct string) bool {
	mimeType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return true
	}

	// Anything prefixed with text/ should not be converted to base64
	if strings.HasPrefix(mimeType, prefixText) {
		return false
	}

	// Range through special cases held in textFormats
	for _, t := range textFormats {
		if t == mimeType {
			return false
		}
	}

	return true
}
