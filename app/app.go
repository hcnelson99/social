package app

import (
	"context"
	// "crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/hcnelson99/social/app/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	SESSION_KEY_LENGTH = 32
)

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

func Start(addr string) {
	var app types.App

	app.Templates = template.Must(template.ParseGlob("./app/templates/*.tmpl"))

	session_key_b64 := os.Getenv("SESSION_KEY")
	if session_key_b64 == "" {
		log.Fatal("Missing session key in environment")
	}
	session_key, err := base64.StdEncoding.DecodeString(session_key_b64)
	if err != nil {
		log.Fatal(err)
	}
	if len(session_key) != SESSION_KEY_LENGTH {
		log.Fatalf("Invalid session key length.")
	}

	app.Store = sessions.NewCookieStore(session_key)

	app.Db, err = getDb()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer app.Db.Close()

	http.Handle("/", GetRouter(&app))

	log.Printf("Starting server on address %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
