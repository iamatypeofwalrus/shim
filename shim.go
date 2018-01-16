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

// Handler is an interface for accepting and responding to API Gateway integration requests in Lambda
type Handler interface {
	Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

// Shim provides a thin layer between your traditional http.Handler based web app and AWS Lambda + API Gateway.
type Shim struct {
	Handler http.Handler
	Log     Log
}

// New returns an initialized Shim
func New(h http.Handler, options ...func(*Shim)) Handler {
	s := &Shim{
		Handler: h,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

// SetDebugLogger is an option function to set the debug logger on a Shim
func SetDebugLogger(l Log) func(*Shim) {
	return func(s *Shim) {
		s.Log = l
	}
}

// Handle converts an APIGatewayProxyRequest into a http.Request, creates a new ResponseWriter,
// and passes them to its http.Handler
func (s *Shim) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	s.Printf("event request: %+v\n", request)
	// TODO: verify that the path here is the raw path with query and path parameters
	req, err := http.NewRequest(
		request.HTTPMethod,
		request.Path,
		strings.NewReader(request.Body),
	)

	if err != nil {
		s.Println("received an error while constructing http request from API Gateway request event")
		return events.APIGatewayProxyResponse{}, err
	}

	s.Printf("http request: %+v", req)
	req = req.WithContext(ctx)

	rw := &ResponseWriter{
		Headers: make(http.Header),
	}
	s.Println("calling ServeHTTP on shim handler")
	s.Handler.ServeHTTP(rw, req)
	s.Printf("received response: %+v\n", rw)

	s.Println("formatting http.Headers into lambda headers")
	headers := FormatHeaders(rw.Headers)

	s.Println("attempting to set default content type")
	SetDefaultContentType(headers, rw.Body.Bytes())

	resp := events.APIGatewayProxyResponse{
		StatusCode: rw.Code,
		Body:       rw.Body.String(),
		Headers:    headers,
	}
	s.Printf("api gateway proxy response: %+v\n", resp)
	return resp, nil
}

// Println is enabled when a Log is passed to Shim
func (s *Shim) Println(v ...interface{}) {
	if s.Log != nil {
		s.Log.Println(v...)
	}
}

// Printf is enabled when a Log is passed to Shim
func (s *Shim) Printf(format string, v ...interface{}) {
	if s.Log != nil {
		s.Log.Printf(format, v...)
	}
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
