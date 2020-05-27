package app

import (
	"github.com/gorilla/mux"
	"github.com/hcnelson99/social/app/types"
	"github.com/hcnelson99/social/app/views"
	"net/http"
)

func GetRouter(app *types.App) *mux.Router {
	staticFileServer := http.FileServer(http.Dir("./app/static"))

	v := views.For(app)

	// TODO: add CSRF everywhere!
	r := mux.NewRouter()
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", staticFileServer))
	r.HandleFunc("/", v.GetComments).Methods("GET")
	r.HandleFunc("/comment", v.PostComment).Methods("POST")
	r.HandleFunc("/login", v.GetLogin).Methods("GET")
	r.HandleFunc("/login/post", v.PostLogin).Methods("POST")
	r.HandleFunc("/login/register", v.Register).Methods("POST")

	return r
}
