package views

import (
	"github.com/hcnelson99/social/app/stores"
	"net/http"
)

const (
	LOGIN_TEMPLATE = "login.tmpl"
)

func renderAuthTemplate(view *viewState, template string, context map[string]interface{}) {
	context["uri"] = view.request.URL.RequestURI()
	view.Templates.ExecuteTemplate(view.response, template, context)
}

func GetLogin(view *viewState) {
	renderAuthTemplate(view, LOGIN_TEMPLATE, map[string]interface{}{})
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
		next, ok := view.request.URL.Query()[CONTINUE_QUERY_KEY]
		if ok && len(next) == 1 {
			view.redirect(next[0])
		} else {
			view.redirect("/")
		}
	} else {
		if authStatus != stores.AUTH_REJECTED {
			context["error"] = "Internal server error while logging in. Please try again later."
		}
		context["invalid_auth"] = true
	}

	renderAuthTemplate(view, LOGIN_TEMPLATE, context)
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

	view.redirect("/")
}
