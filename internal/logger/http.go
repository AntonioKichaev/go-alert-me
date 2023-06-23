package logger

import "net/http"

type responseData struct {
	Status int
	Size   int
}
type loggingResponseWriter struct {
	http.ResponseWriter
	Rd *responseData
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.Rd.Size += size
	return size, err
}
func (lwr *loggingResponseWriter) WriteHeader(statusCode int) {
	lwr.ResponseWriter.WriteHeader(statusCode)
	lwr.Rd.Status = statusCode
}
