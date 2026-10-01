package main

import (
	"log"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter

	statusCode   int
	bytesWritten int
}

func (r *responseRecorder) WriteHeader(
	statusCode int,
) {
	if r.statusCode != 0 {
		return
	}

	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(
	data []byte,
) (int, error) {
	if r.statusCode == 0 {
		r.WriteHeader(http.StatusOK)
	}

	n, err := r.ResponseWriter.Write(data)
	r.bytesWritten += n

	return n, err
}

func loggingMiddleware(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			startedAt := time.Now()

			recorder := &responseRecorder{
				ResponseWriter: w,
			}

			next.ServeHTTP(
				recorder,
				r,
			)

			statusCode := recorder.statusCode

			if statusCode == 0 {
				statusCode = http.StatusOK
			}

			log.Printf(
				"%s %s -> %d | %d байт | %s",
				r.Method,
				r.URL.Path,
				statusCode,
				recorder.bytesWritten,
				time.Since(startedAt).Round(time.Millisecond),
			)
		},
	)
}
