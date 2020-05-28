package views

import (
	"log"
	"net/http"
)

func (app *appViews) GetLogin(w http.ResponseWriter, r *http.Request) {
	app.Templates.ExecuteTemplate(w, "login.tmpl", nil)
}

func (app *appViews) PostLogin(w http.ResponseWriter, r *http.Request) {
	username, success := getPostFormValue(r, "username")
	if !success {
		httpError(w, http.StatusBadRequest)
		return
	}
	password, success := getPostFormValue(r, "password")
	if !success {
		httpError(w, http.StatusBadRequest)
		return
	}

	log.Println(username, password)
}

func (app *appViews) Register(w http.ResponseWriter, r *http.Request) {
	username, success := getPostFormValue(r, "username")
	if !success {
		httpError(w, http.StatusBadRequest)
		return
	}
	password, success := getPostFormValue(r, "password")
	if !success {
		httpError(w, http.StatusBadRequest)
		return
	}

	// TODO: sanity check password and username
	userId, err := app.Stores.NewUser(username, password)
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}

	session, err := app.SessionStore.Get(r, USER_SESSION_NAME)
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}

	sessionKey, err := app.Stores.NewUserSession(userId)
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}

	session.Values["session_key"] = sessionKey
	err = session.Save(r, w)
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
