package router

import (
	"database/sql"
	"net/http"
	"saavedra/services/CuotasIMSS/api"
	"saavedra/services/CuotasIMSS/service"
	"saavedra/services/CuotasIMSS/store"
	"saavedra/utils"
)

func Assambler(mux *http.ServeMux, db *sql.DB) {
	storing := store.New(db)
	servicing := service.New(storing)
	handler := api.New(servicing)

	mux.Handle("/cuotas", utils.AuthMiddleware(http.HandlerFunc(handler.CallListCuotasIMSS)))
	mux.Handle("/cuotas/slice", utils.AuthMiddleware(http.HandlerFunc(handler.CallSliceCuotasIMSS)))
	mux.Handle("/cuotas/search", utils.AuthMiddleware(http.HandlerFunc(handler.CallCuotasSearch)))
	mux.Handle("/cuotas/new", utils.AuthMiddleware(http.HandlerFunc(handler.CallCuotasIMSSNew)))
	mux.Handle("/cuotas/many2one", utils.AuthMiddleware(http.HandlerFunc(handler.CallCuotasIMSSMany2One)))
}
