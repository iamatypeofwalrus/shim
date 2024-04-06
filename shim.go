package shim

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Shim provides a thin layer between your traditional http.Handler based application and AWS Lambda + API Gateway.
type Shim struct {
	Handler http.Handler
	Log     Log
}

// New returns an initialized Shim with the provided http.Handler. If no http.Handler is provided New will use http.DefaultServiceMux
func New(h http.Handler, options ...func(*Shim)) *Shim {
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

func SetDebugWithSlog(l slog.Logger) func(*Shim) {
	return func(s *Shim) {
		s.Log = slogAdapter{Logger: l}
	}
}

// Handle converts an APIGatewayProxyRequest converts an APIGatewayProxyRequest into an http.Request and passes it to the given http.Handler
// along with a ResponseWriter. The response from the handler is converted into an APIGatewayProxyResponse.
func (s *Shim) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	s.printf("event request: %+v\n", request)

	httpReq, err := NewHttpRequestFromAPIGatewayProxyRequest(ctx, request)
	if err != nil {
		s.printf("received an error while constructing http request from API Gateway request event\n")
		return events.APIGatewayProxyResponse{}, err
	}
	s.printf("http request: %+v", httpReq)
	rw := NewResponseWriter()

	s.printf("calling ServeHTTP on shim handler\n")
	s.Handler.ServeHTTP(rw, httpReq)
	s.printf("received response: %+v\n", rw)

	resp := NewAPIGatewayProxyResponse(rw)
	s.printf("api gateway proxy response: %+v\n", resp)
	return resp, nil
}

// HandleRestApiRequests converts an APIGatewayProxyRequest into an http.Request and passes it to the http.Handler. Http responses are converted
// into APIGatewayProxyResponse.
func (s *Shim) HandleRestApiRequests(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return s.Handle(ctx, request)
}

// HandleHttpApiRequests converts an APIGatewayV2HTTPRequest into an http.Request and passes it to the http.Handler. Http responses are converted
// into APIGatewayV2HTTPResponse
func (s *Shim) HandleHttpApiRequests(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	s.printf("shim received event request: %+v", request)

	httpReq, err := NewHttpRequestFromAPIGatewayV2HTTPRequest(ctx, request)
	if err != nil {
		s.printf("received error while converting APIGatewayV2HTTPRequest into http request: %v\n", err)
		return events.APIGatewayV2HTTPResponse{}, err
	}

	s.printf("generated http request: %+v\n", httpReq)

	rw := NewResponseWriter()
	s.printf("calling ServeHTTP on shim handler\n")
	s.Handler.ServeHTTP(rw, httpReq)
	s.printf("received response: %+v\n", rw)

	resp := NewApiGatewayV2HttpResponse(rw)
	s.printf("api gateway v2 http response: %+v\n", resp)

	return resp, nil
}

func (s *Shim) printf(format string, v ...interface{}) {
	if s.Log != nil {
		s.Log.Printf(format, v...)
	}
}
