package stores

import (
	"context"
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/scrypt"
)

/*
   Creates a new user.
*/
func (stores *Stores) NewUser(username, password string) (int, error) {
	const salt_len = 8
	salt := make([]byte, salt_len)
	n, err := rand.Read(salt)
	if err != nil {
		return -1, err
	}
	if n != salt_len {
		// crypto/rand guarantees all bytes provided if err == nil
		return -1, errors.New("bad salt length: crypto/rand package broken")
	}

	// Default values from crypto/scrypt documentation
	hash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return -1, err
	}

	row := stores.db.QueryRow(context.Background(),
		"INSERT INTO users(username, password_hash, password_salt) VALUES ($1, $2, $3) RETURNING id",
		username, hash, salt)

	var userId int
	err = row.Scan(&userId)
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (stores *Stores) NewUserSession(userId int) (string, error) {
	return "testkey", nil
	/*
		const user_session_key_length = 32
		user_session_key := securecookie.GenerateRandomKey(user_session_key_length)
		_, err = app.Stores.Db.Exec(context.Background(),
			"INSERT INTO user_sessions(session_key, user_id) VALUES ($1, $2)",
			user_session_key, user_id)
		if err != nil {
			log.Fatal(err)
		}
	*/
}
