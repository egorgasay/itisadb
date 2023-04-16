package handler

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-storage/internal/grpc-storage/storage"
	"grpc-storage/internal/grpc-storage/usecase"
	api "grpc-storage/pkg/api/storage"
)

type Handler struct {
	api.UnimplementedStorageServer
	logic *usecase.UseCase
}

func New(logic *usecase.UseCase) *Handler {
	return &Handler{logic: logic}
}

func (h *Handler) Set(s api.Storage_SetServer) error {
	for {
		select {
		case <-s.Context().Done():
			return s.Context().Err()
		default:
			recv, err := s.Recv()
			if err != nil {
				return err
			}

			memUsage := h.logic.Set(recv.Key, recv.Value)
			msg := &api.SetResponse{
				Status:    "ok",
				Total:     memUsage.Total,
				Available: memUsage.Available,
			}

			err = s.Send(msg)
			if err != nil {
				return err
			}
		}
	}
}

func (h *Handler) Get(s api.Storage_GetServer) error {
	for {
		select {
		case <-s.Context().Done():
			return s.Context().Err()

		default:
			recv, err := s.Recv()
			if err != nil {
				return err
			}
			var msg *api.GetResponse
			ram, value, err := h.logic.Get(recv.Key)
			if err != nil {
				if errors.Is(err, storage.ErrNotFound) {
					msg = &api.GetResponse{
						Available: ram.Available,
						Total:     ram.Total,
						Value:     status.Error(codes.NotFound, err.Error()).Error(),
					}
				}
				msg = &api.GetResponse{
					Available: ram.Available,
					Total:     ram.Total,
					Value:     err.Error(),
				}
			} else {
				msg = &api.GetResponse{
					Available: ram.Available,
					Total:     ram.Total,
					Value:     value,
				}
			}

			err = s.Send(msg)
			if err != nil {
				return err
			}
		}
	}
}
