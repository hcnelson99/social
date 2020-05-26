package main

import (
	"context"
	"crypto/rand"
	// "crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/scrypt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const user_session_name = "login"

var t *template.Template
var db *pgxpool.Pool
var store sessions.Store

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

	_, err := db.Exec(context.Background(), "INSERT INTO comments(author, text) VALUES ($1, $2)", "Steven Shan", comment)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
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

	session, err := store.Get(r, user_session_name)
	if err != nil {
		log.Fatal(err)
	}
	username := ""
	if session_key, found := session.Values["session_key"]; found == true {
		row := db.QueryRow(context.Background(), "SELECT user_id FROM user_sessions WHERE session_key = $1", session_key)
		var user_id int
		err = row.Scan(&user_id)
		if err != nil {
			log.Fatal(err)
		}

		row = db.QueryRow(context.Background(), "SELECT username FROM users WHERE id = $1", user_id)
		err = row.Scan(&username)
		if err != nil {
			log.Fatal(err)
		}
	}

	t.ExecuteTemplate(w, "index.tmpl",
		map[string]interface{}{
			"comments": comments,
			"username": username,
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

	const salt_len = 8
	salt := make([]byte, salt_len)
	n, err := rand.Read(salt)
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
	if n != salt_len {
		// crypto/rand guarantees all bytes provided if err == nil
		log.Fatal("crypto/rand broken")
	}

	// Default values from crypto/scrypt documentation
	hash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}

	row := db.QueryRow(context.Background(),
		"INSERT INTO users(username, password_hash, password_salt) VALUES ($1, $2, $3) RETURNING id",
		username, hash, salt)
	var user_id int
	err = row.Scan(&user_id)
	if err != nil {
		log.Fatal(err)
	}

	const user_session_key_length = 32
	user_session_key := securecookie.GenerateRandomKey(user_session_key_length)
	_, err = db.Exec(context.Background(),
		"INSERT INTO user_sessions(session_key, user_id) VALUES ($1, $2)",
		user_session_key, user_id)
	if err != nil {
		log.Fatal(err)
	}

	session, err := store.Get(r, user_session_name)
	if err != nil {
		log.Fatal(err)
	}
	session.Values["session_key"] = user_session_key
	err = session.Save(r, w)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
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

	// To generate a session key, run:
	//
	//     func main() {
	//         fmt.Println(base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)))
	//     }

	const session_key_length = 32

	session_key_b64 := os.Getenv("SESSION_KEY")
	if session_key_b64 == "" {
		log.Fatal("Missing session key in environment")
	}
	session_key, err := base64.StdEncoding.DecodeString(session_key_b64)
	if err != nil {
		log.Fatal(err)
	}
	if len(session_key) != session_key_length {
		log.Fatal("Invalid session key length")
	}
	store = sessions.NewCookieStore(session_key)

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
