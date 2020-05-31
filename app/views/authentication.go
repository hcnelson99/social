package views

import (
	"github.com/hcnelson99/social/app/stores"
	"net/http"
)

const (
	LOGIN_TEMPLATE               = "login.tmpl"
	INVALIDATE_SESSIONS_TEMPLATE = "invalidate_sessions.tmpl"
)

func renderAuthTemplate(view *viewState, template string, context map[string]interface{}) {
	if context == nil {
		context = map[string]interface{}{}
	}
	context["uri"] = view.request.URL.RequestURI()
	view.Templates.ExecuteTemplate(view.response, template, context)
}

func GetLogin(view *viewState) {
	// if user is already logged in, redirect to the homepage
	if view.checkLogin() != nil {
		view.redirect(view.routes.Default, nil)
		return
	}

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
			view.redirect(next[0], nil)
		} else {
			view.redirect(view.routes.Default, nil)
		}
	} else {
		if authStatus != stores.AUTH_REJECTED {
			context["error"] = "Internal server error while logging in. Please try again later."
		}
		context["invalid_auth"] = true
	}

	renderAuthTemplate(view, LOGIN_TEMPLATE, context)
}

func PostRegister(view *viewState) {
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

	view.redirect(view.routes.Default, nil)
}

/*
   Logs out the user and redirects them to the homepage.
*/
func GetLogout(view *viewState) {
	view.clearSession()
	view.redirect(view.routes.Default, nil)
}

/*
   Logs the user out of all devices and redirects them to the homepage.

   This increments the session generation number in the table.
*/
func GetInvalidateUserSessions(view *viewState) {
	var queryParams map[string]string

	if user := view.checkLogin(); user != nil {
		if !view.Stores.InvalidateUserSessions(user.UserId) {
			queryParams = map[string]string{
				ERROR_QUERY_KEY: "error-invalidate-user-sessions",
			}
		}
	}

	view.redirect(view.routes.Default, queryParams)
}
