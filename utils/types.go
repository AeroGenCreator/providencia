package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserParty string `json:"party"`
	jwt.RegisteredClaims
}

type LevelConfig struct {
	Level       int      `json:"level"`
	Permissions []string `json:"permissions"`
}

type SaavedraConfigFile struct {
	Brand   map[string]any `json:"brand"`
	Parties [4]string      `json:"parties"`
	Levels  []LevelConfig  `json:"levels"`
}

type ContextKey string

var NoUserPartyPermission = errors.New("Permissions were added to an endpoint but no user party data provided.")
var DuplicatedDataError = errors.New("Duplicated data, check information | Datos duplicados, revisar información")
