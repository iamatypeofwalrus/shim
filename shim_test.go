package shim

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

const (
	helloWorld = "hello, world"
)

var handleFunc = func(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprint(rw, helloWorld)
}

var shimTests = []struct {
	Path    string
	Method  string
	ReqBody string

	HandlerFunc http.HandlerFunc
	HandlePath  string

	Code      int
	ExpectErr bool
}{
	{
		Path:        "/",
		Method:      http.MethodGet,
		ReqBody:     "",
		HandlerFunc: handleFunc,
		HandlePath:  "/",
		Code:        http.StatusOK,
		ExpectErr:   false,
	},
	{
		Path:        "/notFound",
		Method:      http.MethodGet,
		ReqBody:     "",
		HandlerFunc: handleFunc,
		HandlePath:  "/found",
		Code:        http.StatusNotFound,
		ExpectErr:   false,
	},
}

func TestShim(t *testing.T) {
	for _, test := range shimTests {
		mux := http.NewServeMux()
		mux.HandleFunc(
			test.HandlePath,
			test.HandlerFunc,
		)

		s := New(mux)

		requestEvent := events.APIGatewayProxyRequest{
			HTTPMethod: test.Method,
			Path:       test.Path,
			Body:       test.ReqBody,
		}

		ctx := context.Background()
		resp, err := s.Handle(ctx, requestEvent)

		if resp.StatusCode != test.Code {
			t.Errorf("expected response code to be %v but was %v\n", test.Code, resp.StatusCode)
			t.Logf("resp: %+v\n", resp)
		}

		if test.ExpectErr == true && err == nil {
			t.Errorf("expected err but was nil")
		}

		if test.ExpectErr == false && err != nil {
			t.Error("expected err to be nil but was present")
			t.Logf("err: %+v", err)
		}
	}
}

func TestQueryParams(t *testing.T) {
	// Set up Query Params
	qp := make(map[string]string)
	key := "hello"
	value := "world"
	qp[key] = value

	request := events.APIGatewayProxyRequest{
		HTTPMethod:            http.MethodGet,
		Path:                  "/",
		QueryStringParameters: qp,
	}

	// Construct Mux and Shim
	mux := http.NewServeMux()
	var receivedQueryParams url.Values
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		receivedQueryParams = req.URL.Query()
		fmt.Fprint(w, "yup")
	})
	s := New(mux)

	// The magic!
	resp, err := s.Handle(context.Background(), request)

	// assertions
	if err != nil {
		t.Errorf("expected error to be nil but was %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be 200 but was %v", resp.StatusCode)
	}

	if len(receivedQueryParams) != 1 {
		t.Errorf("expected number of query params to be 1 but was %v\n", len(receivedQueryParams))
	}

	queryVal, ok := receivedQueryParams[key]
	if !ok {
		t.Errorf("expected query param val to be present for key %v\n", key)
	}

	if queryVal[0] != value {
		t.Errorf("expected query string value to be %v but was %v\n", value, queryVal[0])
	}
}

func TestBase64(t *testing.T) {
	respBody := "Goodbye, world"
	gzippedBody, err := gzipString(respBody)
	if err != nil {
		t.Fatalf("unable to gzip string: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(gzippedBody))
	})

	s := New(mux)

	body := "hello, world"
	event := events.APIGatewayProxyRequest{
		Path:            "/",
		HTTPMethod:      http.MethodPost,
		IsBase64Encoded: true,
		Body:            base64.StdEncoding.EncodeToString([]byte(body)),
	}

	resp, err := s.Handle(context.TODO(), event)
	if err != nil {
		t.Fatal("exepcted error from Handle to be nil but was", err)
	}

	if !resp.IsBase64Encoded {
		t.Fatal("expected IsBase64Encoded to be true but was false")
	}

	decodedBody, err := base64.StdEncoding.DecodeString(resp.Body)
	if err != nil {
		t.Fatal("expected error from decoding base64 string to be nil but was", err)
	}

	gunzippedBody, err := gunzipBytes(decodedBody)
	if err != nil {
		t.Fatal("expected error from gunzipping string to be nil but was", err)
	}

	if gunzippedBody != respBody {
		t.Error("expected decodedBody and respBody to be the same")
		t.Logf("respBody: %v", respBody)
		t.Logf("decodedBody: %v", decodedBody)
	}
}

func TestHandleHttpApiRequests(t *testing.T) {
	// Create a sample APIGatewayV2HTTPRequest event with base64 encoded body
	body := "Hello, World!"
	event := events.APIGatewayV2HTTPRequest{
		Version:               "2.0",
		RouteKey:              "GET /hello",
		RawPath:               "/hello",
		RawQueryString:        "name=John&age=30",
		Headers:               map[string]string{"Content-Type": "application/json"},
		QueryStringParameters: map[string]string{"name": "John", "age": "30"},
		Body:                  body,
		IsBase64Encoded:       false,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(body))
	})

	s := New(mux)

	// Call the function under test
	resp, err := s.HandleHttpApiRequests(context.Background(), event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if string(resp.Body) != body {
		t.Errorf("expected body '%s', got '%s'", body, string(resp.Body))
	}
}
