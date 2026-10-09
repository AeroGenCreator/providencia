package api

import (
	"log"
	"net/http"
	"saavedra/services/Saavedra/Session/service"
	"saavedra/utils"
	"time"
)

type EndpointHandler struct {
	service service.Service
}

func New(service service.Service) *EndpointHandler {
	return &EndpointHandler{service: service}
}

// ROUTE: "/session"
func (e EndpointHandler) CallSession(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		id, ok1 := r.Context().Value(utils.UserId).(string)
		name, ok2 := r.Context().Value(utils.UserName).(string)
		party, ok3 := r.Context().Value(utils.UserParty).(string)
		token, ok4 := r.Context().Value(utils.UserToken).(string)
		if !ok1 || !ok2 || !ok3 || !ok4 {
			log.Println("Error extracting context values")
			http.Error(w, "Error extracting context values", http.StatusInternalServerError)
			return
		}
		log.Println("Attempting refresh user session token...")
		newToken, err := e.service.RefreshToken(id, name, party, token)
		if err != nil {
			log.Printf("Error generating new user session token (%v)", err.Error())
			http.Error(w, "Error generating new user session token", http.StatusInternalServerError)
			return
		}
		if newToken == "" {
			if err = e.service.ExpellUser(id); err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    newToken,
			Expires:  time.Now().Add(12 * time.Hour),
			Path:     "/",
			HttpOnly: true,
			Secure:   utils.IsProduction,
			SameSite: http.SameSiteLaxMode,
		})
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodPut:
		id, ok := r.Context().Value(utils.UserId).(string)
		if !ok {
			log.Println("Error extracting context keys.")
			http.Error(w, "Error extracting context keys.", http.StatusInternalServerError)
			return
		}
		if err := e.service.ExpellUser(id); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Expires:  time.Now().Add(12 * time.Hour),
			Path:     "/",
			HttpOnly: true,
			Secure:   utils.IsProduction,
			SameSite: http.SameSiteLaxMode,
		})
		w.WriteHeader(http.StatusOK)
		return
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
}
