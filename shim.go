package shim

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Handler is an interface for accepting and responding to API Gateway integration requests
type Handler interface {
	Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

// Shim provides a thin layer between your traditional http.Handler based application and AWS Lambda + API Gateway.
type Shim struct {
	Handler http.Handler
	Log     Log
}

// New returns an initialized Shim with the provided http.Handler. If no http.Handler is provided New will use http.DefaultServiceMux
func New(h http.Handler, options ...func(*Shim)) Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	s := &Shim{
		Handler: h,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

// SetDebugLogger is an option function to set the debug logger on a Shim. The debug logger gives insight into the event received from
// APIGateway and how shim transforms the request and response.
func SetDebugLogger(l Log) func(*Shim) {
	return func(s *Shim) {
		s.Log = l
	}
}

// Handle converts an APIGatewayProxyRequest converts an APIGatewayProxyRequest into an http.Request and passes it to the given http.Handler
// along with a ResponseWriter. The response from the handler is converted into an APIGatewayProxyResponse.
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
