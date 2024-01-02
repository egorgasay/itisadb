package logic

import (
	"go.uber.org/zap"
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

func (l *Logic) hasPermission(userID int, level models.Level) bool {
	// always ok when security is disabled
	if !l.cfg.Security.On {
		return true
	}

	// ok when security is not mandatory for Default level
	if !l.cfg.Security.MandatoryAuthorization && level == constants.DefaultLevel {
		return true
	}

	userLevel, err := l.storage.GetUserLevel(userID)
	if err != nil {
		l.logger.Warn("failed to get user level", zap.Error(err))
		return false
	}

	return userLevel >= level
}
