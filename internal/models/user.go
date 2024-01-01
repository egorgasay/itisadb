package models

type User struct {
	ID       uint   `json:"id"`
	Login    string `json:"username"`
	Password string `json:"password"`
	Level    Level  `json:"level"`
	Active   bool   `json:"active"`
}

type MetaData struct {
	UserID   uint  `json:"user_id"`
	ServerID int32 `json:"server_id"`
}
