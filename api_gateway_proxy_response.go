package shim

import (
	"encoding/base64"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// NewAPIGatewayProxyResponse converts a shim.ResponseWriter into an events.APIGatewayProxyResponse
func NewAPIGatewayProxyResponse(rw *ResponseWriter) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		StatusCode: rw.Code,
	}

	httpHeaders := rw.Headers
	setContentTypeIfNotPresent(httpHeaders, rw.Body.Bytes())

	headers := formatHeaders(httpHeaders)
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
