package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreyxaxa/PubSub_gRPC_Service/config"
	"github.com/andreyxaxa/PubSub_gRPC_Service/internal/controller/grpc"
	"github.com/andreyxaxa/PubSub_gRPC_Service/pkg/grpcserver"
	"github.com/andreyxaxa/PubSub_gRPC_Service/pkg/logger"
	"github.com/andreyxaxa/PubSub_gRPC_Service/pkg/subpub"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// SubPub
	sp := subpub.NewSubPub()

	// gRPC Server
	grpcServer := grpcserver.New(grpcserver.Port(cfg.GRPC.Port))
	grpc.NewRouter(grpcServer.App, sp, l) // DI

	// Start server
	go grpcServer.Start()
	l.Info("gRPC server started on %s", cfg.GRPC.Port)

	// Graceful Shutdown
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err := <-grpcServer.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	// Shutdown
	err := grpcServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	} else {
		l.Info("gRPC server shutdown complete")
	}
}
