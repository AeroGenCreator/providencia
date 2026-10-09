package router

import (
	"database/sql"
	"net/http"
	"saavedra/services/Saavedra/Login/api"
	"saavedra/services/Saavedra/Login/service"
	"saavedra/services/Saavedra/Login/store"
	"saavedra/utils"
)

func Assambler(mux *http.ServeMux, db *sql.DB) {
	storing := store.New(db)
	servicing := service.New(storing)
	handler := api.New(servicing)

	mux.Handle("/", utils.LowLevelMiddleware(http.HandlerFunc(handler.CallRoot)))
	mux.HandleFunc("/philosophy", handler.CallPhilosophy)
	mux.HandleFunc("/tutorial", handler.CallTutorial)
	mux.HandleFunc("/obtain", handler.CallObtain)
	mux.HandleFunc("/docs", handler.CallDocs)
}
