package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"saavedra/services/Employee/service"
	"saavedra/services/Employee/types"
	"saavedra/utils"
	"strconv"
)

type EndpointHandler struct {
	service service.Service
}

func New(service service.Service) EndpointHandler {
	return EndpointHandler{service: service}
}

// ROUTE: "/employee/slice": viaja numero de pagina(page) en la URL.
func (e EndpointHandler) CallSliceEmployees(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		page := r.URL.Query().Get("page")
		records, err := e.service.SliceEmployee(page)
		if err != nil {
			log.Printf("Error '/employee/slice' (%v)", err.Error())
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

// ROUTE: '/employee'
func (e EndpointHandler) CallEmployees(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(1, party)
		if err != nil {
			log.Printf("Error '/employee' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tpl, err := template.ParseFiles("services/Employee/views/list.html")
		if err != nil {
			log.Printf("Error '/employee' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tpl.Execute(w, map[string]string{"Service": "Empleados"})
		if err != nil {
			log.Printf("Error '/employee' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE: "/employee/new"
func (e EndpointHandler) CallEmployeesNew(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(1, party)
		if err != nil {
			log.Printf("Error '/employee/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tpl, err := template.ParseFiles("services/Employee/views/new.html")
		if err != nil {
			log.Printf("Error '/employee/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		breadcrumb := map[string]string{"Service": "Empleados", "View1": "Nuevo"}
		if err = tpl.Execute(w, breadcrumb); err != nil {
			log.Printf("Error '/employee/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(2, party)
		if err != nil {
			log.Printf("Error '/employee/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var record types.Employee
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			log.Printf("Error '/employee/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = e.service.CreateEmployee(&record)
		if err == utils.DuplicatedDataError {
			log.Print(utils.DuplicatedDataError.Error())
			http.Error(w, utils.DuplicatedDataError.Error(), http.StatusConflict)
			return
		} else if err != nil {
			log.Printf("Error '/employee/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE: "/employee/record"
func (e EndpointHandler) CallEmployeesRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(1, party)
		if err != nil {
			log.Printf("Error '/employee/record' (%v)", err.Error())
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
			log.Printf("Error '/employee/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tpl, err := template.ParseFiles("services/Employee/views/record.html")
		if err != nil {
			log.Printf("Error '/employee/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		breadcrumb := map[string]any{"Id": intId, "Service": "Empleados", "View1": "Editar Empleado"}
		if err = tpl.Execute(w, breadcrumb); err != nil {
			log.Printf("Error '/employee/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(1, party)
		if err != nil {
			log.Printf("Error '/employee/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var body types.OnlyId
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Printf("Error '/employee/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		record, err := e.service.ReadEmployee(body.Id)
		if err != nil {
			log.Printf("Error '/employee/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(record)
	case http.MethodPut:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(3, party)
		if err != nil {
			log.Printf("Error '/employee/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var employee types.Employee
		if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
			log.Printf("Error '/employee/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := e.service.UpdateEmployee(&employee); err != nil {
			log.Printf("Error '/employee/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		party := r.Context().Value(utils.UserParty).(string)
		valid, err := utils.AddPermissions(4, party)
		if err != nil {
			log.Printf("Error '/employee/new' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var body types.OnlyId
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Printf("Error '/employee/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := e.service.DeleteEmployee(body.Id); err != nil {
			log.Printf("Error '/employee/record' (%v)", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

// ROUTE: "/employee/search"
func (e EndpointHandler) CallEmployeesSearch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pattern := r.URL.Query().Get("pattern")
		page := r.URL.Query().Get("page")
		slice, err := e.service.SearchEmployee(pattern, page)
		if err != nil {
			log.Printf("Error '/employee/search' (%v)", err.Error())
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
