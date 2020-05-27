package views

import (
	"github.com/hcnelson99/social/app/types"
	"net/http"
)

const (
	USER_SESSION_NAME = "login"
)

type appViews struct {
	*types.App
}

type Views interface {
	GetComments(http.ResponseWriter, *http.Request)
	PostComment(http.ResponseWriter, *http.Request)
	GetLogin(http.ResponseWriter, *http.Request)
	PostLogin(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
}

func For(app *types.App) Views {
	return &appViews{app}
}

func httpError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
