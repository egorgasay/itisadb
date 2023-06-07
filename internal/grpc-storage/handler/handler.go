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
	logic usecase.IUseCase
}

func New(logic usecase.IUseCase) *Handler {
	return &Handler{logic: logic}
}

func (h *Handler) Set(ctx context.Context, r *api.SetRequest) (*api.SetResponse, error) {
	ram, err := h.logic.Set(r.Key, r.Value, r.Unique)
	if errors.Is(err, storage.ErrAlreadyExists) {
		return &api.SetResponse{
			Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
		}, status.Error(codes.AlreadyExists, err.Error())
	}

	return &api.SetResponse{
		Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
	}, err
}

func (h *Handler) Get(ctx context.Context, r *api.GetRequest) (*api.GetResponse, error) {
	ram, value, err := h.logic.Get(r.Key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.GetResponse{
				Ram:   &api.Ram{Total: ram.Total, Available: ram.Available},
				Value: err.Error(),
			}, status.Error(codes.NotFound, err.Error())
		}
		return &api.GetResponse{
			Ram:   &api.Ram{Total: ram.Total, Available: ram.Available},
			Value: err.Error(),
		}, err
	}

	return &api.GetResponse{
		Ram:   &api.Ram{Total: ram.Total, Available: ram.Available},
		Value: value,
	}, nil
}

func (h *Handler) SetToIndex(ctx context.Context, r *api.SetToIndexRequest) (*api.SetResponse, error) {
	ram, err := h.logic.SetToIndex(r.Name, r.Key, r.Value, r.Unique)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return &api.SetResponse{
				Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
			}, status.Error(codes.AlreadyExists, err.Error())
		}
	}

	return &api.SetResponse{
		Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
	}, err
}

func (h *Handler) GetFromIndex(ctx context.Context, r *api.GetFromIndexRequest) (*api.GetResponse, error) {
	ram, value, err := h.logic.GetFromIndex(r.Name, r.Key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.GetResponse{
				Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
			}, status.Error(codes.NotFound, err.Error())
		}
		// TODO: handle one more error here
		return &api.GetResponse{
			Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
		}, err
	}

	return &api.GetResponse{
		Ram:   &api.Ram{Total: ram.Total, Available: ram.Available},
		Value: value,
	}, nil
}

func (h *Handler) GetIndex(ctx context.Context, r *api.GetIndexRequest) (*api.GetIndexResponse, error) {
	ram, index, err := h.logic.GetIndex(r.Name)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.GetIndexResponse{
				Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
			}, status.Error(codes.NotFound, err.Error())
		}
		return &api.GetIndexResponse{
			Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
		}, err
	}
	return &api.GetIndexResponse{
		Ram:   &api.Ram{Total: ram.Total, Available: ram.Available},
		Index: index,
	}, nil
}

func (h *Handler) NewIndex(ctx context.Context, r *api.NewIndexRequest) (*api.NewIndexResponse, error) {
	ram, err := h.logic.NewIndex(r.Name)

	if err != nil {
		return &api.NewIndexResponse{}, err
	}
	return &api.NewIndexResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}, nil
}

func (h *Handler) AttachToIndex(ctx context.Context, r *api.AttachToIndexRequest) (*api.AttachToIndexResponse, error) {
	ram, err := h.logic.AttachToIndex(r.Dst, r.Src)
	if err != nil {
		return &api.AttachToIndexResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}, err
	}
	return &api.AttachToIndexResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}, nil
}

func (h *Handler) DeleteIndex(ctx context.Context, r *api.DeleteIndexRequest) (*api.DeleteIndexResponse, error) {
	ram, err := h.logic.DeleteIndex(r.Index)
	if err != nil {
		// TODO: handle error
		// TODO: useless &api.Ram{Total: ram.Total, Available: ram.Available} ??
		return &api.DeleteIndexResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}, err
	}
	return &api.DeleteIndexResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}, nil
}

func (h *Handler) Delete(ctx context.Context, r *api.DeleteRequest) (*api.DeleteResponse, error) {
	ram, err := h.logic.Delete(r.Key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.DeleteResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}, status.Error(codes.NotFound, storage.ErrNotFound.Error())
		}
		return &api.DeleteResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}, err
	}
	return &api.DeleteResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}, nil
}

func (h *Handler) Size(ctx context.Context, r *api.IndexSizeRequest) (*api.IndexSizeResponse, error) {
	ram, size, err := h.logic.Size(r.Name)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.IndexSizeResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}, status.Error(codes.NotFound, storage.ErrNotFound.Error())
		}
		return &api.IndexSizeResponse{
			Ram: &api.Ram{Total: ram.Total, Available: ram.Available},
		}, err
	}
	return &api.IndexSizeResponse{
		Ram:  &api.Ram{Total: ram.Total, Available: ram.Available},
		Size: size,
	}, nil
}

func (h *Handler) DeleteAttr(ctx context.Context, r *api.DeleteAttrRequest) (*api.DeleteAttrResponse, error) {
	ram, err := h.logic.DeleteAttr(r.Name, r.Key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &api.DeleteAttrResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}, status.Error(codes.NotFound, storage.ErrNotFound.Error())
		}
		return &api.DeleteAttrResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}, err
	}
	return &api.DeleteAttrResponse{Ram: &api.Ram{Total: ram.Total, Available: ram.Available}}, nil
}
