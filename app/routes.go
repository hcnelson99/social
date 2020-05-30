package app

import (
	"github.com/gorilla/mux"
	"github.com/hcnelson99/social/app/types"
	"github.com/hcnelson99/social/app/views"
	"net/http"
)

func GetRouter(app *types.App) *mux.Router {
	staticFileServer := http.FileServer(http.Dir("./app/static"))

	getView := func(v views.ViewFunction) views.HttpHandler {
		return views.Get(app, v)
	}

	// TODO: add CSRF everywhere!
	r := mux.NewRouter()
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", staticFileServer))
	r.HandleFunc("/", getView(views.GetComments)).Methods("GET")
	r.HandleFunc("/comment", getView(views.PostComment)).Methods("POST")
	r.HandleFunc(
		"/login",
		views.GetMethods(app, views.HandlerMap{
			GET:  views.GetLogin,
			POST: views.PostLogin,
		}),
	).Methods("GET", "POST")
	r.HandleFunc("/login/register", getView(views.Register)).Methods("POST")

	return r
}
