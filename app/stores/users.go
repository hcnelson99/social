package stores

import (
	"context"
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/scrypt"
)

/*
   Creates a new user.

   Returns the id (primary key) of the created user and a session generation
   number. See `NewUserSession` for details about the session generation.
*/
func (stores *Stores) NewUser(username, password string) (int, int, error) {
	const salt_len = 8
	salt := make([]byte, salt_len)
	n, err := rand.Read(salt)
	if err != nil {
		return -1, -1, err
	}
	if n != salt_len {
		// crypto/rand guarantees all bytes provided if err == nil
		return -1, -1, errors.New("bad salt length: crypto/rand package broken")
	}

	// Default values from crypto/scrypt documentation
	hash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return -1, -1, err
	}

	row := stores.db.QueryRow(context.Background(),
		"INSERT INTO users(username, password_hash, password_salt) VALUES ($1, $2, $3) RETURNING id, session_generation",
		username, hash, salt)

	var userId int
	var sessionGeneration int
	err = row.Scan(&userId, &sessionGeneration)
	if err != nil {
		return -1, -1, err
	}

	return userId, sessionGeneration, nil
}

/*
   Creates a new login session for a user.

   Returns a session_generation number. This generation number incremented
   to invalidate existing sessions. It is only valid when it is equal to the
   session_generation column in the user table.
*/
func (stores *Stores) NewUserSession(userId int) (int, error) {
	row, err := stores.db.Query(context.Background(),
		"SELECT session_generation FROM users WHERE id=$1",
		userId)
	if err != nil {
		return -1, err
	}

	var sessionGeneration int
	if err := row.Scan(&sessionGeneration); err != nil {
		return -1, err
	}

	return sessionGeneration, nil
}
