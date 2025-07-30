package middleware

import (
	"log"
	"net/http"
	"time"
)

var middlewares = []func(next http.Handler) http.Handler{
	// Last middleware executes first
	ErrorRecovery,
	Logger,
}

func Chain(next http.HandlerFunc) http.Handler {
	n := http.Handler(next)
	for _, m := range middlewares {
		n = m(n)
	}
	return n
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := statusRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(&rec, r)
		duration := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rec.status, duration)
	})
}

func ErrorRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Caught internal error:", r)
				http.Error(w, "Internal Server Error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}
