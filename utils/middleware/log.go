package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

func LoggerMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			recorder := httptest.NewRecorder()
			next.ServeHTTP(recorder, r)

			for k, v := range recorder.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(recorder.Code)
			recorder.Body.WriteTo(w)

			responseTime := time.Since(start).Seconds()
			logMessage := fmt.Sprintf("%s - [%s] - \"%s %s %s\" %d %s - [%s]\n",
				r.RemoteAddr,
				time.Now().Format(time.RFC1123),
				r.Method,
				r.URL.Path,
				r.Proto,
				recorder.Code,
				r.UserAgent(),
				fmt.Sprintf("%.9fµs", responseTime),
			)
			log.Print(logMessage)
		})
	}
}

func ApplyMiddleware(h http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	handler := http.Handler(h)
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler.ServeHTTP
}
