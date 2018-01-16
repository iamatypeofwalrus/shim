package shim

import (
	"fmt"
	"net/http"
	"testing"
)

func TestShim(t *testing.T) {
	msg := "hello, world"
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, msg)
	})

	s := &Shim{
		Handler: mux,
	}

	event
}
