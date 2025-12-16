package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"shortlink/internal/api/http/handler"
	"shortlink/internal/shortener"
	"time"
)

type Server struct {
	httpServer *http.Server
	service    *shortener.Service
	logger     *log.Logger
}

func NewServer(port string, service *shortener.Service) *Server {
	logger := log.New(os.Stdout, "[HTTP Server] ", log.LstdFlags|log.Lshortfile)
	linkAPIHandler := handler.NewLinkAPI(service, logger)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/links", linkAPIHandler.CreateLink)
	mux.HandleFunc("GET /", linkAPIHandler.RedirectLink)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	})

	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + port,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		service: service,
		logger:  logger,
	}
}

func (s *Server) Start() error {
	s.logger.Printf("Server listening on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()

}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Printf("Server shutting down\n")
	return s.httpServer.Shutdown(ctx)
}
