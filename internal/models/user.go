package models

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"username"`
	Password string `json:"password"`
	Level    Level  `json:"level"`
	Active   bool   `json:"active"`
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
