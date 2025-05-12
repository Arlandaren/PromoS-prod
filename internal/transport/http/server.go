package server

import (
	"context"
	"solution/internal/shared/config"
)

type Server struct {
	addr         string
	serverRouter *MainRouter
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		addr:         cfg.Server.Addr,
		serverRouter: NewRouter(),
	}
}

func (s *Server) StartHttpServer(ctx context.Context) error {
	s.serverRouter.SetContext(ctx)

	s.serverRouter.RouteInit()

	err := s.serverRouter.router.Run(s.addr)
	if err != nil {
		return err
	}
	return nil
}
