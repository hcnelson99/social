// Interface for setting and getting user session values
package views

import (
	"github.com/gorilla/sessions"
	"github.com/hcnelson99/social/app/stores"
	"log"
	"net/http"
	"net/url"
	"path"
)

const (
	USER_SESSION_NAME      = "user"
	SESSION_USER_ID_KEY    = "user_id"
	SESSION_GENERATION_KEY = "session_generation"

	// URL query key indicating a URI to redirect to
	CONTINUE_QUERY_KEY = "continue"
)

func getUserSession(view *viewState) *sessions.Session {
	session, err := view.SessionStore.Get(view.request, USER_SESSION_NAME)
	if session == nil {
		log.Print("getting session returned nil", err)
	} else if err != nil {
		log.Print("error decoding session", err)
	}
	return session
}

/*
   Sets the user session so the user stays authenticated.

   Returns true if the session was created successfully, and false on failure.
*/
func (view *viewState) setUserSession(userId, sessionGeneration int) bool {
	session := getUserSession(view)
	if session == nil {
		return false
	}
	session.Values[SESSION_USER_ID_KEY] = userId
	session.Values[SESSION_GENERATION_KEY] = sessionGeneration
	if err := session.Save(view.request, view.response); err != nil {
		log.Print("couldn't save user session", err)
		return false
	}
	return true
}

/*
   Clears the user's current session.
*/
func (view *viewState) clearSession() {
	if session := getUserSession(view); session != nil {
		session.Options.MaxAge = -1
	}
}

func (view *viewState) redirect(uri string) {
	url := path.Join("/", uri)
	http.Redirect(view.response, view.request, url, http.StatusFound)
}

/*
   Returns a pointer to a user object if the user is authenticated or else nil.
*/
func (view *viewState) checkLogin() *stores.User {
	session := getUserSession(view)
	if session == nil {
		return nil
	}
	userId, gotUserId := session.Values[SESSION_USER_ID_KEY]
	sessionGen, gotSessionGen := session.Values[SESSION_GENERATION_KEY]

	if !gotUserId || !gotSessionGen {
		return nil
	}

	userIdVal, userIdValid := userId.(int)
	sessionGenVal, sessionGenValid := sessionGen.(int)
	if userIdValid && sessionGenValid {
		return view.Stores.CheckUserSession(userIdVal, sessionGenVal)
	}

	// clear session if authentication failed so we don't repeatedly query
	// the database to check if the user is logged in
	view.clearSession()

	return nil
}

/*
   Redirects user to the login page.

   This also handles directing the user back to the page they were on after they log in.
*/
func (view *viewState) gotoLogin() {
	base, err := url.Parse(view.routes.Login)
	if err != nil || view.request.Method != "GET" {
		view.redirect(view.routes.Login)
		return
	}

	params := url.Values{}
	uri := view.request.URL.RequestURI()
	if uri != "/" {
		params.Add(CONTINUE_QUERY_KEY, uri)
	}
	base.RawQuery = params.Encode()

	view.redirect(base.String())
}

