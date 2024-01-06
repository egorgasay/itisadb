package logic

import (
	"github.com/egorgasay/gost"
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

func (l *Logic) hasPermission(claimsOpt gost.Option[models.UserClaims], level models.Level) bool {
	// always ok when security is disabled
	if !l.cfg.Security.On {
		return true
	}

	// ok when security is not mandatory for Default level
	if !l.cfg.Security.MandatoryAuthorization && level == constants.DefaultLevel {
		return true
	}

	if claimsOpt.IsNone() {
		return false
	}

	claims := claimsOpt.Unwrap()

	return claims.Level >= level
}
