package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"itisadb/internal/constants"
)

func (h *Handler) AuthMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	h.logger.Debug("Request", zap.String("method", info.FullMethod))

	if !h.security.MandatoryAuthorization || info.FullMethod == "/api.ItisaDB/Authenticate" {
		return handler(ctx, req)
	}

	token, err := getToken(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := h.session.AuthByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, constants.UserKey, claims)

	res, err := handler(ctx, req)
	if err != nil {
		h.logger.Error("Failed to perform request", zap.String("method", info.FullMethod), zap.Error(err))
		return nil, err
	}

	return res, nil
}
