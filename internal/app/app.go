package app

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"boilerplate-skeletoncode/internal/config"
	httpdelivery "boilerplate-skeletoncode/internal/delivery/http"
	"boilerplate-skeletoncode/internal/infrastructure/persistence/memory"
	"boilerplate-skeletoncode/internal/usecase"
)

type App struct {
	server          *http.Server
	shutdownTimeout time.Duration
	closers         []io.Closer
}

func New(cfg config.Config) *App {
	userRepository := memory.NewUserRepository()
	userCreatedHooks, closers := buildUserCreatedHooks(cfg)
	userUsecase := usecase.NewUserService(userRepository, userCreatedHooks...)
	router := httpdelivery.NewRouter(userUsecase)

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		server:          server,
		shutdownTimeout: cfg.ShutdownTimeout,
		closers:         closers,
	}
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, a.shutdownTimeout)
	defer cancel()

	errs := make([]error, 0, len(a.closers)+1)

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		errs = append(errs, err)
	}

	for _, closer := range a.closers {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
