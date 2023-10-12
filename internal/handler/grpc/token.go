package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func getToken(ctx context.Context) (token string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}

	values := md.Get("token")
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "no tokens in token")
	}

	return values[0], nil
}
