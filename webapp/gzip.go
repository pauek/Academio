package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type gzipHandler struct {
	handler http.Handler
	noExpire bool
}

func gzipped(w http.ResponseWriter, r *http.Request, noExpire bool, fn http.HandlerFunc) {
	if noExpire {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 60 * 60 * 24 * 365))
	}
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		fn(w, r)
		return
	}
	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close()
	fn(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
}

func (H *gzipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gzipped(w, r, H.noExpire, func(w http.ResponseWriter, r *http.Request) {
		H.handler.ServeHTTP(w, r)
	})
}

func Gzipped(handler http.Handler) http.Handler {
	return &gzipHandler{handler: handler}
}

func GzippedNoExpire(handler http.Handler) http.Handler {
	return &gzipHandler{handler: handler, noExpire: true}
}

func GzippedFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gzipped(w, r, false, func (w http.ResponseWriter, r *http.Request) { 
			fn(w, r) 
		})
	}
}