package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var db *pgxpool.Pool

func badRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "bad request")
}

func internalServerError(c *gin.Context) {
	c.String(http.StatusInternalServerError, "internal server error")
}

func newComment(c *gin.Context) {
	comment, success := c.GetPostForm("comment")
	if !success {
		badRequest(c)
	}

	_, err := db.Query(context.Background(), "INSERT INTO comments(author, text) VALUES ($1, $2)", "Steven Shan", comment)
	if err != nil {
		log.Fatal(err)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func index(c *gin.Context) {
	rows, err := db.Query(context.Background(), "SELECT author, date, text FROM comments")
	if err != nil {
		log.Fatal(err)
		badRequest(c)
		return
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

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"comments": comments,
	})
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	db, err = pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", index)
	r.POST("/", newComment)
	r.Static("/static", "./static")
	r.Run()
}
