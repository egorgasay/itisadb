package security

import (
	"github.com/egorgasay/gost"
	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

type SecurityService struct {
	cfg config.SecurityConfig
}

func NewSecurityService(cfg config.SecurityConfig) *SecurityService {
	return &SecurityService{
		cfg: cfg,
	}
}

func (l *SecurityService) HasPermission(claimsOpt gost.Option[models.UserClaims], level models.Level) bool {
	// always ok when security is disabled
	if !l.cfg.MandatoryAuthorization {
		return true
	}

	// ok when security is not mandatory for Default level
	if level == constants.DefaultLevel {
		return true
	}

	if claimsOpt.IsNone() {
		return false
	}

	claims := claimsOpt.Unwrap()

	return claims.Level >= level
}
