package views

import (
	"context"
	"crypto/rand"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/scrypt"
	"log"
	"net/http"
)

func (appl *appViews) GetLogin(w http.ResponseWriter, r *http.Request) {
	appl.Templates.ExecuteTemplate(w, "login.tmpl", nil)
}

func (appl *appViews) PostLogin(w http.ResponseWriter, r *http.Request) {
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

func (appl *appViews) Register(w http.ResponseWriter, r *http.Request) {
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

	row := appl.Db.QueryRow(context.Background(),
		"INSERT INTO users(username, password_hash, password_salt) VALUES ($1, $2, $3) RETURNING id",
		username, hash, salt)
	var user_id int
	err = row.Scan(&user_id)
	if err != nil {
		log.Fatal(err)
	}

	const user_session_key_length = 32
	user_session_key := securecookie.GenerateRandomKey(user_session_key_length)
	_, err = appl.Db.Exec(context.Background(),
		"INSERT INTO user_sessions(session_key, user_id) VALUES ($1, $2)",
		user_session_key, user_id)
	if err != nil {
		log.Fatal(err)
	}

	session, err := appl.Store.Get(r, USER_SESSION_NAME)
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
