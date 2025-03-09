package server

import "context"

type Server struct {
	log Logger
}

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Application interface { 
	// TODO
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		log: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}
