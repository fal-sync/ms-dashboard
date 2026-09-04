package httpdelivery

import (
	"context"
	"net/http"

	"boilerplate-skeletoncode/internal/domain"
	"boilerplate-skeletoncode/internal/usecase"
)

type UserUsecase interface {
	CreateUser(rctx context.Context, input usecase.CreateUserInput) (domain.User, error)
	GetUser(rctx context.Context, id string) (domain.User, error)
	ListUsers(rctx context.Context) ([]domain.User, error)
}

func NewRouter(userUsecase UserUsecase) http.Handler {
	userHandler := newUserHandler(userUsecase)
	healthHandler := newHealthHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler.check)
	mux.HandleFunc("POST /users", userHandler.create)
	mux.HandleFunc("GET /users", userHandler.list)
	mux.HandleFunc("GET /users/{id}", userHandler.getByID)

	return mux
}
