package shim

import (
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
