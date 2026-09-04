package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"boilerplate-skeletoncode/internal/app"
	"boilerplate-skeletoncode/internal/config"
)

func main() {
	cfg := config.Load()

	application := app.New(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)

	go func() {
		errCh <- application.Run()
	}()

	select {
	case <-ctx.Done():
		if err := application.Shutdown(context.Background()); err != nil {
			log.Fatalf("shutdown server: %v", err)
		}
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("run server: %v", err)
		}
	}
}
