package views

import (
	"net/http"
)

type commentForm struct {
	Comment string
}

func PostComment(view *viewState) {
	var commentData commentForm

	if view.parseForm(&commentData) != nil {
		httpError(view.response, http.StatusBadRequest)
		return
	}

	if view.Stores.NewComment(commentData.Comment) != nil {
		httpError(view.response, http.StatusBadRequest)
		return
	}

	view.redirect(view.routes.Default, nil)
}

func GetComments(view *viewState) {
	comments, err := view.Stores.GetAllComments()
	if err != nil {
		httpError(view.response, http.StatusBadRequest)
		return
	}

	context := map[string]interface{}{
		"comments": comments,
	}

	user := view.checkLogin()
	if user != nil {
		context["username"] = user.Username
	}

	view.Templates.ExecuteTemplate(
		view.response,
		"index.tmpl",
		context,
	)
}
