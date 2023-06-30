package mgzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressResponseWriter struct {
	w  http.ResponseWriter
	gw *gzip.Writer
}

func newCompressResponseWriter(w http.ResponseWriter) *compressResponseWriter {
	return &compressResponseWriter{
		w:  w,
		gw: gzip.NewWriter(w),
	}
}
func (crw *compressResponseWriter) Write(b []byte) (int, error) {

	return crw.gw.Write(b)

}
func (crw *compressResponseWriter) Header() http.Header {
	return crw.w.Header()
}
func (crw *compressResponseWriter) WriteHeader(statusCode int) {
	ct := crw.w.Header().Get("Content-Type")
	if ct == "application/json" || strings.Contains(ct, "text") {
		crw.w.Header().Set("Content-Encoding", _encoding)
	}
	crw.w.WriteHeader(statusCode)
}
func (crw *compressResponseWriter) Close() error {
	return crw.gw.Close()
}

type decompressRequest struct {
	r  io.ReadCloser
	gz *gzip.Reader
}

func newDecompressRequest(r io.ReadCloser) (*decompressRequest, error) {
	f, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &decompressRequest{r: r, gz: f}, nil
}
func (dr *decompressRequest) Close() error {
	err := dr.r.Close()
	if err != nil {
		return err
	}
	return dr.gz.Close()
}
func (dr *decompressRequest) Read(data []byte) (int, error) {
	return dr.gz.Read(data)

}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalWriter := w
		canEncodingGzip := strings.Contains(r.Header.Get("Accept-Encoding"), _encoding)
		if canEncodingGzip {
			enc := newCompressResponseWriter(w)
			originalWriter = enc
			defer enc.Close()
		}

		isEncodedGzip := r.Header.Get("Content-Encoding") == _encoding
		if isEncodedGzip {
			body, err := newDecompressRequest(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = body
			defer body.Close()

		}

		next.ServeHTTP(originalWriter, r)

	})
}
