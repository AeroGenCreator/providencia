package router

import (
	"database/sql"
	"fmt"
	"net/http"
	"saavedra/services/CuotasIMSS/api"
	"saavedra/services/CuotasIMSS/service"
	"saavedra/services/CuotasIMSS/store"
	"saavedra/utils"
)

func Assambler(mux *http.ServeMux, db *sql.DB) {

	if err := service.CargarCuotasBloques(); err != nil {
		fmt.Printf("Advertencia de carga...")
		fmt.Println(err.Error())
	}
	if err := service.CargarEditables(); err != nil {
		fmt.Printf("Advertencia de carga...")
		fmt.Println(err.Error())
	}
	if err := service.CargarFactores(); err != nil {
		fmt.Printf("Advertencia de carga...")
		fmt.Println(err.Error())
	}
	if err := service.CargarDiasVacaciones(); err != nil {
		fmt.Printf("Advertencia de carga...")
		fmt.Println(err.Error())
	}
	if err := service.CargarCesantiaVejez(); err != nil {
		fmt.Printf("Advertencia de carga...")
		fmt.Println(err.Error())
	}
	if err := service.CargarPrimaRiesgo(); err != nil {
		fmt.Printf("Advertencia de carga...")
		fmt.Println(err.Error())
	}

	storing := store.New(db)
	servicing := service.New(storing)
	handler := api.New(servicing)

	mux.Handle("/cuotas", utils.AuthMiddleware(http.HandlerFunc(handler.CallListCuotasIMSS)))
	mux.Handle("/cuotas/slice", utils.AuthMiddleware(http.HandlerFunc(handler.CallSliceCuotasIMSS)))
	mux.Handle("/cuotas/search", utils.AuthMiddleware(http.HandlerFunc(handler.CallCuotasSearch)))
	mux.Handle("/cuotas/new", utils.AuthMiddleware(http.HandlerFunc(handler.CallCuotasIMSSNew)))
	mux.Handle("/cuotas/config", utils.AuthMiddleware(http.HandlerFunc(handler.CallCuotasIMSSConfig)))
	mux.Handle("/cuotas/record", utils.AuthMiddleware(http.HandlerFunc(handler.CallCuotasIMSSRecord)))
	mux.Handle("/cuotas/many2one", utils.AuthMiddleware(http.HandlerFunc(handler.CallCuotasIMSSMany2One)))
	mux.Handle("/cuotas/employees", utils.AuthMiddleware(http.HandlerFunc(handler.CallEmployeesNames)))
}
