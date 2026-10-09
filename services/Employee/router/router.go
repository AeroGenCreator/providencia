package router

import (
	"database/sql"
	"net/http"
	"saavedra/services/Employee/api"
	"saavedra/services/Employee/service"
	"saavedra/services/Employee/store"
	"saavedra/utils"
)

func Assambler(mux *http.ServeMux, db *sql.DB) {
	storing := store.New(db)
	servicing := service.New(storing)
	handler := api.New(servicing)

	mux.Handle("/employee", utils.AuthMiddleware(http.HandlerFunc(handler.CallEmployees)))
	mux.Handle("/employee/new", utils.AuthMiddleware(http.HandlerFunc(handler.CallEmployeesNew)))
	mux.Handle("/employee/slice", utils.AuthMiddleware(http.HandlerFunc(handler.CallSliceEmployees)))
	mux.Handle("/employee/record", utils.AuthMiddleware(http.HandlerFunc(handler.CallEmployeesRecord)))
	mux.Handle("/employee/search", utils.AuthMiddleware(http.HandlerFunc(handler.CallEmployeesSearch)))

}
