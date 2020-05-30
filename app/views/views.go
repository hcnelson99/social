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

type HandlerMap struct {
	GET  ViewFunction
	POST ViewFunction
}

func defaultView(view *viewState) {
	httpError(view.response, http.StatusMethodNotAllowed)
}

func callHandler(viewFunc ViewFunction, view *viewState) {
	if viewFunc != nil {
		viewFunc(view)
	} else {
		defaultView(view)
	}
}

func GetMethods(app *types.App, handlers HandlerMap) HttpHandler {
	return Get(app, func(view *viewState) {
		var viewFunc ViewFunction
		switch view.request.Method {
		case "GET":
			viewFunc = handlers.GET
		case "POST":
			viewFunc = handlers.POST
		}
		callHandler(viewFunc, view)
	})
}

func httpError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
