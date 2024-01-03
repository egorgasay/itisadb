package balancer

import (
	"context"
	"fmt"
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

func (c *Balancer) CreateUser(ctx context.Context, userID int, user models.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// TODO: determine server?

	if !c.hasPermission(userID, user.Level) {
		return constants.ErrForbidden
	}

	user.Active = true
	_, err := c.storage.CreateUser(user)
	if err != nil {
		c.logger.Warn("failed to create user", zap.Error(err), zap.String("user", user.Login))
		return err
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteCreateUser(user)
	}

	return nil
}

func (c *Balancer) DeleteUser(ctx context.Context, userID int, login string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// TODO: determine server

	targetID, targetUser, err := c.storage.GetUserByName(login)
	if err != nil {
		c.logger.Warn("failed to get user for delete", zap.Error(err), zap.String("user", login))
		return err
	}

	if !c.hasPermission(userID, targetUser.Level) {
		return constants.ErrForbidden
	}

	err = c.storage.DeleteUser(targetID)
	if err != nil {
		c.logger.Warn("failed to delete user", zap.Error(err), zap.String("user", login))
		return err
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteDeleteUser(login)
	}

	return nil
}

func (c *Balancer) ChangePassword(ctx context.Context, userID int, login, password string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// TODO: determine server

	targetID, targetUser, err := c.storage.GetUserByName(login)
	if err != nil {
		c.logger.Warn("failed to get user for change password", zap.Error(err), zap.String("user", login))
		return err
	}

	if !c.hasPermission(userID, targetUser.Level) {
		return constants.ErrForbidden
	}

	targetUser.Password = password
	err = c.storage.SaveUser(targetID, targetUser)
	if err != nil {
		c.logger.Warn("failed to change password", zap.Error(err), zap.String("user", login))
		return err
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteCreateUser(targetUser)
	}

	return nil
}

func (c *Balancer) ChangeLevel(ctx context.Context, userID int, login string, level models.Level) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// TODO: determine server

	if !c.hasPermission(userID, level) {
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

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteCreateUser(targetUser)
	}

	return nil
}

func (c *Balancer) hasPermission(userID int, level models.Level) bool {
	// always ok when security is disabled
	if !c.cfg.Security.On {
		return true
	}

	// ok when security is not mandatory for Default level
	if !c.cfg.Security.MandatoryAuthorization && level == constants.DefaultLevel {
		return true
	}

	userLevel, err := c.storage.GetUserLevel(userID)
	if err != nil {
		c.logger.Warn("failed to get user level", zap.Error(err))
		return false
	}

	return userLevel >= level
}
