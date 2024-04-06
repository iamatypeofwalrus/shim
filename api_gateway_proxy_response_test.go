package shim

import (
	"encoding/base64"
	"net/http"
	"sort"
	"strings"
	"testing"
)

func TestFormatHeadersCanonicalizesKeyNames(t *testing.T) {
	h := make(http.Header)
	encoding := "gzip"
	h["accept-encoding"] = []string{"gzip"}
	f := formatHeaders(h)

	v, ok := f["Accept-Encoding"]
	if !ok {
		t.Error("expected to get key Accept-Encoding to have a value")
	}

	if v != encoding {
		t.Errorf("expected encoding to be %v but was %v", encoding, v)
	}
}

func TestFormatHeadersHandlesMultipleValusForAKey(t *testing.T) {
	h := make(http.Header)
	gzip := "gzip"
	deflate := "deflate"
	encodings := []string{gzip, deflate}
	h["Accept-Encoding"] = encodings
	f := formatHeaders(h)

	v, ok := f["Accept-Encoding"]
	if !ok {
		t.Fatal("expected value for Accept-Encoding to be present")
	}

	encodingVals := strings.Split(v, multipleValueSeperator)
	if len(encodingVals) != len(encodings) {
		t.Fatalf("expected number of encoded values to be %v but was %v\n", len(encodings), len(encodingVals))
	}

	sort.Strings(encodings)
	sort.Strings(encodingVals)

	for i, encoding := range encodings {
		if encodingVals[i] != encoding {
			t.Errorf("expected encoding at %v to be %v but was %v\n", i, encoding, encodingVals[i])
		}
	}
}

func TestNewAPIGatewayProxyResponseConvertsNonTextResponsesToBase64(t *testing.T) {
	input := "hello, world"
	body, err := gzipString(input)
	if err != nil {
		t.Fatalf("unable to gzip string: %v", err)
	}

	rw := NewResponseWriter()
	rw.Write([]byte(body))
	rw.Headers.Set(httpHeaderContentType, "application/octet-stream")
	resp := NewAPIGatewayProxyResponse(rw)

	if !resp.IsBase64Encoded {
		t.Fatal("expected IsBase64Encoded to be true but was false")
	}

	decodedBody, err := base64.StdEncoding.DecodeString(resp.Body)
	if err != nil {
		t.Fatal("expected error from base64 decode to be nil but was", err)
	}

	decodedGunzippedInput, err := gunzipBytes(decodedBody)
	if err != nil {
		t.Fatal("expected error from gunzip to be nil but was", err)
	}

	if input != decodedGunzippedInput {
		t.Errorf("expected decodedBody to be %v but was %v", input, decodedGunzippedInput)
	}
}

func TestNewAPIGatewayProxyResponseDoesNotConvertTextResponsesToBase64(t *testing.T) {
	body := "hello, world"
	rw := NewResponseWriter()
	rw.Write([]byte(body))
	resp := NewAPIGatewayProxyResponse(rw)

	if resp.IsBase64Encoded {
		t.Error("expected response not be base64 encoded")
	}

	if resp.Body != body {
		t.Errorf("expected body to be %v but was %v", body, resp.Body)
	}
}

func TestShouldConvertToBase64(t *testing.T) {
	cases := []struct {
		in  string
		out bool
	}{
		{in: "//", out: true},
		{in: "text/plain", out: false},
		{in: "application/protobuf", out: true},
		{in: "application/json", out: false},
		{in: "application/json; charset=utf-8", out: false},
	}

	for _, c := range cases {
		out := shouldConvertToBase64(c.in)
		if out != c.out {
			t.Errorf("for %v expected %v but was %v", c.in, c.out, out)
		}
	}
}
