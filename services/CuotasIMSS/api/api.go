package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"saavedra/services/CuotasIMSS/service"
	"saavedra/services/CuotasIMSS/types"
	"saavedra/utils"
)

type EndpointHandler struct {
	service service.Service
}

func New(service service.Service) EndpointHandler {
	return EndpointHandler{service: service}
}

// ROUTE: "/cuotas/slice"
func (e EndpointHandler) CallSliceCuotasIMSS(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		page := r.URL.Query().Get("page")
		data, err := e.service.SliceCuotasIMSS(page)
		if err != nil {
			log.Printf("Error /cuotas/slice (%v)", err.Error())
			http.Error(w, "Error /cuotas/slice ", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
}

// ROUTE: "/cuotas/search"
func (e EndpointHandler) CallCuotasSearch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pattern := r.URL.Query().Get("pattern")
		page := r.URL.Query().Get("page")
		data, err := e.service.SearchCuotasIMSS(pattern, page)
		if err != nil {
			log.Printf("Error /cuotas/search (%v)", err.Error())
			http.Error(w, "Error /cuotas/search ", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
}

// ROUTE: "/cuotas"
func (e EndpointHandler) CallListCuotasIMSS(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(1, party)
		if err != nil {
			log.Printf("Error /cuotas - permissions - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tpl, err := template.ParseFiles("services/CuotasIMSS/views/list.html")
		if err != nil {
			log.Printf("Error /cuotas - FileParsing - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tpl.Execute(w, map[string]string{"Service": "Cuotas Obrero Patronal"}); err != nil {
			log.Printf("Error /cuotas - RenderingFile - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// "/cuotas/new"
func (e EndpointHandler) CallCuotasIMSSNew(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(2, party)
		if err != nil {
			log.Printf("Message (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tpl, err := template.ParseFiles("services/CuotasIMSS/views/new.html")
		if err != nil {
			log.Printf("Error /cuotas/new - ParsingFile - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = tpl.Execute(w, map[string]string{
			"Service": "Cuotas Obrero Patronal", "Nuevo": "Nuevo Registro",
		}); err != nil {
			log.Printf("Error /cuotas/new - RenderingFile - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(2, party)
		if err != nil {
			log.Printf("Message (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var peticion types.Peticion
		if err := json.NewDecoder(r.Body).Decode(&peticion); err != nil {
			log.Printf("Error /cuotas/new - DecodingBody - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respuestaCuotas, err := e.service.CalcularCuotas(&peticion)
		if err != nil {
			log.Printf("Error /cuotas/new - Calculating - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(respuestaCuotas)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// "/cuotas/many2one"
func (e EndpointHandler) CallCuotasIMSSMany2One(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := e.service.Many2OneCuotasIMSS()
		if err != nil {
			log.Printf("Error /cuotas/many2one - ServingMany2One - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}
