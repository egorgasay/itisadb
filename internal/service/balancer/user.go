package balancer

import (
	"context"
	"fmt"

	"github.com/egorgasay/gost"
	"go.uber.org/zap"
	"itisadb/internal/constants"
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

	if !c.security.HasPermission(claims, user.Level) {
		return constants.ErrForbidden
	}

	user.Active = true
	r := c.storage.NewUser(user)
	if r.IsErr() {
		c.logger.Warn("failed to create user", zap.Error(r.Error()), zap.String("user", user.Login))
		return r.Error()
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteNewUser(user)
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

	r := c.storage.GetUserByName(login)
	if r.IsErr() {
		c.logger.Warn("failed to get user for delete", zap.Error(r.Error()), zap.String("user", login))
		return r.Error()
	}

	targetUser := r.Unwrap()

	if !c.security.HasPermission(claims, targetUser.Level) {
		return constants.ErrForbidden
	}

	if err := c.storage.DeleteUser(targetUser.ID).Error(); err != nil {
		c.logger.Warn("failed to delete user", zap.Error(err), zap.String("user", login))
		return err
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteDeleteUser(login)
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

	res := c.storage.GetUserByName(login)
	if res.IsErr() {
		c.logger.Warn("failed to get user for change password", zap.Error(res.Error()), zap.String("user", login))
		return res.Error()
	}

	targetUser := res.Unwrap()

	if !c.security.HasPermission(claims, targetUser.Level) {
		return constants.ErrForbidden
	}

	targetUser.Password = password
	if r := c.storage.SaveUser(targetUser); r.IsErr() {
		c.logger.Warn("failed to change level", zap.Error(r.Error()), zap.String("user", login))
		return r.Error()
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteNewUser(targetUser)
	}

	if err := c.servers.Iter(func(server domains.Server) error {
		ctx := context.TODO()

		targetUser.Password = password
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

	if !c.security.HasPermission(claims, level) {
		return constants.ErrForbidden
	}

	res := c.storage.GetUserByName(login)
	if res.IsErr() {
		c.logger.Warn("failed to get user for change level", zap.Error(res.Error()), zap.String("user", login))
		return res.Error()
	}

	targetUser := res.Unwrap()
	targetUser.Level = level

	if err := c.storage.SaveUser(targetUser).Error(); err != nil {
		c.logger.Warn("failed to change level", zap.Error(err), zap.String("user", login))
		return err
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteNewUser(targetUser)
	}

	if err := c.servers.Iter(func(server domains.Server) error {
		ctx := context.TODO()

		targetUser.Level = level
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
