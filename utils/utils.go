package utils

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	// GenerateFromPassword: Espera un arreglo de bytes.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// Se retorna la encriptación como string.
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {

	// Valida una constraseña contra un hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(id int, name, party string) (string, error) {

	// Parsea token de sesión en bytes.
	var jwtKey = []byte(SessionToken)
	expirationTime := time.Now().Add(5 * time.Minute)
	// Tipado: Estructura de un Claims para generación firmada del JWT
	claims := &Claims{
		UserId:    strconv.Itoa(id),
		UserName:  name,
		UserParty: party,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	// Generacion del JWT - Metodo de firmado HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Parseo del JWT a string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Token + Header + Cookies
		jwtKey := []byte(SessionToken)
		requestedWith := r.Header.Get("X-Requested-With")
		c, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Extracción del token desde la petición & validación.
		strToken := c.Value
		claims := &Claims{}
		_, err = jwt.ParseWithClaims(strToken, claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})
		if err != nil {
			if requestedWith != "jsFrontendComponent" {
				log.Println("Token is expired plus no JS component used in request.")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Información de usuario agregada a la petición.
		ctx := context.WithValue(r.Context(), UserId, claims.UserId)
		ctx = context.WithValue(ctx, UserName, claims.UserName)
		ctx = context.WithValue(ctx, UserParty, claims.UserParty)
		ctx = context.WithValue(ctx, UserToken, strToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Total paginas de la tabla: a = Cuenta de registros de la tabla, b = Numero de registros por slice.
func CalculateTotalPages(a, b int) int {
	return (a + b - 1) / b
}
