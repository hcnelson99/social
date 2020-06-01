package views

import (
	"github.com/hcnelson99/social/app/stores"
)

const (
	LOGIN_TEMPLATE               = "login.tmpl"
	REGISTER_TEMPLATE            = "register.tmpl"
	INVALIDATE_SESSIONS_TEMPLATE = "invalidate_sessions.tmpl"
)

type authForm struct {
	Username string
	Password string
}

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

	renderAuthTemplate(view, LOGIN_TEMPLATE, nil)
}

func PostLogin(view *viewState) {
	context := map[string]interface{}{}

	var loginData authForm
	if view.parseForm(&loginData) == nil {
		username := loginData.Username
		password := loginData.Password

		user, sessionGen, authStatus := view.Stores.Login(username, password)
		if authStatus == stores.AUTH_VALIDATED && view.setUserSession(user.UserId, sessionGen) {
			next, ok := view.request.URL.Query()[CONTINUE_QUERY_KEY]
			if ok && len(next) == 1 {
				view.redirect(next[0], nil)
			} else {
				view.redirect(view.routes.Default, nil)
			}
		} else if authStatus != stores.AUTH_REJECTED {
			context[TEMPLATE_ERROR] = view.language.ERROR_LOGGING_IN
		} else {
			context[TEMPLATE_ERROR] = view.language.LOGIN_REJECTED
		}
	} else {
		context[TEMPLATE_ERROR] = view.language.ERROR_LOGGING_IN
	}

	renderAuthTemplate(view, LOGIN_TEMPLATE, context)
}

func GetRegister(view *viewState) {
	if view.checkLogin() != nil {
		view.redirect(view.routes.Default, nil)
		return
	}

	renderAuthTemplate(view, REGISTER_TEMPLATE, nil)
}

func PostRegister(view *viewState) {
	var registerData authForm
	context := map[string]interface{}{}
	if view.parseForm(&registerData) == nil {
		username := registerData.Username
		password := registerData.Password

		// TODO: sanity check password and username
		user, sessionGeneration := view.Stores.NewUser(username, password)
		if user == nil {
			context[TEMPLATE_ERROR] = view.language.USERNAME_TAKEN
		} else if view.setUserSession(user.UserId, sessionGeneration) {
			view.redirect(view.routes.Default, nil)
			return
		} else {
			context[TEMPLATE_ERROR] = view.language.ERROR_REGISTERING
		}
	} else {
		context[TEMPLATE_ERROR] = view.language.ERROR_REGISTERING
	}

	renderAuthTemplate(view, REGISTER_TEMPLATE, context)
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
	if user := view.checkLogin(); user != nil {
		if !view.Stores.InvalidateUserSessions(user.UserId) {
			queryParams := map[string]string{
				ERROR_QUERY_KEY: ERROR_INVALIDATE_USER_SESSIONS,
			}
			view.redirect(view.routes.Error, queryParams)
			return
		}
	}

	view.redirect(view.routes.Default, nil)
}
