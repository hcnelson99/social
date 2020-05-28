// Main code for database stores. This is an abstraction between the business
// logic in the views and the SQL queries.
package stores

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

// Establishes a connection pool to the database
func connectToDb() (*pgxpool.Pool, error) {
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

type Stores struct {
	db *pgxpool.Pool
}

// Initializes a connection to database and returns an abstraction for
// interacting to the database. This abstraction is used instead of directly
// writing and executing SQL.
func (stores *Stores) Init() error {
	db, err := connectToDb()
	stores.db = db
	return err
}

// Cleans up stores and closes connection to database.
func (stores *Stores) Close() {
	stores.db.Close()
}
