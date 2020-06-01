package app

import (
	// "crypto/subtle"
	"encoding/base64"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/hcnelson99/social/app/types"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	SESSION_KEY_LENGTH = 32
)

func Start(addr string) {
	var app types.App

	app.Templates = template.Must(template.ParseGlob("./app/templates/*.tmpl"))

	app.SchemaDecoder = schema.NewDecoder()
	if app.SchemaDecoder == nil {
		log.Fatal("Could not initialize schema decoder")
	}

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

	app.SessionStore = sessions.NewCookieStore(session_key)

	err = app.Stores.Init()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer app.Stores.Close()

	http.Handle("/", GetRouter(&app))

	log.Printf("Starting server on address %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
