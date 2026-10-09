package router

import "net/http"

// Endpoint para la carpeta assets
const GlobalAssetsPath = "/assets/"

func ServeGlobalAssetsDirectory(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("./assets/"))
	mux.Handle(GlobalAssetsPath, http.StripPrefix(GlobalAssetsPath, fs))
}
