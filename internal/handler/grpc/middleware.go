package grpc

import (
	"context"
	"google.golang.org/grpc"
)

func (h *Handler) AuthMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/api.ItisaDB/Authenticate" {
		return handler(ctx, req)
	}

	token, err := getToken(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := h.session.AuthByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, "userID", uint(userID))

	return handler(ctx, req)
}
