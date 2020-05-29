// Interface for setting and getting user session values
package views

import (
	"github.com/gorilla/sessions"
)

const (
	USER_SESSION_NAME      = "user"
	SESSION_USER_ID_KEY    = "user_id"
	SESSION_GENERATION_KEY = "session_generation"
)

func getUserSession(view *viewState) *sessions.Session {
	session, _ := view.SessionStore.Get(view.request, USER_SESSION_NAME)
	// TODO log err
	return session
}

func (view *viewState) setUserSession(userId, sessionGeneration int) error {
	session := getUserSession(view)
	session.Values[SESSION_USER_ID_KEY] = userId
	session.Values[SESSION_GENERATION_KEY] = sessionGeneration
	return session.Save(view.request, view.response)
}
