package core

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

func (c *Core) Authenticate(ctx context.Context, login string, password string) (string, error) {
	token, err := c.session.AuthByPassword(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate: %w", err)
	}

	return token, nil
}

func (c *Core) CreateUser(ctx context.Context, userID int, user models.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if user.Level > constants.SecretLevel || user.Level < constants.DefaultLevel {
		return constants.ErrForbidden
	}

	creator, err := c.storage.GetUserByID(int(userID))
	if err != nil {
		c.logger.Warn("failed to get user for delete", zap.Error(err), zap.String("user", user.Login))
		return err
	}

	if creator.Level < user.Level {
		c.logger.Warn("can't create user", zap.Error(err), zap.String("user", user.Login))
		return constants.ErrForbidden
	}

	_, err = c.storage.CreateUser(user)
	if err != nil {
		c.logger.Warn("failed to create user", zap.Error(err), zap.String("user", user.Login))
		return err
	}

	if c.cfg.TransactionLoggerConfig.On {
		c.tlogger.WriteCreateUser(user)
	}

	return nil
}

func (c *Core) DeleteUser(ctx context.Context, userID int, login string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	targetID, targetUser, err := c.storage.GetUserByName(login)
	if err != nil {
		c.logger.Warn("failed to get user for delete", zap.Error(err), zap.String("user", login))
		return err
	}

	us, err := c.storage.GetUserByID(int(userID))
	if err != nil {
		c.logger.Warn("failed to get user for delete", zap.Error(err), zap.String("user", login))
		return err
	}

	if targetUser.Level > us.Level {
		c.logger.Warn("can't delete user", zap.Error(err), zap.String("user", login))
		return constants.ErrForbidden
	}

	err = c.storage.DeleteUser(targetID)
	if err != nil {
		c.logger.Warn("failed to delete user", zap.Error(err), zap.String("user", login))
		return err
	}

	if c.cfg.TransactionLoggerConfig.On {
		c.tlogger.WriteDeleteUser(login)
	}

	return nil
}

func (c *Core) ChangePassword(ctx context.Context, userID int, login, password string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	targetID, targetUser, err := c.storage.GetUserByName(login)
	if err != nil {
		c.logger.Warn("failed to get user for change password", zap.Error(err), zap.String("user", login))
		return err
	}

	us, err := c.storage.GetUserByID(userID)
	if err != nil {
		c.logger.Warn("failed to get user for change password", zap.Error(err), zap.String("user", login))
		return err
	}

	if targetUser.Level > us.Level {
		c.logger.Warn("can't change password", zap.Error(err), zap.String("user", login))
		return constants.ErrForbidden
	}

	targetUser.Password = password
	err = c.storage.SaveUser(targetID, targetUser)
	if err != nil {
		c.logger.Warn("failed to change password", zap.Error(err), zap.String("user", login))
		return err
	}

	if c.cfg.TransactionLoggerConfig.On {
		c.tlogger.WriteCreateUser(targetUser)
	}

	return nil
}

func (c *Core) ChangeLevel(ctx context.Context, userID int, login string, level int8) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	targetID, targetUser, err := c.storage.GetUserByName(login)
	if err != nil {
		c.logger.Warn("failed to get user for change level", zap.Error(err), zap.String("user", login))
		return err
	}

	us, err := c.storage.GetUserByID(int(userID))
	if err != nil {
		c.logger.Warn("failed to get user for change level", zap.Error(err), zap.String("user", login))
		return err
	}

	if targetUser.Level > us.Level {
		c.logger.Warn("can't change level", zap.Error(err), zap.String("user", login))
		return constants.ErrForbidden
	}

	targetUser.Level = models.Level(level)
	err = c.storage.SaveUser(targetID, targetUser)
	if err != nil {
		c.logger.Warn("failed to change level", zap.Error(err), zap.String("user", login))
		return err
	}

	if c.cfg.TransactionLoggerConfig.On {
		c.tlogger.WriteCreateUser(targetUser)
	}

	return nil
}
