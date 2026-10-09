package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

const (
	UserId    ContextKey = "id"
	UserName  ContextKey = "name"
	UserParty ContextKey = "party"
	UserToken ContextKey = "token"
)

func InsertAdmin(db *sql.DB) error {
	hashedPassword, err := HashPassword(AdminPassword)
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Error when hashing admin password.")
	}
	// Se injectan las credenciales de admnistrador.
	q1 := `
	INSERT INTO users(name, email, password, party)
	VALUES(?, ?, ?, ?)
	ON CONFLICT DO UPDATE SET
	name = excluded.name,
	email = excluded.email,
	password = excluded.password,
	party = excluded.party;`

	_, err = db.Exec(q1, AdminName, AdminEmail, hashedPassword, "admin")
	if err != nil {
		log.Panicln(err.Error())
		log.Fatal("Admin injection error.")
	}
	// Almacenar de manera global el 'id' del administrador.
	q2 := "SELECT id FROM users WHERE name = ? AND email = ?;"
	var id int
	if err = db.QueryRow(q2, AdminName, AdminEmail).Scan(&id); err != nil {
		log.Fatalf("Error; query admin users (%v)", err.Error())
	}
	AdminId = id
	return nil
}

func LowLevelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		// 1. Manejo de errores al obtener cookies.
		if err == http.ErrNoCookie {
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			log.Println("Error LowLevelMiddleware cookies handler.", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tokenStr := c.Value
		// 2. Si la cookie esta vacia se continua a la siguiente función.
		if tokenStr == "" {
			next.ServeHTTP(w, r)
			return
		}
		// 3. Si existe valor en la cookie, se parsea el token y se extraen valores.
		claims := &Claims{}
		parser := jwt.NewParser()
		_, _, err = parser.ParseUnverified(tokenStr, claims)
		if err != nil {
			log.Println("Error parsing token from request into claims structure.", err.Error())
			next.ServeHTTP(w, r)
			return
		}
		// 4. Inyección exitosa del valor en el contexto
		ctx := context.WithValue(r.Context(), UserId, claims.UserId)
		ctx = context.WithValue(ctx, UserName, claims.UserName)
		ctx = context.WithValue(ctx, UserParty, claims.UserParty)
		ctx = context.WithValue(ctx, UserToken, tokenStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoadSaavedraConfigFile() {
	// Lectura y parseo del JSON saavedra de configuración.
	fileBytes, err := os.ReadFile("config/saavedra.json")
	if err != nil {
		log.Fatalf("Error loading saavedra config file: (%v)", err.Error())
	}
	var saavedra SaavedraConfigFile
	if err = json.Unmarshal(fileBytes, &saavedra); err != nil {
		log.Fatalf("Error parsing json file: (%v)", err.Error())
	}
	// Se cargan datos en buffer para facil acceso. (Permisos) (Niveles de Permisos)
	Parties = saavedra.Parties
	Levels = saavedra.Levels
}
