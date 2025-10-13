package httpsh

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"treeTwo/domain"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/redis"
)

type Server struct {
	port     string
	logger   zerolog.Logger
	service  domain.Service
	storage  domain.Storage
	handlers domain.Handlers
	redis    *redis.Client
	srv      *http.Server
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewServer(port string, logger zerolog.Logger, service domain.Service, storage domain.Storage, handlers domain.Handlers, redis *redis.Client) *Server {
	return &Server{
		port:     port,
		logger:   logger,
		service:  service,
		storage:  storage,
		handlers: handlers,
		redis:    redis,
	}
}

func (s *Server) RunServer() error {
	if s.port == "" {
		s.logger.Error().Msg(("Port is not set"))
		os.Exit(1)
	}
	mux := http.NewServeMux()
	RegisterRoutes(mux, s.handlers)
	mux.Handle("/", &handleDef{})
	srv := &http.Server{
		Addr:         "0.0.0.0:" + s.port,
		Handler:      withCORS(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	s.srv = srv

	s.logger.Info().Msg("Starting server on port: " + s.port)
	return s.srv.ListenAndServe()
}

type handleDef struct{}

func (h *handleDef) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	statusCode := http.StatusNotFound
	if r.Method == http.MethodOptions {
		statusCode = http.StatusOK
	}

	w.WriteHeader(statusCode)
	response := ErrorResponse{
		Message: "Endpoint not found or method not allowed",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
