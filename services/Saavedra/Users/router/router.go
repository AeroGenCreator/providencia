package router

import (
	"database/sql"
	"net/http"
	"saavedra/services/Saavedra/Users/api"
	"saavedra/services/Saavedra/Users/service"
	"saavedra/services/Saavedra/Users/store"
	"saavedra/utils"
)

func Assambler(mux *http.ServeMux, db *sql.DB) {
	storing := store.New(db)
	servicing := service.New(storing)
	handler := api.New(servicing)

	mux.Handle("/users", utils.AuthMiddleware(http.HandlerFunc(handler.CallUsers)))
	mux.Handle("/users/new", utils.AuthMiddleware(http.HandlerFunc(handler.CallUsersNew)))
	mux.Handle("/users/slice", utils.AuthMiddleware(http.HandlerFunc(handler.CallSliceUsers)))
	mux.Handle("/users/record", utils.AuthMiddleware(http.HandlerFunc(handler.CallUsersRecord)))
	mux.Handle("/users/search", utils.AuthMiddleware(http.HandlerFunc(handler.CallUsersSearch)))
	mux.Handle("/users/many2one", utils.AuthMiddleware(http.HandlerFunc(handler.CallUsersMany2One)))
}
