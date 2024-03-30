package models

import "github.com/egorgasay/gost"

type User struct {
	ID       int                 `json:"id"`
	Login    string              `json:"username"`
	Password string              `json:"password"`
	Level    Level               `json:"level"`
	Active   bool                `json:"active"`
	changeID gost.RwLock[uint64] `json:"-"`
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

func (u *User) GetChangeID() uint64 {
	u.changeID.RBorrow()
	return u.changeID.ReadAndRelease()
}

func (u *User) SetChangeID(syncID uint64) {
	u.changeID.SetWithLock(syncID)
}
