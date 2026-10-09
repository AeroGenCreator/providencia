package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"saavedra/services/Saavedra/Login/service"
	"saavedra/services/Saavedra/Users/types"
	"saavedra/utils"
	"time"
)

type EndpointHandler struct {
	service service.Service
}

func New(service service.Service) *EndpointHandler {
	return &EndpointHandler{service: service}
}

// ROUTE: "/"
func (e EndpointHandler) CallRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tpl, err := template.ParseFiles("services/Saavedra/Login/views/login.html")
		if err != nil {
			log.Println("Error Path: '/', Method: Get", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = tpl.Execute(w, nil); err != nil {
			log.Println("Error Path: '/', Method: Get", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		var user types.UserStr
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		credentials, err := e.service.Login(&user)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if credentials.Token == "" {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    credentials.Token,
			Expires:  time.Now().Add(12 * time.Hour),
			Path:     "/",
			HttpOnly: true,
			Secure:   utils.IsProduction,
			SameSite: http.SameSiteLaxMode,
		})
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodHead:
		w.WriteHeader(http.StatusOK)
		return
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
}

// ROUTE: "/obtain"
func (e EndpointHandler) CallObtain(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tpl, err := template.ParseFiles("services/Saavedra/Login/views/obtain.html")
		if err != nil {
			log.Printf("Error parsing (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tpl.Execute(w, map[string]string{"Title": "Obtener Saavedra"}); err != nil {
			log.Printf("Error rendering (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE: "/tutorial"
func (e EndpointHandler) CallTutorial(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tpl, err := template.ParseFiles("services/Saavedra/Login/views/tutorial.html")
		if err != nil {
			log.Printf("Error parsing (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tpl.Execute(w, map[string]string{"Title": "Tutorial"}); err != nil {
			log.Printf("Error rendering (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE: "/docs"
func (e EndpointHandler) CallDocs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tpl, err := template.ParseFiles("services/Saavedra/Login/views/docs.html")
		if err != nil {
			log.Printf("Error parsing (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tpl.Execute(w, map[string]string{"Title": "Documentación"}); err != nil {
			log.Printf("Error rendering (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE: "/philosophy"
func (e EndpointHandler) CallPhilosophy(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tpl, err := template.ParseFiles("services/Saavedra/Login/views/philosophy.html")
		if err != nil {
			log.Printf("Error parsing (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tpl.Execute(w, map[string]string{"Title": "Filosofía"}); err != nil {
			log.Printf("Error rendering (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}
