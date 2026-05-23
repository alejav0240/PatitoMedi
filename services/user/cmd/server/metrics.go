package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type metrics struct {
	requests      atomic.Uint64
	errors        atomic.Uint64
	loginsOK      atomic.Uint64
	loginsFailed  atomic.Uint64
	registrations atomic.Uint64
}

func (a *app) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "user_service_requests_total %d\n", a.metrics.requests.Load())
	fmt.Fprintf(w, "user_service_errors_total %d\n", a.metrics.errors.Load())
	fmt.Fprintf(w, "user_service_logins_success_total %d\n", a.metrics.loginsOK.Load())
	fmt.Fprintf(w, "user_service_logins_failed_total %d\n", a.metrics.loginsFailed.Load())
	fmt.Fprintf(w, "user_service_registrations_total %d\n", a.metrics.registrations.Load())
}

func (a *app) instrument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.metrics.requests.Add(1)
		rw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		if rw.status >= 400 {
			a.metrics.errors.Add(1)
		}
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
