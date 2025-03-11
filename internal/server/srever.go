package server

import (
	"context"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/go-chi/chi/v5"
	"log"
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
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Application interface {
	// TODO
}

func NewServer(logger *logger.Logger, app Application) *Server {
	router := chi.NewRouter()
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
		log.Println("Старт сервера, порт " + s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка запуска сервера: %v", err)
		}
	}()
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Остановка сервера")
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка остановки сервера: %v", err)
		return err
	}
	log.Println("Сервер остановлен!")
	return nil
}
