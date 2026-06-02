package server

import (
	"currency-exchange-converter/internal/config"
	"fmt"
	"log/slog"
	"net/http"
)

type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	http   *http.Server
}

func New(cfg *config.Config, logger *slog.Logger, router http.Handler) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
		http: &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.App.Port),
			Handler:           router,
			ReadTimeout:       cfg.Server.ReadTimeout,
			WriteTimeout:      cfg.Server.WriteTimeout,
			IdleTimeout:       cfg.Server.IdleTimeout,
			ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
			MaxHeaderBytes:    cfg.Server.MaxHeaderBytes,
		},
	}
}

func (s *Server) Run() error {
	s.logger.Info("server started",
		"port", s.cfg.App.Port,
		"log_level", s.cfg.Log.Level,
		"sqlite_path", s.cfg.SQLite.Path,
	)
	return s.http.ListenAndServe()
}
