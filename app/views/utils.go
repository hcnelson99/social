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

		if err := session.Save(view.request, view.response); err != nil {
			log.Print("failed to clear session", err)
		}
	}
}

func (view *viewState) redirect(redirectPath string, queryParams map[string]string) {
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	url := path.Join("/", redirectPath)
	if rawParams := params.Encode(); rawParams != "" {
		url = url + "?" + rawParams
	}
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
	var queryParams map[string]string

	uri := view.request.URL.RequestURI()
	if view.request.Method == "GET" && uri != view.routes.Default {
		queryParams = map[string]string{
			CONTINUE_QUERY_KEY: uri,
		}
	}

	view.redirect(view.routes.Login, queryParams)
}

/*
   Wrapper to parse POST form data.

   Arguments:
       form: a pointer to the form struct to write the data to (must be non-nil)

   Returns nil if the form was parsed successfully, else return an error.
*/
func (view *viewState) parseForm(form interface{}) error {
	err := view.request.ParseForm()
	if err != nil {
		return err
	}

	err = view.SchemaDecoder.Decode(form, view.request.PostForm)
	if err != nil {
		return err
	}

	return nil
}
