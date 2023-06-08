package handler

import "itisadb/internal/dirdb/usecase"

type Handler struct {
	logic *usecase.IUseCase
}

func New(logic *usecase.IUseCase) *Handler {
	return &Handler{logic: logic}
}
