package server

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"patitomedi/video-call/internal/hub"
)

type Server struct {
	hub    *hub.Hub
	mux    *http.ServeMux
}

func New(h *hub.Hub) *Server {
	s := &Server{hub: h, mux: http.NewServeMux()}
	s.mux.HandleFunc("/ws/video", s.handleWS)
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.Handle("/metrics", promhttp.Handler())
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "video-call"})
}
