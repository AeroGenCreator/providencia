package router

import (
	"database/sql"
	"net/http"
	"saavedra/services/Saavedra/Session/api"
	"saavedra/services/Saavedra/Session/service"
	"saavedra/services/Saavedra/Session/store"
	"saavedra/utils"
)

func Assambler(mux *http.ServeMux, db *sql.DB) {
	storing := store.New(db)
	servicing := service.New(storing)
	handler := api.New(servicing)

	mux.Handle("/session", utils.LowLevelMiddleware(http.HandlerFunc(handler.CallSession)))
}
