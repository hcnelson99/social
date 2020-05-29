package views

import (
	"log"
	"net/http"
)

func GetLogin(view *viewState) {
	view.Templates.ExecuteTemplate(view.response, "login.tmpl", nil)
}

func PostLogin(view *viewState) {
	username, success := getPostFormValue(view.request, "username")
	if !success {
		httpError(view.response, http.StatusBadRequest)
		return
	}
	password, success := getPostFormValue(view.request, "password")
	if !success {
		httpError(view.response, http.StatusBadRequest)
		return
	}

	log.Println(username, password)
}

func Register(view *viewState) {
	username, success := getPostFormValue(view.request, "username")
	if !success {
		httpError(view.response, http.StatusBadRequest)
		return
	}
	password, success := getPostFormValue(view.request, "password")
	if !success {
		httpError(view.response, http.StatusBadRequest)
		return
	}

	// TODO: sanity check password and username
	_, sessionGeneration, err := view.Stores.NewUser(username, password)
	if err != nil {
		httpError(view.response, http.StatusInternalServerError)
		return
	}

	session, err := view.SessionStore.Get(view.request, USER_SESSION_NAME)
	if err != nil {
		httpError(view.response, http.StatusInternalServerError)
		return
	}

	session.Values["session_generation"] = sessionGeneration
	err = session.Save(view.request, view.response)
	if err != nil {
		httpError(view.response, http.StatusInternalServerError)
		return
	}

	http.Redirect(view.response, view.request, "/", http.StatusFound)
}
