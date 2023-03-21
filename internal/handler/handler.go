package handler

import (
	"context"
	"github.com/egorgasay/grpc-storage/internal/usecase"
	api "github.com/egorgasay/grpc-storage/pkg/api"
)

type Handler struct {
	api.UnimplementedStorageServer
	logic *usecase.UseCase
}

func (h *Handler) Set(ctx context.Context, r *api.SetRequest) (*api.SetResponse, error) {
	h.logic.Set(r.Key, r.Value)
	return &api.SetResponse{
		Status: "ok",
	}, nil
}

func (h *Handler) Get(ctx context.Context, r *api.GetRequest) (*api.GetResponse, error) {
	value, err := h.logic.Get(r.Key)
	if err != nil {
		return &api.GetResponse{
			Value: err.Error(),
		}, err
	}

	return &api.GetResponse{
		Value: value,
	}, nil
}

func New(logic *usecase.UseCase) *Handler {
	return &Handler{logic: logic}
}
