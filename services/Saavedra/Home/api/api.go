package api

import (
	"html/template"
	"log"
	"net/http"
	"saavedra/utils"
)

// ROUTE "/home"
func CallHome(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userName, ok := r.Context().Value(utils.UserName).(string)
		if !ok {
			log.Println("Error extracting context keys")
			http.Error(w, "Error extracting context keys", http.StatusInternalServerError)
			return
		}
		tpl, err := template.ParseFiles("services/Saavedra/Home/views/home.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err = tpl.Execute(w, map[string]string{"UserName": userName}); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodHead:
		w.WriteHeader(http.StatusOK)
		return
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
}
