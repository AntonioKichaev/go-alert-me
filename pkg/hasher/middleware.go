package hasher

import (
	"bytes"
	"io"
	"net/http"
)

const _shaHeader = "HashSHA256"

func HasherMiddleware(key string) func(http.Handler) http.Handler {
	h := NewHasher(key)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			originalWriter := w
			sha := r.Header.Get(_shaHeader)

			if len(sha) != 0 {
				buf, _ := io.ReadAll(r.Body)
				sign := h.Sign(buf)
				if sign != sha {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(buf))
			}
			if len(key) != 0 {
				originalWriter = newHashShaResponseWriter(w, h)
			}

			next.ServeHTTP(originalWriter, r)

		})
	}
}

type hashShaResponseWriter struct {
	w    http.ResponseWriter
	hash *Hasher
	sign string
}

func newHashShaResponseWriter(w http.ResponseWriter, h *Hasher) *hashShaResponseWriter {
	return &hashShaResponseWriter{
		w:    w,
		hash: h,
	}
}
func (crw *hashShaResponseWriter) Write(b []byte) (int, error) {
	crw.sign = crw.hash.Sign(b)
	return crw.w.Write(b)

}
func (crw *hashShaResponseWriter) Header() http.Header {
	return crw.w.Header()
}
func (crw *hashShaResponseWriter) WriteHeader(statusCode int) {
	crw.w.Header().Set(_shaHeader, crw.sign)
	crw.w.WriteHeader(statusCode)
}
