package models

import "github.com/egorgasay/gost"

type User struct {
	ID       int                 `json:"id"`
	Login    string              `json:"username"`
	Password string              `json:"password"`
	Level    Level               `json:"level"`
	Active   bool                `json:"active"`
	syncID   gost.RwLock[uint64] `json:"-"`
}

func (u User) ExtractClaims() UserClaims {
	return UserClaims{
		ID:    u.ID,
		Level: u.Level,
	}
}

type UserClaims struct {
	ID    int   `json:"id"`
	Level Level `json:"level"`
}

func (u *User) GetSyncID() uint64 {
	u.syncID.RBorrow()
	return u.syncID.ReadAndRelease()
}

func (u *User) SetSyncID(syncID uint64) {
	u.syncID.SetWithLock(syncID)
}
