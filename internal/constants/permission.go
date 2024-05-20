package constants

import "itisadb/internal/models"

const (
	// DefaultLevel no mandatory encryption, auth
	DefaultLevel models.Level = iota
	// RestrictedLevel no mandatory encryption, but mandatory auth
	RestrictedLevel
	// SecretLevel mandatory encryption, Auth
	SecretLevel
)

const (
	MaxLevel = SecretLevel
	MinLevel = DefaultLevel
)
