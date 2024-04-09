package shim

import (
	"bytes"
	"net/http"
)

var headerContentType = "Content-Type"

// NewResponseWriter returns an ResponseWriter with the headers properly initialized
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		Headers: make(http.Header),
	}
}

// ResponseWriter adheres to the http.ResponseWriter interface and makes the HTTP Status Code, Headers, and Body publicly
// accessible
type ResponseWriter struct {
	Code    int
	Headers http.Header
	Body    bytes.Buffer
}

// Header adheres the http.ResponseWriter interface
func (rw *ResponseWriter) Header() http.Header {
	return rw.Headers
}

// Write adheres to the io.Writer interface
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if rw.Code == 0 {
		rw.WriteHeader(http.StatusOK)
	}

	if rw.Header().Get(headerContentType) == "" {
		rw.Header().Set(headerContentType, http.DetectContentType(b))
	}

	return rw.Body.Write(b)
}

// WriteHeader adheres to the http.ResponseWriter interface
func (rw *ResponseWriter) WriteHeader(c int) {
	rw.Code = c
}
