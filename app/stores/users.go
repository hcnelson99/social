package stores

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"golang.org/x/crypto/scrypt"
	"log"
)

type User struct {
	Username string
	UserId   int
}

const SALT_LENGTH = 8

type AuthStatus int

const (
	AUTH_VALIDATED AuthStatus = iota
	AUTH_REJECTED
	AUTH_ERROR
)

/*
   Initializes a `User` object. Never returns nil.
*/
func newUser(username string, userId int) *User {
	return &User{
		Username: username,
		UserId:   userId,
	}
}

/*
   Generates a salt and hashes a given password with it.

   Returns a concatentation of the salt and the hashed password so only one
   column is required to store it.
*/
func hashPassword(password string) []byte {
	salt := make([]byte, SALT_LENGTH)
	n, err := rand.Read(salt)
	if err != nil {
		log.Print("couldn't create a salt", err)
		return nil
	}
	if n != SALT_LENGTH {
		// crypto/rand guarantees all bytes provided if err == nil
		log.Print("bad salt length: crypto/rand package broken")
		return nil
	}

	// Default values from crypto/scrypt documentation
	// https://godoc.org/golang.org/x/crypto/scrypt
	hash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		log.Print("couldn't hash password", err)
		return nil
	}

	return append(salt, hash...)
}

/*
   Checks if a plaintext password matches a hashed password.

   Uses the salt that is embedded in the hash by `hashPassword`.
*/
func checkPassword(hash []byte, password string) AuthStatus {
	if len(hash) < SALT_LENGTH {
		log.Printf(
			"invalid password hash (length %d < %d)", len(hash), SALT_LENGTH)
		return AUTH_ERROR
	}
	salt := hash[:SALT_LENGTH]
	passwordHash := hash[SALT_LENGTH:]

	// Default values from crypto/scrypt documentation
	// https://godoc.org/golang.org/x/crypto/scrypt
	newHash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		log.Print("couldn't hash password", err)
		return AUTH_ERROR
	}

	if subtle.ConstantTimeCompare(passwordHash, newHash) == 1 {
		return AUTH_VALIDATED
	} else {
		return AUTH_REJECTED
	}
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

	var username string

	if err := row.Scan(&username); err != nil {
		return nil
	}

	return newUser(username, userId)
}

/*
   Creates a new user.

   Returns a pointer to a `User` struct and a session generation number upon
    success. See `NewUserSession` for details about the session generation.
*/
func (stores *Stores) NewUser(username, password string) (*User, int) {
	hash := hashPassword(password)
	if hash == nil {
		return nil, -1
	}

	row := stores.db.QueryRow(context.Background(),
		"INSERT INTO users(username, password) VALUES ($1, $2) RETURNING id, session_generation",
		username, hash)

	var userId int
	var sessionGeneration int
	if err := row.Scan(&userId, &sessionGeneration); err != nil {
		log.Print("couldn't create user in database", err)
		return nil, -1
	}

	if userId < 0 || sessionGeneration < 0 {
		log.Print("database inconsistent: user_id and session_generation should be >= 0")
		return nil, -1
	}

	return newUser(username, userId), sessionGeneration
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

/*
   Authenticates a user.

   If the returned `AuthStatus` is `AUTH_VALIDATED`, the user was successfully
   authenticated and a `User` object and session number will be returned.
*/
func (stores *Stores) Login(username, password string) (*User, int, AuthStatus) {
	row := stores.db.QueryRow(context.Background(),
		"SELECT id, password, session_generation FROM users WHERE username=$1",
		username)

	var userId int
	var passwordHash []byte
	var sessionGeneration int
	if err := row.Scan(&userId, &passwordHash, &sessionGeneration); err != nil {
		return nil, -1, AUTH_ERROR
	}

	status := checkPassword(passwordHash, password)
	if status == AUTH_VALIDATED {
		return newUser(username, userId), sessionGeneration, status
	}

	return nil, -1, status
}
