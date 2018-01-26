package shim

import (
	"context"
	"net/http"

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
	s.printf("event request: %+v\n", request)

	httpReq, err := NewHTTPRequest(ctx, request)
	if err != nil {
		s.println("received an error while constructing http request from API Gateway request event")
		return events.APIGatewayProxyResponse{}, err
	}
	s.printf("http request: %+v", httpReq)
	rw := NewResponseWriter()

	s.println("calling ServeHTTP on shim handler")
	s.Handler.ServeHTTP(rw, httpReq)
	s.printf("received response: %+v\n", rw)

	resp := NewAPIGatewayProxyResponse(rw)
	s.printf("api gateway proxy response: %+v\n", resp)
	return resp, nil
}

func (s *Shim) println(v ...interface{}) {
	if s.Log != nil {
		s.Log.Println(v...)
	}
}

func (s *Shim) printf(format string, v ...interface{}) {
	if s.Log != nil {
		s.Log.Printf(format, v...)
	}
}
