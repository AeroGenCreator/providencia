package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	employeeRouter "saavedra/services/Employee/router"
	employeeSQL "saavedra/services/Employee/sql"
	homeRouter "saavedra/services/Saavedra/Home/router"
	loginRouter "saavedra/services/Saavedra/Login/router"
	serveAssets "saavedra/services/Saavedra/ServeFiles/router"
	sessionRouter "saavedra/services/Saavedra/Session/router"
	sessionSQL "saavedra/services/Saavedra/Session/sql"
	usersRouter "saavedra/services/Saavedra/Users/router"
	usersSQL "saavedra/services/Saavedra/Users/sql"
	utils "saavedra/utils"

	_ "github.com/glebarez/go-sqlite"
	"github.com/joho/godotenv"
)

func main() {

	// Leer variables desde fichero .env
	if err := godotenv.Load(); err != nil {
		log.Printf("No such file '.env' reading from system environment variables.")
	}

	// Carga las variables como constantes globales accesibles desde cualquier parte del proyecto.
	utils.InstantiateEnvironmentVariables()

	// Cargar metadata desde JSON saavedra configuración.
	utils.LoadSaavedraConfigFile()

	// Creación: directorio base de datos
	if err := os.MkdirAll("./db", 0755); err != nil {
		log.Fatalf("Error when creating database directory: (%v)", err)
	}

	// Instancia conexion a la base de datos.
	db, err := sql.Open("sqlite", "./db/"+utils.DatabaseName+".db")
	if err != nil {
		log.Fatalf("Error when opening database connection: (%v)", err)
	}

	// === Registro de Esquemas SQL (TABLAS) ===
	if err = usersSQL.CreateSchema(db); err != nil {
		log.Fatal(err.Error())
	}
	if err = sessionSQL.CreateSchema(db); err != nil {
		log.Fatal(err.Error())
	}
	if err = employeeSQL.CreateSchema(db); err != nil {
		log.Fatal(err.Error())
	}

	// Crear usuario admin o remplazar credenciales.
	utils.InsertAdmin(db)

	// Instanciar servicio GO
	mux := http.NewServeMux()

	// Registrar directorio de assets
	serveAssets.ServeGlobalAssetsDirectory(mux)

	// === Registro de Servicios ===
	loginRouter.Assambler(mux, db)
	sessionRouter.Assambler(mux, db)
	homeRouter.Assambler(mux)
	usersRouter.Assambler(mux, db)
	employeeRouter.Assambler(mux, db)

	// Levantar Servidor
	fmt.Printf("🚀 Server listening on: http://%v:%v", utils.Host, utils.Port)
	if err = http.ListenAndServe(utils.Host+":"+utils.Port, mux); err != nil {
		log.Fatalf("Server has stopped, error: (%v)", err)
	}

}
