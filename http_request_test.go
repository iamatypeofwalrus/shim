package shim

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestNewHTTPRequestPassesQueryStrings(t *testing.T) {
	key := "query"
	val := "param"
	qp := map[string]string{
		key: val,
	}
	event := events.APIGatewayProxyRequest{
		Path:                  "/api",
		HTTPMethod:            http.MethodGet,
		QueryStringParameters: qp,
	}

	req, err := NewHTTPRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("expected error from New HTTPRequest to be nil")
		t.Logf("err: %s\n", err)
	}

	vals := req.URL.Query()
	if len(vals) < 1 {
		t.Error("expected query params values map to be larger than zero")
	}

	urlVals, ok := vals[key]
	if !ok {
		t.Error("expected key in urls params to be present")
	}

	if len(urlVals) != 1 {
		t.Error("expected number of query param values to be exactly one")
	}

	if urlVals[0] != val {
		t.Errorf("expected val to be %v but was %v", val, urlVals[0])
	}
}

func TestNewHTTPRequestDecodesBase64Bodies(t *testing.T) {
	body := "hello, world"
	encodedBody := base64.StdEncoding.EncodeToString([]byte(body))

	event := events.APIGatewayProxyRequest{
		Path:            "/api",
		HTTPMethod:      http.MethodGet,
		Body:            encodedBody,
		IsBase64Encoded: true,
	}

	req, err := NewHTTPRequest(context.TODO(), event)
	if err != nil {
		t.Fatal("execpted error from NewHTTPRequest to be nil but was ", err)
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal("execpted read error to be nil but was ", err)
	}
	req.Body.Close()

	if string(reqBody) != body {
		t.Error("exepcted request body to equal unenccoded body")
		t.Logf("body: %v\n", body)
		t.Logf("reqBody: %v\n", string(reqBody))
	}
}

func TestNewHTTPRequestPassesHeaders(t *testing.T) {
	key := "hello"
	val := "world"
	event := events.APIGatewayProxyRequest{
		Path:       "/api",
		HTTPMethod: http.MethodPost,
		Headers: map[string]string{
			key: val,
		},
	}

	req, err := NewHTTPRequest(context.TODO(), event)
	if err != nil {
		t.Fatal("exepected error from NewHTTPRequest to be nil but was ", err)
	}

	if headerVal := req.Header.Get(key); headerVal != val {
		t.Errorf("expected headerVal to be %v but was %v", val, headerVal)
	}
}

func TestNewHTTPRequestSetsContentLength(t *testing.T) {
	body := "hello, world"
	event := events.APIGatewayProxyRequest{
		Path:       "/api",
		HTTPMethod: http.MethodGet,
		Body:       body,
	}

	req, err := NewHTTPRequest(context.TODO(), event)
	if err != nil {
		t.Fatal("expected error from NewHTTPRequest to be nil but was ", err)
	}

	if contentLength := req.Header.Get(contentLength); contentLength == "" {
		t.Error("expected content length to be set")
	}
}
