package balancer

import (
	"context"
	"fmt"

	"github.com/egorgasay/gost"
	"go.uber.org/zap"
	"itisadb/internal/domains"
	"itisadb/internal/models"
)

func (c *Balancer) Authenticate(ctx context.Context, login string, password string) (string, error) {
	token, err := c.session.AuthByPassword(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate: %w", err)
	}

	return token, nil
}

func (c *Balancer) NewUser(ctx context.Context, claims gost.Option[models.UserClaims], user models.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := c.servers.Iter(func(server domains.Server) error {
		if server.IsOffline() {
			return nil
		}

		ctx := context.TODO()

		if r := server.NewUser(ctx, claims, user); r.IsErr() {
			c.logger.Warn("failed to create user", zap.Error(r.Error()), zap.String("user", user.Login))
			return r.Error()
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c *Balancer) DeleteUser(ctx context.Context, claims gost.Option[models.UserClaims], login string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := c.servers.Iter(func(server domains.Server) error {
		ctx := context.TODO()

		if r := server.DeleteUser(ctx, claims, login); r.IsErr() {
			c.logger.Warn("failed to delete user", zap.Error(r.Error()), zap.String("user", login))
			return r.Error()
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c *Balancer) ChangePassword(ctx context.Context, claims gost.Option[models.UserClaims], login, password string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := c.servers.Iter(func(server domains.Server) error {
		ctx := context.TODO()

		if r := server.ChangePassword(ctx, claims, login, password); r.IsErr() {
			c.logger.Warn("failed to change password", zap.Error(r.Error()), zap.String("user", login))
			return r.Error()
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c *Balancer) ChangeLevel(ctx context.Context, claims gost.Option[models.UserClaims], login string, level models.Level) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := c.servers.Iter(func(server domains.Server) error {
		ctx := context.TODO()

		if err := server.ChangeLevel(ctx, claims, login, level).Error(); err != nil {
			c.logger.Warn("failed to change level", zap.Error(err), zap.String("user", login))
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
