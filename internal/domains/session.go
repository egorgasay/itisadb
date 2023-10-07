package domains

import (
	"context"
)

type Session interface {
	AuthByToken(ctx context.Context, token string) (string, error)
	AuthByPassword(ctx context.Context, username, password string) (string, error)
	Create(ctx context.Context, guid int) (string, error)
}
