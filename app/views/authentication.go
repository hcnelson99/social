package views

import (
	"github.com/hcnelson99/social/app/stores"
	"net/http"
)

const (
	LOGIN_TEMPLATE = "login.tmpl"
)

func GetLogin(view *viewState) {
	view.Templates.ExecuteTemplate(view.response, LOGIN_TEMPLATE, nil)
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

	context := map[string]interface{}{
		"invalid_auth": false,
	}

	user, sessionGen, authStatus := view.Stores.Login(username, password)
	if authStatus == stores.AUTH_VALIDATED && view.setUserSession(user.UserId, sessionGen) {
		http.Redirect(view.response, view.request, "/", http.StatusFound)
	} else {
		if authStatus == stores.AUTH_ERROR {
			context["error"] = "Internal server error while logging in. Please try again later."
		}
		context["invalid_auth"] = true
	}

	view.Templates.ExecuteTemplate(view.response, LOGIN_TEMPLATE, context)
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
	user, sessionGeneration := view.Stores.NewUser(username, password)
	if user == nil {
		httpError(view.response, http.StatusInternalServerError)
		return
	}

	if view.setUserSession(user.UserId, sessionGeneration) {
		httpError(view.response, http.StatusInternalServerError)
		return
	}

	http.Redirect(view.response, view.request, "/", http.StatusFound)
}
