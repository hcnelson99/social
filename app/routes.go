package app

import (
	"github.com/gorilla/mux"
	"github.com/hcnelson99/social/app/types"
	"github.com/hcnelson99/social/app/views"
	"net/http"
)

const (
	DEFAULT_ROUTE = "/"
	LOGIN_ROUTE   = "/login"
	ERROR_ROUTE   = "/error"
)

func GetRouter(app *types.App) *mux.Router {
	staticFileServer := http.FileServer(http.Dir("./app/static"))

	routes := views.RouteConfig{
		Default: DEFAULT_ROUTE,
		Login:   LOGIN_ROUTE,
		Error:   ERROR_ROUTE,
	}

	getView := func(v views.ViewFunction) views.HttpHandler {
		return views.Get(app, routes, v)
	}

	// TODO: add CSRF everywhere!
	r := mux.NewRouter()
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", staticFileServer))
	r.HandleFunc(DEFAULT_ROUTE, getView(views.GetComments)).Methods("GET")
	r.HandleFunc("/comment", getView(views.PostComment)).Methods("POST")
	r.HandleFunc(
		LOGIN_ROUTE,
		views.GetMethods(app, routes, views.HandlerMap{
			GET:  views.GetLogin,
			POST: views.PostLogin,
		}),
	).Methods("GET", "POST")
	r.HandleFunc(
		"/register",
		views.GetMethods(app, routes, views.HandlerMap{
			GET:  views.GetRegister,
			POST: views.PostRegister,
		}),
	).Methods("GET", "POST")
	r.HandleFunc("/error", getView(views.GetError)).Methods("GET")
	r.HandleFunc("/logout", getView(views.GetLogout)).Methods("GET")
	r.HandleFunc("/logout-all", getView(views.GetInvalidateUserSessions)).Methods("GET")

	return r
}
