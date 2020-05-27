package types

import (
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v4/pgxpool"
	"html/template"
)

type App struct {
	Templates *template.Template
	Db        *pgxpool.Pool
	Store     sessions.Store
}
