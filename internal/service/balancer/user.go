package balancer

import (
	"context"
	"fmt"

	"github.com/egorgasay/gost"
	"go.uber.org/zap"
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

func (c *Balancer) Authenticate(ctx context.Context, login string, password string) (string, error) {
	token, err := c.session.AuthByPassword(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate: %w", err)
	}

	return token, nil
}

func (c *Balancer) CreateUser(ctx context.Context, claims gost.Option[models.UserClaims], user models.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// TODO: determine server?

	if !c.security.HasPermission(claims, user.Level) {
		return constants.ErrForbidden
	}

	user.Active = true
	_, err := c.storage.CreateUser(user)
	if err != nil {
		c.logger.Warn("failed to create user", zap.Error(err), zap.String("user", user.Login))
		return err
	}

	if c.cfg.TransactionLogger.On { // TODO: ???
		c.tlogger.WriteCreateUser(user)
	}

	return nil
}

func (c *Balancer) DeleteUser(ctx context.Context, claims gost.Option[models.UserClaims], login string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// TODO: determine server

	targetID, targetUser, err := c.storage.GetUserByName(login)
	if err != nil {
		c.logger.Warn("failed to get user for delete", zap.Error(err), zap.String("user", login))
		return err
	}

	if !c.security.HasPermission(claims, targetUser.Level) {
		return constants.ErrForbidden
	}

	err = c.storage.DeleteUser(targetID)
	if err != nil {
		c.logger.Warn("failed to delete user", zap.Error(err), zap.String("user", login))
		return err
	}

	if c.cfg.TransactionLogger.On { // TODO: ???
		c.tlogger.WriteDeleteUser(login)
	}

	return nil
}

func (c *Balancer) ChangePassword(ctx context.Context, claims gost.Option[models.UserClaims], login, password string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// TODO: determine server

	targetID, targetUser, err := c.storage.GetUserByName(login)
	if err != nil {
		c.logger.Warn("failed to get user for change password", zap.Error(err), zap.String("user", login))
		return err
	}

	if !c.security.HasPermission(claims, targetUser.Level) {
		return constants.ErrForbidden
	}

	targetUser.Password = password
	err = c.storage.SaveUser(targetID, targetUser)
	if err != nil {
		c.logger.Warn("failed to change password", zap.Error(err), zap.String("user", login))
		return err
	}

	if c.cfg.TransactionLogger.On { // TODO: ???
		c.tlogger.WriteCreateUser(targetUser)
	}

	return nil
}

func (c *Balancer) ChangeLevel(ctx context.Context, claims gost.Option[models.UserClaims], login string, level models.Level) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// TODO: determine server

	if !c.security.HasPermission(claims, level) {
		return constants.ErrForbidden
	}

	targetID, targetUser, err := c.storage.GetUserByName(login)
	if err != nil {
		c.logger.Warn("failed to get user for change level", zap.Error(err), zap.String("user", login))
		return err
	}

	targetUser.Level = level
	err = c.storage.SaveUser(targetID, targetUser)
	if err != nil {
		c.logger.Warn("failed to change level", zap.Error(err), zap.String("user", login))
		return err
	}

	if c.cfg.TransactionLogger.On { // TODO: ???
		c.tlogger.WriteCreateUser(targetUser)
	}

	return nil
}
