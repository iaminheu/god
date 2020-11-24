package handler

import (
	"compress/gzip"
	"git.zc0901.com/go/god/api/httpx"
	"net/http"
	"strings"
)

const gzipEncoding = "gzip"

// API Gzip 中间件
func GzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get(httpx.ContentEncoding), gzipEncoding) {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			r.Body = reader
		}

		next.ServeHTTP(w, r)
	})
}
