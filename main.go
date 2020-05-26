package main

import (
	"context"
	"crypto/rand"
	// "crypto/subtle"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/scrypt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var t *template.Template
var db *pgxpool.Pool

func httpError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func getPostFormValue(r *http.Request, key string) (string, bool) {
	r.ParseForm()
	values, success := r.PostForm[key]
	if !success || len(values) != 1 {
		return "", false
	}
	return values[0], true
}

func postComment(w http.ResponseWriter, r *http.Request) {
	comment, success := getPostFormValue(r, "comment")
	if !success {
		httpError(w, http.StatusBadRequest)
		return
	}

	_, err := db.Query(context.Background(), "INSERT INTO comments(author, text) VALUES ($1, $2)", "Steven Shan", comment)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusFound)
}

func getComments(w http.ResponseWriter, r *http.Request) {
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

func getLogin(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "login.tmpl", nil)
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	username, success := getPostFormValue(r, "username")
	if !success {
		httpError(w, http.StatusBadRequest)
		return
	}
	password, success := getPostFormValue(r, "password")
	if !success {
		httpError(w, http.StatusBadRequest)
		return
	}

	log.Println(username, password)
}

func register(w http.ResponseWriter, r *http.Request) {
	username, success := getPostFormValue(r, "username")
	if !success {
		httpError(w, http.StatusBadRequest)
		return
	}
	password, success := getPostFormValue(r, "password")
	if !success {
		httpError(w, http.StatusBadRequest)
		return
	}

	// TODO: sanity check password and username

	salt_len := 8
	salt := make([]byte, salt_len)
	n, err := rand.Read(salt)
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
	if n != 8 {
		// crypto/rand guarantees all bytes provided if err == nil
		log.Fatal("crypto/rand broken")
	}

	// Default values from crypto/scrypt documentation
	hash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}

	_, err = db.Query(context.Background(), "INSERT INTO users(username, password_hash, password_salt) VALUES ($1, $2, $3)", username, hash, salt)
	if err != nil {
		// TODO: even though we have a uniqueness constraint on usernames, this
		// Query doesnt seem to give an error when we try to insert a duplicate
		// username. why is this?? We would like for it to do so so we do not
		// have a TOCTOU with checking for uniqueness of username (since
		// scrypting is probably more expensive than a db query)
		log.Fatal(err)
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

	t = template.Must(template.ParseGlob("./templates/*.tmpl"))

	var err error

	db, err = getDb()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// TODO: add CSRF everywhere!
	r := mux.NewRouter()
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", getComments).Methods("GET")
	r.HandleFunc("/comment", postComment).Methods("POST")
	r.HandleFunc("/login", getLogin).Methods("GET")
	r.HandleFunc("/login/post", postLogin).Methods("POST")
	r.HandleFunc("/login/register", register).Methods("POST")

	http.Handle("/", r)
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
