package stores

import (
	"context"
	"crypto/rand"
	"golang.org/x/crypto/scrypt"
	"log"
)

type User struct {
	Username string
}

/*
   Checks user session values in database to get the username of the currently
   logged in user and makes sure their session hasn't been invalidated.

   Returns a non-nil User pointer or an error.
*/
func (stores *Stores) CheckUserSession(userId, sessionGeneration int) *User {
	row := stores.db.QueryRow(context.Background(),
		"SELECT username FROM users WHERE id=$1 AND session_generation=$2",
		userId, sessionGeneration)

	var user User

	if err := row.Scan(&user.Username); err != nil {
		return nil
	}

	return &user
}

/*
   Creates a new user.

   Returns the id (primary key) of the created user and a session generation
   number. See `NewUserSession` for details about the session generation.
*/
func (stores *Stores) NewUser(username, password string) (int, int) {
	const salt_len = 8
	salt := make([]byte, salt_len)
	n, err := rand.Read(salt)
	if err != nil {
		log.Print("couldn't create a salt", err)
		return -1, -1
	}
	if n != salt_len {
		// crypto/rand guarantees all bytes provided if err == nil
		log.Print("bad salt length: crypto/rand package broken")
		return -1, -1
	}

	// Default values from crypto/scrypt documentation
	hash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		log.Print("couldn't hash password", err)
		return -1, -1
	}

	row := stores.db.QueryRow(context.Background(),
		"INSERT INTO users(username, password_hash, password_salt) VALUES ($1, $2, $3) RETURNING id, session_generation",
		username, hash, salt)

	var userId int
	var sessionGeneration int
	err = row.Scan(&userId, &sessionGeneration)
	if err != nil {
		log.Print("couldn't create user in database", err)
		return -1, -1
	}

	if userId < 0 || sessionGeneration < 0 {
		log.Print("database schema inconsistent: " + "user_id and session_generation should be >= 0")
	}

	return userId, sessionGeneration
}

/*
   Creates a new login session for a user.

   Returns a session_generation number. This generation number incremented
   to invalidate existing sessions. It is only valid when it is equal to the
   session_generation column in the user table.
*/
func (stores *Stores) NewUserSession(userId int) int {
	row := stores.db.QueryRow(context.Background(),
		"SELECT session_generation FROM users WHERE id=$1",
		userId)

	var sessionGeneration int
	if err := row.Scan(&sessionGeneration); err != nil {
		return -1
	}

	return sessionGeneration
}
