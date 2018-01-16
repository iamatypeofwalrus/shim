package shim

import (
	"fmt"
	"net/http"
	"testing"
)

func TestResponseWriterSetCodeAndWrite(t *testing.T) {
	rw := NewResponseWriter()
	rw.WriteHeader(http.StatusAccepted)
	fmt.Fprint(rw, "hello, world")
	if rw.Code != http.StatusAccepted {
		t.Errorf("expected status code to be %v but was %v", http.StatusAccepted, rw.Code)
	}
}

func TestResponseWriterFirstWrite(t *testing.T) {
	rw := NewResponseWriter()
	fmt.Fprint(rw, "hello, world")
	if rw.Code != http.StatusOK {
		t.Error("expected first write to response writer to set status code")
	}
}
