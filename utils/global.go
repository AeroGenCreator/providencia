package utils

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
)

var SessionToken string
var DatabaseName string
var AdminName string
var AdminEmail string
var AdminPassword string
var RecordsPerSlice int
var IsProduction bool
var Host string
var Port string
var AdminId int
var Parties [4]string
var Levels []LevelConfig

func InstantiateEnvironmentVariables() {

	// Carga de cadenas de texto
	SessionToken = os.Getenv("SESSION_TOKEN")
	DatabaseName = os.Getenv("DATABASE_NAME")
	AdminName = os.Getenv("ADMIN_NAME")
	AdminEmail = os.Getenv("ADMIN_EMAIL")
	AdminPassword = os.Getenv("ADMIN_PASSWORD")
	Host = os.Getenv("HOST")
	Port = os.Getenv("PORT")

	// Validar datos de funcionamiento
	if SessionToken == "" || AdminName == "" || AdminEmail == "" || AdminPassword == "" {
		log.Fatalln(`
			Missing critical data; one of the following data is missing:
			SESSION_TOKEN, ADMIN_NAME, ADMIN_EMAIL, ADMIN_PASSWORD`)
	}

	// Garantizar opcionales de servidor
	if DatabaseName == "" {
		DatabaseName = "database"
	}
	if Host == "" {
		Host = "0.0.0.0"
	}
	if Port == "" {
		Port = "8080"
	}

	// Obtencion de variables booleanas
	IsProduction, err := strconv.ParseBool(os.Getenv("IS_PRODUCTION"))
	if err != nil {
		log.Print("Invalid datatype given to IS_PRODUCTION: Default 'false'.")
		IsProduction = false
	}
	fmt.Println("🔐 Project status set to: ", IsProduction)

	// Obtencion de variables numericas
	RecordsPerSlice, err = strconv.Atoi(os.Getenv("RECORDS_PER_SLICE"))
	if err != nil {
		log.Print("Invalid datatype given to RECORDS_PER_SLICE: Default: '10 records per slice'")
		RecordsPerSlice = 10
	}
	fmt.Println("📦 Records per slice set to: ", RecordsPerSlice)

}

func AddPermissions(level int, userParty string) (bool, error) {
	// Si no se dan permisos se asume el nivel mas alto de estos.
	if level == 0 {
		level = 4
	}
	// Olvidar pasar grupo(party) de permisos de usuario petición, se notifica.
	if userParty == "" {
		return false, NoUserPartyPermission
	}
	// Se valida permiso de usuario desde petición. Si el grupo (party) es valido devuelve 'true'
	for _, dict := range Levels {
		if dict.Level == level {
			hasPermission := slices.Contains(dict.Permissions, userParty)
			return hasPermission, nil
		}
	}
	// Si no se dan permisos siempre se asume mayor protección.
	return false, nil
}
