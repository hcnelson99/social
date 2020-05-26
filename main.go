package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var t *template.Template
var db *pgxpool.Pool

func badRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "400 bad request")
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 page not found")
}

func getPostFormValue(r *http.Request, key string) (string, bool) {
	r.ParseForm()
	values, success := r.PostForm[key]
	if !success || len(values) != 1 {
		return "", false
	}
	return values[0], true
}

func putIndex(w http.ResponseWriter, r *http.Request) {
	comment, success := getPostFormValue(r, "comment")
	if !success {
		badRequest(w)
		return
	}

	_, err := db.Query(context.Background(), "INSERT INTO comments(author, text) VALUES ($1, $2)", "Steven Shan", comment)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusFound)
}

func getIndex(w http.ResponseWriter) {
	rows, err := db.Query(context.Background(), "SELECT author, date, text FROM comments")
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

	t.ExecuteTemplate(w, "index.tmpl",
		map[string]interface{}{
			"comments": comments,
		},
	)
}

func index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getIndex(w)
		break
	case "POST":
		putIndex(w, r)
		break
	default:
		notFound(w)
		break
	}
}

func getDb() (*pgxpool.Pool, error) {
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	url := os.Getenv("DATABASE_URL")
	name := os.Getenv("DATABASE_NAME")

	var connUrl string
	if username == "" && password == "" {
		connUrl = fmt.Sprintf(
			"postgres://%s/%s",
			url,
			name,
		)
	} else {
		connUrl = fmt.Sprintf(
			"postgres://%s:%s@%s/%s",
			username,
			password,
			url,
			name,
		)
	}

	log.Printf("Connecting to database: postgres://%s/%s\n", url, name)

	return pgxpool.Connect(context.Background(), connUrl)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	t = template.Must(template.ParseGlob("./templates/*"))

	var err error

	db, err = getDb()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", index)

	http.Handle("/", r)
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
