package types

import (
	"github.com/gorilla/sessions"
	"github.com/hcnelson99/social/app/stores"
	"html/template"
)

type App struct {
	Templates    *template.Template
	Stores       stores.Stores
	SessionStore sessions.Store
}
