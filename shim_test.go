package shim

import (
	"context"
	"fmt"
	"net/http"
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

		s := &Shim{
			Handler: mux,
		}

		ctx := context.Background()
		requestEvent := events.APIGatewayProxyRequest{
			HTTPMethod: test.Method,
			Path:       test.Path,
			Body:       test.ReqBody,
		}

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
