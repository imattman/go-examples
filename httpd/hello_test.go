package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloPath(t *testing.T) {
	paths := []struct {
		uri  string
		body string
	}{
		{"/hello/", "Hello, World!"},
		{"/hello/a", "Hello, a!"},
		{"/hello/Charlie", "Hello, Charlie!"},
		{"/hello/foo/bar/baz", "Hello, foo/bar/baz!"},
	}

	// convert the func to HTTP handler and invoke!
	handler := http.HandlerFunc(hello)

	for i, p := range paths {
		// Create a request to pass to our handler.
		// Use nil to represent zero query parameters
		req, err := http.NewRequest("GET", "http://example.com"+p.uri, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to capture handler behavior
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Errorf("%d: Unexpected response code for URI %q: got %v, want %v",
				i,
				p.uri,
				status,
				http.StatusOK)
		}

		t.Logf("body: %q", w.Body.String())
		if body := w.Body.String(); body != p.body {
			t.Errorf("%d: Unexpected body for URI %q: got %q, want %q",
				i,
				p.uri,
				body,
				p.body)
		}
	}
}
