package shim

import (
	"mime"
	"net/http"
	"strings"
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

func setContentTypeIfNotPresent(headers http.Header, body []byte) {
	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", http.DetectContentType(body))
	}
}
