package models

type User struct {
	Login    string `json:"username"`
	Password string `json:"password"`
	Level    Level  `json:"level"`
}
