package views

import (
	"log"
	"net/http"
)

func getPostFormValue(r *http.Request, key string) (string, bool) {
	r.ParseForm()
	values, success := r.PostForm[key]
	if !success || len(values) != 1 {
		return "", false
	}
	return values[0], true
}

func (app *appViews) PostComment(w http.ResponseWriter, r *http.Request) {
	comment, success := getPostFormValue(r, "comment")
	if !success {
		httpError(w, http.StatusBadRequest)
		return
	}

	if app.Stores.NewComment(comment) != nil {
		httpError(w, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *appViews) GetComments(w http.ResponseWriter, r *http.Request) {
	comments, err := app.Stores.GetAllComments()
	if err != nil {
		httpError(w, http.StatusBadRequest)
		return
	}

	_, err = app.SessionStore.Get(r, USER_SESSION_NAME)
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

	app.Templates.ExecuteTemplate(w, "index.tmpl",
		map[string]interface{}{
			"comments": comments,
			"username": "Steven Shan",
		},
	)
}
