package shim

import (
	"encoding/base64"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func NewApiGatewayV2HttpResponse(rw *ResponseWriter) events.APIGatewayV2HTTPResponse {
	resp := events.APIGatewayV2HTTPResponse{
		StatusCode: rw.Code,
	}

	headers := rw.Headers
	headers[httpHeaderContentType] = []string{http.DetectContentType(rw.Body.Bytes())}

	resp.MultiValueHeaders = headers

	if shouldConvertToBase64(rw.Headers.Get(httpHeaderContentType)) {
		resp.Body = base64.StdEncoding.EncodeToString(rw.Body.Bytes())
		resp.IsBase64Encoded = true
	} else {
		resp.Body = string(rw.Body.String())
	}

	return resp
}
