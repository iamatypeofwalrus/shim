package shim

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// NewAPIGatewayProxyResponse converts a shim.ResponseWriter into an events.APIGatewayProxyResponse
func NewAPIGatewayProxyResponse(rw *ResponseWriter) events.APIGatewayProxyResponse {
	headers := formatHeaders(rw.Headers)

	setDefaultContentType(headers, rw.Body.Bytes())

	return events.APIGatewayProxyResponse{
		StatusCode: rw.Code,
		Body:       rw.Body.String(),
		Headers:    headers,
	}
}

func formatHeaders(h http.Header) map[string]string {
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

func setDefaultContentType(lambdaHeaders map[string]string, body []byte) {
	if _, ok := lambdaHeaders[contentType]; !ok {
		lambdaHeaders[contentType] = http.DetectContentType(body)
	}
}
