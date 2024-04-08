package shim

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestNewHttpRequestFromAPIGatewayV2HTTPRequest(t *testing.T) {
	ctx := context.Background()

	// Create a sample APIGatewayV2HTTPRequest event
	event := events.APIGatewayV2HTTPRequest{
		Version:               "2.0",
		RouteKey:              "GET /hello",
		RawPath:               "/hello",
		RawQueryString:        "name=John&age=30",
		Headers:               map[string]string{"Content-Type": "application/json"},
		QueryStringParameters: map[string]string{"name": "John", "age": "30"},
	}

	// Call the function under test
	req, err := NewHttpRequestFromAPIGatewayV2HTTPRequest(ctx, event)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify the created request
	if req.Method != http.MethodGet {
		t.Errorf("expected GET method, got %s", req.Method)
	}

	if req.URL.Path != "/hello" {
		t.Errorf("expected path '/hello', got %s", req.URL.Path)
	}

	if req.URL.RawQuery != "name=John&age=30" {
		t.Errorf("expected query string 'name=John&age=30', got %s", req.URL.RawQuery)
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected 'Content-Type' header to be 'application/json', got %s", req.Header.Get("Content-Type"))
	}

	if req.Context() != ctx {
		t.Error("expected context to be the same")
	}
}

func TestNewHttpRequestFromAPIGatewayV2HTTPRequest_Base64(t *testing.T) {
	ctx := context.Background()

	// Create a sample APIGatewayV2HTTPRequest event with base64 encoded body
	body := "Hello, World!"
	encodedBody := base64.StdEncoding.EncodeToString([]byte(body))
	event := events.APIGatewayV2HTTPRequest{
		Version:               "2.0",
		RouteKey:              "GET /hello",
		RawPath:               "/hello",
		RawQueryString:        "name=John&age=30",
		Headers:               map[string]string{"Content-Type": "application/json"},
		QueryStringParameters: map[string]string{"name": "John", "age": "30"},
		Body:                  encodedBody,
		IsBase64Encoded:       true,
	}

	// Call the function under test
	req, err := NewHttpRequestFromAPIGatewayV2HTTPRequest(ctx, event)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify the created request
	decodedBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Errorf("unexpected error when reading request body: %v", err)
	}

	if strings.TrimSpace(string(decodedBody)) != body {
		t.Errorf("expected body '%s', got '%s'", body, string(decodedBody))
	}
}

func TestNewHttpRequestFromAPIGatewayV2HTTPRequest_XRayAndLambdaRequestIDs(t *testing.T) {
	requestID := "request-id"
	traceID := "trace-id"

	// Create a sample APIGatewayV2HTTPRequest event
	event := events.APIGatewayV2HTTPRequest{
		Version:        "2.0",
		RouteKey:       "GET /hello",
		RawPath:        "/hello",
		RawQueryString: "name=John&age=30",
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"x-amzn-trace-id": traceID,
		},
		QueryStringParameters: map[string]string{"name": "John", "age": "30"},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RequestID: requestID,
		},
	}

	// Call the function under test
	req, err := NewHttpRequestFromAPIGatewayV2HTTPRequest(context.Background(), event)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check if the X-Ray ID and the Lambda request ID are added as headers to the request
	if req.Header.Get("x-amzn-trace-id") != traceID {
		t.Errorf("expected X-Ray trace ID '%s', got '%s'", traceID, req.Header.Get("X-Amzn-Trace-Id"))
	}

	if req.Header.Get("x-request-id") != requestID {
		t.Errorf("expected Lambda request ID '%s', got '%s'", requestID, req.Header.Get("Lambda-Runtime-Aws-Request-Id"))
	}
}
