package views

import (
	"net/http"
)

func getPostFormValue(request *http.Request, key string) (string, bool) {
	request.ParseForm()
	values, success := request.PostForm[key]
	if !success || len(values) != 1 {
		return "", false
	}
	return values[0], true
}

func PostComment(view *viewState) {
	comment, success := getPostFormValue(view.request, "comment")
	if !success {
		httpError(view.response, http.StatusBadRequest)
		return
	}

	if view.Stores.NewComment(comment) != nil {
		httpError(view.response, http.StatusBadRequest)
		return
	}

	http.Redirect(view.response, view.request, "/", http.StatusFound)
}

func GetComments(view *viewState) {
	comments, err := view.Stores.GetAllComments()
	if err != nil {
		httpError(view.response, http.StatusBadRequest)
		return
	}

	view.Templates.ExecuteTemplate(
		view.response,
		"index.tmpl",
		map[string]interface{}{
			"comments": comments,
			"username": "Steven Shan",
		},
	)
}
