package views

import (
	"log"
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

	_, err = view.SessionStore.Get(view.request, USER_SESSION_NAME)
	if err != nil {
		log.Fatal(err)
	}
	/*
		username := ""
		if session_key, found := session.Values["session_key"]; found == true {
			row := app.Stores.Db.QueryRow(context.Background(), "SELECT user_id FROM user_sessions WHERE session_key = $1", session_key)
			var user_id int
			err = row.Scan(&user_id)
			if err != nil {
				log.Fatal(err)
			}

			row = app.Stores.Db.QueryRow(context.Background(), "SELECT username FROM users WHERE id = $1", user_id)
			err = row.Scan(&username)
			if err != nil {
				log.Fatal(err)
			}
		}
	*/

	view.Templates.ExecuteTemplate(
		view.response,
		"index.tmpl",
		map[string]interface{}{
			"comments": comments,
			"username": "Steven Shan",
		},
	)
}
