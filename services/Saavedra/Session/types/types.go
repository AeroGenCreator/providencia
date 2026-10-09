package types

type Session struct {
	Id     int    `json:"id"`
	UserId int    `json:"userId"`
	Token  string `json:"token"`
}
