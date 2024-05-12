package domains

import (
	"github.com/egorgasay/gost"
	"itisadb/internal/models"
)

type SecurityService interface {
	HasPermission(claimsOpt gost.Option[models.UserClaims], level models.Level) bool
	Encrypt(val string) (string, error)
	Decrypt(val string) (string, error)
}
