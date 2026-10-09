package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"saavedra/services/Saavedra/Users/service"
	"saavedra/services/Saavedra/Users/types"
	"saavedra/utils"
	"strconv"
)

type EndpointHandler struct {
	service service.Service
}

func New(service service.Service) EndpointHandler {
	return EndpointHandler{service: service}
}

// ROUTE: "/users/slice": viaja numero de pagina(page) en la URL.
func (e EndpointHandler) CallSliceUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		page := r.URL.Query().Get("page")
		records, err := e.service.SliceUser(page)
		if err != nil {
			log.Printf("Error '/users/slice' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(records)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE: '/users'
func (e EndpointHandler) CallUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(0, party)
		if err != nil {
			log.Printf("Error '/users' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tpl, err := template.ParseFiles("services/Saavedra/Users/views/list.html")
		if err != nil {
			log.Printf("Error '/users' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tpl.Execute(w, map[string]string{"Service": "Usuarios"})
		if err != nil {
			log.Printf("Error '/users' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE: "/users/new"
func (e EndpointHandler) CallUsersNew(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(0, party)
		if err != nil {
			log.Printf("Error '/users' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tpl, err := template.ParseFiles("services/Saavedra/Users/views/new.html")
		if err != nil {
			log.Printf("Error '/users/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		breadcrumb := map[string]string{"Service": "Usuarios", "View1": "Nuevo"}
		if err = tpl.Execute(w, breadcrumb); err != nil {
			log.Printf("Error '/users/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		var record types.User
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			log.Printf("Error '/users/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := e.service.CreateUser(&record); err != nil {
			log.Printf("Error '/users/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE: "users/record"
func (e EndpointHandler) CallUsersRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(0, party)
		if err != nil {
			log.Printf("Error '/users' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id := r.URL.Query().Get("id")
		intId, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("Error '/users/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tpl, err := template.ParseFiles("services/Saavedra/Users/views/record.html")
		if err != nil {
			log.Printf("Error '/users/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		breadcrumb := map[string]any{"Id": intId, "Service": "Usuarios", "View1": "Editar Usuario"}
		if err = tpl.Execute(w, breadcrumb); err != nil {
			log.Printf("Error '/users/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		var body types.Body
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Printf("Error '/users/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		record, err := e.service.ReadUser(body.Id)
		if err != nil {
			log.Printf("Error '/users/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(record)
	case http.MethodPut:
		var user types.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Printf("Error '/users/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := e.service.UpdateUser(&user); err != nil {
			log.Printf("Error '/users/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		var body types.Body
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Printf("Error '/users/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rId := r.Context().Value(utils.UserId).(string)
		err := e.service.DeleteUser(body.Id, rId)
		switch err {
		case types.ErrAdminProtection:
			w.WriteHeader(http.StatusUnauthorized)
			return
		case types.ErrSelfProtection:
			w.WriteHeader(http.StatusUnauthorized)
			return
		case nil:
			w.WriteHeader(http.StatusOK)
		default:
			log.Printf("Error '/users/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE: "/users/many2one"
func (e EndpointHandler) CallUsersMany2One(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		many2one := e.service.Many2One()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(many2one)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE: "/users/search"
func (e EndpointHandler) CallUsersSearch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pattern := r.URL.Query().Get("pattern")
		page := r.URL.Query().Get("page")
		slice, err := e.service.SearchUser(pattern, page)
		if err != nil {
			log.Printf("Error '/users/search' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(slice)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}
