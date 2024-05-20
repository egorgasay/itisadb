package constants

import "time"

const (

	// AccessTTL TODO: handle expiration of tokens
	AccessTTL  = 24 * time.Hour
	RefreshTTL = 7 * 24 * time.Hour
)

const (
	GUID  = "guid"
	IAT   = "iat"
	LEVEL = "level"
)

const NoUser = 0

const UserKey = "user-claims"
