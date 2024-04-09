package shim

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"testing"
)

func TestNewApiGatewayV2HttpResponse_Base64(t *testing.T) {
	body, err := gzipString("hello, world")
	if err != nil {
		t.Fatalf("unable to gzip string: %v", err)
	}

	rw := NewResponseWriter()
	rw.Write([]byte(body))
	rw.Headers.Set(httpHeaderContentType, "application/octet-stream")
	resp := NewApiGatewayV2HttpResponse(rw)

	if resp.StatusCode != rw.Code {
		t.Errorf("expected status code %d, got %d", rw.Code, resp.StatusCode)
	}

	if resp.Body != base64.StdEncoding.EncodeToString(rw.Body.Bytes()) {
		t.Errorf("expected body to be base64 encoded")
	}

	if !resp.IsBase64Encoded {
		t.Errorf("expected IsBase64Encoded to be true")
	}

	if resp.Headers[httpHeaderContentType] == "" {
		t.Errorf("expected Content-Type header to be set")
	}
}

func TestNewApiGatewayV2HttpResponse_NoBase64(t *testing.T) {
	rw := NewResponseWriter()
	rw.Write([]byte("hello, world"))

	resp := NewApiGatewayV2HttpResponse(rw)

	if resp.StatusCode != rw.Code {
		t.Errorf("expected status code %d, got %d", rw.Code, resp.StatusCode)
	}

	if resp.Body != rw.Body.String() {
		t.Errorf("expected body to be '%s', got '%s'", rw.Body.String(), resp.Body)
	}

	if resp.IsBase64Encoded {
		t.Errorf("expected IsBase64Encoded to be false")
	}

	if resp.Headers[httpHeaderContentType] == "" {
		t.Errorf("expected Content-Type header to be set")
	}
}

func gzipString(input string) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write([]byte(input)); err != nil {
		return nil, fmt.Errorf("unable to write to gzip writer: %w", err)
	}

	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("unable to close gzip writer: %w", err)
	}

	return buf.Bytes(), nil
}

func gunzipBytes(input []byte) (string, error) {
	buf := bytes.NewBuffer(input)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return "", fmt.Errorf("unable to create gzip reader: %w", err)
	}
	defer gz.Close()

	res, err := io.ReadAll(gz)
	if err != nil {
		return "", fmt.Errorf("unable to read from gzip reader: %w", err)
	}

	return string(res), nil
}
