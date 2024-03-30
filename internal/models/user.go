package models

type User struct {
	ID       int      `json:"id"`
	Login    string   `json:"username"`
	Password string   `json:"password"`
	Level    Level    `json:"level"`
	Active   bool     `json:"active"`
	changeID changeID `json:"-"`
}

type changeID struct {
	id uint64
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
	return u.changeID.id
}

func (u *User) SetChangeID(syncID uint64) {
	u.changeID.id = syncID
}
