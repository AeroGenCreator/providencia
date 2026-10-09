package router

import (
	"net/http"
	"saavedra/services/Saavedra/Home/api"
	"saavedra/utils"
)

func Assambler(mux *http.ServeMux) {
	mux.Handle("/home", utils.AuthMiddleware(http.HandlerFunc(api.CallHome)))
}
