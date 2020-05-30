// Interface for setting and getting user session values
package views

import (
	"github.com/gorilla/sessions"
	"github.com/hcnelson99/social/app/stores"
	"log"
)

const (
	USER_SESSION_NAME      = "user"
	SESSION_USER_ID_KEY    = "user_id"
	SESSION_GENERATION_KEY = "session_generation"
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

   Returns false if the session was created successfully, and true on failure.
*/
func (view *viewState) setUserSession(userId, sessionGeneration int) bool {
	session := getUserSession(view)
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
	getUserSession(view).Options.MaxAge = -1
}

/*
   Returns a pointer to a user object if the user is authenticated or else nil.
*/
func (view *viewState) checkLogin() *stores.User {
	session := getUserSession(view)
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
