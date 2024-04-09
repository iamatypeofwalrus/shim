package shim

import (
	"encoding/base64"
	"strings"
	"unicode/utf8"

	"github.com/aws/aws-lambda-go/events"
)

func NewApiGatewayV2HttpResponse(rw *ResponseWriter) events.APIGatewayV2HTTPResponse {

	var output string
	isBase64Encoded := false

	bytes := rw.Body.Bytes()

	if utf8.Valid(bytes) {
		output = string(bytes)
	} else {
		output = base64.StdEncoding.EncodeToString(bytes)
		isBase64Encoded = true
	}

	headers := make(map[string]string)
	cookies := make([]string, 0)

	for key, values := range rw.Headers {
		if key == "Set-Cookie" {
			cookies = append(cookies, values...)
		} else {
			headers[key] = strings.Join(values, ",")
		}
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode:      rw.Code,
		Headers:         headers,
		Cookies:         cookies,
		IsBase64Encoded: isBase64Encoded,
		Body:            output,
	}
}
