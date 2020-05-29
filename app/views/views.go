package views

import (
	"github.com/hcnelson99/social/app/types"
	"net/http"
)

type viewState struct {
	*types.App
	response http.ResponseWriter
	request  *http.Request
}

type ViewFunction = func(*viewState)
type HttpHandler = func(http.ResponseWriter, *http.Request)

func Get(app *types.App, viewFunc ViewFunction) HttpHandler {
	return func(response http.ResponseWriter, request *http.Request) {
		viewFunc(&viewState{app, response, request})
	}
}

func httpError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
