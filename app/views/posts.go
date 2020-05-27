package views

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

func getPostFormValue(r *http.Request, key string) (string, bool) {
	r.ParseForm()
	values, success := r.PostForm[key]
	if !success || len(values) != 1 {
		return "", false
	}
	return values[0], true
}

func (appl *appViews) PostComment(w http.ResponseWriter, r *http.Request) {
	comment, success := getPostFormValue(r, "comment")
	if !success {
		httpError(w, http.StatusBadRequest)
		return
	}

	_, err := appl.Db.Exec(context.Background(), "INSERT INTO comments(author, text) VALUES ($1, $2)", "Steven Shan", comment)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (appl *appViews) GetComments(w http.ResponseWriter, r *http.Request) {

	rows, err := appl.Db.Query(context.Background(), "SELECT author, date, text FROM comments")
	if err != nil {
		log.Fatal(err)
	}
	type Comment struct {
		Author string
		Date   *time.Time
		First  string
		Rest   []string
	}
	var comments []Comment

	for rows.Next() {
		var comment Comment
		var text string
		if err := rows.Scan(&comment.Author, &comment.Date, &text); err != nil {
			log.Fatal(err)
		}

		paragraphs := strings.Split(text, "\n")
		comment.First = paragraphs[0]
		comment.Rest = paragraphs[1:]

		comments = append(comments, comment)
	}

	session, err := appl.Store.Get(r, USER_SESSION_NAME)
	if err != nil {
		log.Fatal(err)
	}
	username := ""
	if session_key, found := session.Values["session_key"]; found == true {
		row := appl.Db.QueryRow(context.Background(), "SELECT user_id FROM user_sessions WHERE session_key = $1", session_key)
		var user_id int
		err = row.Scan(&user_id)
		if err != nil {
			log.Fatal(err)
		}

		row = appl.Db.QueryRow(context.Background(), "SELECT username FROM users WHERE id = $1", user_id)
		err = row.Scan(&username)
		if err != nil {
			log.Fatal(err)
		}
	}

	appl.Templates.ExecuteTemplate(w, "index.tmpl",
		map[string]interface{}{
			"comments": comments,
			"username": username,
		},
	)
}
