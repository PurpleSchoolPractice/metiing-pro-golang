package server

import (
	"context"

	"errors"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type Server struct {
	log    Logger
	router *chi.Mux
	http   *http.Server
	app    Application
}

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Application interface {
	// TODO
}

func NewServer(logger *logger.Logger, app Application, router *chi.Mux) *Server {
	router.Use(middleware.RequestLogger(logger))

	server := &Server{
		log:    logger,
		router: router,
		app:    app,
	}
	return server
}

func (s *Server) Start(ctx context.Context) error {
	s.http = &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}
	go func() {

		s.log.Info("Starting server on " + s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("Error starting server: %v", err)
		}
	}()
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {

	s.log.Info("Stopping server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		s.log.Error("Error shutting down server: %v", err)
		return err
	}
	s.log.Info("Server stopped")
	return nil
}
