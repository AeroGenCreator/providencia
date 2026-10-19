package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"saavedra/services/CuotasIMSS/service"
	"saavedra/services/CuotasIMSS/types"
	"saavedra/utils"
	"strconv"
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
		name := r.URL.Query().Get("name")
		year := r.URL.Query().Get("year")
		data, err := e.service.SearchCuotasIMSS(name, year)
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

// ROUTE: "/cuotas/many2one"
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

// ROUTE: "/cuotas/employees"
func (e EndpointHandler) CallEmployeesNames(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := e.service.EmployeesNames()
		if err != nil {
			log.Printf("Error /cuotas/employees - ServingEmployeesNames - (%v)", err.Error())
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
		if err := tpl.Execute(w, map[string]string{"Service": "Cuotas"}); err != nil {
			log.Printf("Error /cuotas - RenderingFile - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE "/cuotas/new"
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
			"Service": "Cuotas", "Nuevo": "Nuevo Registro",
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
		cuotas, err := e.service.CalcularCuotas(&peticion)
		if err == types.SalaryError {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err != nil {
			log.Printf("Error /cuotas/new - Calculating - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cuotas)
	case http.MethodPost:
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
		if err = e.service.CreateCuotasIMSS(&peticion); err == types.SalaryError {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err != nil {
			log.Printf("Error /cuotas/new - CreatingCuotasIMSS - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE "cuotas/record"
func (e EndpointHandler) CallCuotasIMSSRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(1, party)
		if err != nil {
			log.Printf("Error '/cuotas/record' (%v)", err.Error())
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
			log.Printf("Error '/cuotas/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tpl, err := template.ParseFiles("services/CuotasIMSS/views/record.html")
		if err != nil {
			log.Printf("Error '/cuotas/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		breadcrumb := map[string]any{"Id": intId, "Service": "Cuotas", "View1": "Leer Registro"}
		if err = tpl.Execute(w, breadcrumb); err != nil {
			log.Printf("Error '/cuotas/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPut:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(1, party)
		if err != nil {
			log.Printf("Error '/cuotas/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var id types.Record
		if err = json.NewDecoder(r.Body).Decode(&id); err != nil {
			log.Printf("Error '/cuotas/record' - JSONDecoding - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		record, err := e.service.ReadCuotasIMSS(id.Id)
		if err != nil {
			log.Printf("Error '/cuotas/record' - ReadingRecord - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(record)
	case http.MethodDelete:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(4, party)
		if err != nil {
			log.Printf("Error '/cuotas/record' - Permissions - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var id types.Record
		if err = json.NewDecoder(r.Body).Decode(&id); err != nil {
			log.Printf("Error '/cuotas/record' - JSONDecoding - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = e.service.DeleteCuotasIMSS(id.Id); err != nil {
			log.Printf("Error '/cuotas/record' - DeleteRecord - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE: "/cuotas/config"
func (e EndpointHandler) CallCuotasIMSSConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(3, party)
		if err != nil {
			log.Printf("Error '/cuotas/record' - Permissions - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tpl, err := template.ParseFiles("services/cuotasIMSS/views/config.html")
		if err != nil {
			log.Printf("Error '/cuotas/config' - ParseTemplate - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = tpl.Execute(w, map[string]string{"Service": "Cuotas", "View1": "Configuración"}); err != nil {
			log.Printf("Error '/cuotas/config' - RenderTemplate - (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}
