package handler

import "github.com/egorgasay/grpc-storage/internal/usecase"

type Handler struct {
	logic *usecase.UseCase
}

func New(logic *usecase.UseCase) *Handler {
	return &Handler{logic: logic}
}
