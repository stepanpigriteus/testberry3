package httpsh

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"threeFive/domain"
	"threeFive/internal/db"

	"github.com/rs/zerolog"
)

type Server struct {
	port    string
	logger  zerolog.Logger
	service domain.Service
	db      *db.DB

	handlers domain.Handlers
	srv      *http.Server
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewServer(port string, logger zerolog.Logger, service domain.Service, handlers domain.Handlers, db *db.DB) *Server {
	return &Server{
		port:     port,
		logger:   logger,
		service:  service,
		db:       db,
		handlers: handlers,
	}
}

func (s *Server) RunServer(ctx context.Context) error {
	if s.port == "" {
		s.logger.Error().Msg("Port is not set")
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

	serverErr := make(chan error, 1)

	go func() {
		s.logger.Info().Msg("Starting server on port: " + s.port)
		// Отправляем только реальные ошибки, не ErrServerClosed
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		s.logger.Error().Err(err).Msg("Server error")
		return err
	case <-stop:
		s.logger.Info().Msg("Shutdown signal received")
	case <-ctx.Done():
		s.logger.Info().Msg("Context cancelled")
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.logger.Info().Msg("Shutting down server gracefully...")

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		s.logger.Error().Err(err).Msg("Server shutdown failed, forcing close")
		if closeErr := s.srv.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}

	// Закрываем соединение с БД
	if s.db != nil {
		s.logger.Info().Msg("Closing database connection...")
		if err := s.db.Close(); err != nil {
			s.logger.Error().Err(err).Msg("Failed to close database")
			return err
		}
	}

	s.logger.Info().Msg("Server exited correctly")
	return nil
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
