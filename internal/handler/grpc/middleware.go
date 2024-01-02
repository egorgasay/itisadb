package grpc

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func (h *Handler) AuthMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	h.logger.Debug("Request", zap.String("method", info.FullMethod))

	if !h.conf.On || info.FullMethod == "/api.ItisaDB/Authenticate" {
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
