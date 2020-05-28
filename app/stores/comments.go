package stores

import (
	"context"
	"log"
	"strings"
	"time"
)

type Comment struct {
	Author string
	Date   *time.Time
	First  string
	Rest   []string
}

func (stores *Stores) GetAllComments() ([]Comment, error) {
	rows, err := stores.db.Query(context.Background(),
		"SELECT author, date, text FROM comments")
	if err != nil {
		return []Comment{}, err
	}

	var comments []Comment

	for rows.Next() {
		var comment Comment
		var text string

		if err := rows.Scan(&comment.Author, &comment.Date, &text); err != nil {
			log.Fatal("Invalid database schema.", err)
		}

		paragraphs := strings.Split(text, "\n")
		comment.First = paragraphs[0]
		comment.Rest = paragraphs[1:]

		comments = append(comments, comment)
	}

	return comments, nil
}

func (stores *Stores) NewComment(comment string) error {
	_, err := stores.db.Exec(
		context.Background(),
		"INSERT INTO comments(author, text) VALUES ($1, $2)",
		"Steven Shan", comment)
	return err
}
