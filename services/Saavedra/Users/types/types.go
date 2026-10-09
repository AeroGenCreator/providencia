package types

import (
	"errors"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Party    string `json:"party"`
}

type UserStr struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Party    string `json:"party"`
}

type UserSlice struct {
	Records     []*User `json:"records"`
	HasNextPage bool    `json:"hasNextPage"`
}

type Many2One struct {
	Parties [4]string `json:"parties"`
}

type Body struct {
	Id int `json:"id"`
}

var ErrAdminProtection = errors.New("admin protection: standard admin user cannot be deleted")
var ErrSelfProtection = errors.New("self account protection: can't delete this account since is the current logged account")
