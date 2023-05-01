package handler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/internal/grpc-storage/storage"
	"itisadb/internal/grpc-storage/usecase"
	api "itisadb/pkg/api/storage"
)

type Handler struct {
	api.UnimplementedStorageServer
	logic *usecase.UseCase
}

func New(logic *usecase.UseCase) *Handler {
	return &Handler{logic: logic}
}

func (h *Handler) Set(ctx context.Context, r *api.SetRequest) (*api.SetResponse, error) {
	memUsage, err := h.logic.Set(r.Key, r.Value, r.Unique)
	if errors.Is(err, storage.ErrAlreadyExists) {
		return &api.SetResponse{
			Total:     memUsage.Total,
			Available: memUsage.Available,
		}, status.Error(codes.AlreadyExists, err.Error())
	}

	return &api.SetResponse{
		Total:     memUsage.Total,
		Available: memUsage.Available,
	}, err
}

func (h *Handler) Get(ctx context.Context, r *api.GetRequest) (*api.GetResponse, error) {
	ram, value, err := h.logic.Get(r.Key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.GetResponse{
				Available: ram.Available,
				Total:     ram.Total,
				Value:     err.Error(),
			}, status.Error(codes.NotFound, err.Error())
		}
		return &api.GetResponse{
			Available: ram.Available,
			Total:     ram.Total,
			Value:     err.Error(),
		}, err
	}

	return &api.GetResponse{
		Available: ram.Available,
		Total:     ram.Total,
		Value:     value,
	}, nil
}

func (h *Handler) SetToIndex(ctx context.Context, r *api.SetToIndexRequest) (*api.SetResponse, error) {
	memUsage, err := h.logic.SetToIndex(r.Name, r.Key, r.Value)
	if errors.Is(err, storage.ErrAlreadyExists) {
		return &api.SetResponse{
			Total:     memUsage.Total,
			Available: memUsage.Available,
		}, status.Error(codes.AlreadyExists, err.Error())
	}

	return &api.SetResponse{
		Total:     memUsage.Total,
		Available: memUsage.Available,
	}, err
}

func (h *Handler) GetFromIndex(ctx context.Context, r *api.GetFromIndexRequest) (*api.GetResponse, error) {
	ram, value, err := h.logic.GetFromIndex(r.Name, r.Key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.GetResponse{
				Available: ram.Available,
				Total:     ram.Total,
			}, status.Error(codes.NotFound, err.Error())
		}
		// TODO: handle one more error here
		return &api.GetResponse{
			Available: ram.Available,
			Total:     ram.Total,
		}, err
	}

	return &api.GetResponse{
		Available: ram.Available,
		Total:     ram.Total,
		Value:     value,
	}, nil
}

func (h *Handler) GetIndex(ctx context.Context, r *api.GetIndexRequest) (*api.GetIndexResponse, error) {
	ram, index, err := h.logic.GetIndex(r.Name)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.GetIndexResponse{
				Available: ram.Available,
				Total:     ram.Total,
			}, status.Error(codes.NotFound, err.Error())
		}
		return &api.GetIndexResponse{
			Available: ram.Available,
			Total:     ram.Total,
		}, err
	}
	return &api.GetIndexResponse{
		Available: ram.Available,
		Total:     ram.Total,
		Index:     index,
	}, nil
}

func (h *Handler) NewIndex(ctx context.Context, r *api.NewIndexRequest) (*api.NewIndexResponse, error) {
	_, err := h.logic.NewIndex(r.Name)

	// TODO: handle ram
	if err != nil {
		return &api.NewIndexResponse{}, err
	}
	return &api.NewIndexResponse{}, nil
}

func (h *Handler) Size(ctx context.Context, r *api.IndexSizeRequest) (*api.IndexSizeResponse, error) {
	ram, size, err := h.logic.Size(r.Name)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.IndexSizeResponse{
				Available: ram.Available,
				Total:     ram.Total,
			}, status.Error(codes.NotFound, storage.ErrNotFound.Error())
		}
		return &api.IndexSizeResponse{
			Available: ram.Available,
			Total:     ram.Total,
		}, err
	}
	return &api.IndexSizeResponse{
		Available: ram.Available,
		Total:     ram.Total,
		Size:      size,
	}, nil
}
