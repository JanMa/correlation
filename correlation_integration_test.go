package correlation

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/urfave/negroni"
)

func TestIntegration(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "bar")
	})

	correlationMiddleware := New(Options{})

	n := negroni.New()
	n.Use(negroni.HandlerFunc(correlationMiddleware.HandlerFuncWithNext))
	n.UseHandler(mux)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	n.ServeHTTP(res, req)

	expectEq(t, res.Code, http.StatusOK)
	expectEq(t, res.Body.String(), "bar")
	expectNeq(t, res.Header().Get(correlationIDHeader), "")
	expectNeq(t, req.Header.Get(correlationIDHeader), "")
}

func TestIntegrationForRequestOnly(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "bar")
	})

	correlationMiddleware := New(Options{})

	n := negroni.New()
	n.Use(negroni.HandlerFunc(correlationMiddleware.HandlerFuncWithNextForRequestOnly))
	n.UseHandler(mux)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	n.ServeHTTP(res, req)

	expectEq(t, res.Code, http.StatusOK)
	expectEq(t, res.Body.String(), "bar")
	expectEq(t, res.Header().Get(correlationIDHeader), "")
	expectNeq(t, req.Header.Get(correlationIDHeader), "")
}
