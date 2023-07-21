package rest

import (
	"github.com/testit-tms/webhook-bot/internal/config"
	"golang.org/x/exp/slog"
)

type Server struct {
	cfg config.HTTPServer
	log *slog.Logger
}

func New(log *slog.Logger, cfg config.HTTPServer) *Server {
	return &Server{
		log: log,
		cfg: cfg,
	}
}
