package server

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type ServerHTTP struct {
	server *http.Server
	Logger *zap.Logger
}

func NewServerHTTP(logger *zap.Logger, handler http.Handler, adr string) *ServerHTTP {
	return &ServerHTTP{

		server: &http.Server{
			Addr:    adr,
			Handler: handler,
		},

		Logger: logger,
	}
}

func (sh *ServerHTTP) Start() error {
	sh.Logger.Info("Starting HTTP Server")

	err := http.ListenAndServe(sh.server.Addr, sh.server.Handler)
	if err != nil {
		sh.Logger.Error("Error starting HTTP server", zap.Error(err))
		return err
	}

	sh.Logger.Info("The server has been started successfully.")
	return nil
}

func (sh *ServerHTTP) Stop(ctx context.Context) error {
	sh.Logger.Info("Shutting down the server...")
	err := sh.server.Shutdown(ctx)
	if err != nil {
		sh.Logger.Error("Error shutting down the server: ", zap.Error(err))
		return err
	}
	sh.Logger.Info("Server has been shut down successfully")
	return nil
}

func (sh *ServerHTTP) Restart(ctx context.Context) error {

	if err := sh.Stop(ctx); err != nil {
		sh.Logger.Error("Error shutting down the server: ", zap.Error(err))
		return err
	}

	sh.server = &http.Server{
		Addr:    sh.server.Addr,
		Handler: sh.server.Handler,
	}

	go func() {
		err := sh.Start()
		if err != nil && err != http.ErrServerClosed {
			sh.Logger.Error("Could not listen on %s: %v\n", zap.String("address", sh.server.Addr), zap.Error(err))
		}
	}()

	return nil
}
