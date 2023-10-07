package grpc

import (
	"context"
	"google.golang.org/grpc"
)

func (h *Handler) AuthMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/api.ItisaDB/Authenticate" {
		return handler(ctx, req)
	}

	_, err := getToken(ctx)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}
