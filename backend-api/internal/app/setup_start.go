package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func (srv *Service) StartServe() error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.startHTTPServer(ctx); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.startGRPCServer(ctx); err != nil {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Printf("shutdown signal received")
	case err := <-errChan:
		log.Printf("server error: %v", err)
		stop()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("error shutting down HTTP server: %v", err)
	} else {
		log.Printf("http server stopped")
	}

	srv.grpcServer.GracefulStop()
	log.Printf("grpc server stopped")

	if srv.db != nil {
		if err := srv.db.Close(); err != nil {
			log.Printf("error closing db: %v", err)
		} else {
			log.Printf("db connection closed")
		}
	}

	wg.Wait()
	log.Printf("all servers exited gracefully")

	return nil
}

func (srv *Service) startHTTPServer(ctx context.Context) error {
	if srv.httpRouter == nil {
		return fmt.Errorf("router is not initialized, call SetupAPIHandlers first")
	}

	addr := fmt.Sprintf("%s:%d", srv.config.HTTPServer.Host, srv.config.HTTPServer.Port)
	srv.httpServer = &http.Server{
		Addr:    addr,
		Handler: srv.httpRouter,
	}

	log.Printf("starting http server on %s", addr)
	return srv.httpServer.ListenAndServe()
}

func (srv *Service) startGRPCServer(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", srv.config.GRPCServer.Host, srv.config.GRPCServer.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	log.Printf("starting grpc server on %s", addr)
	return srv.grpcServer.Serve(lis)
}
