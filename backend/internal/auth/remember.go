package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
)

func HashToken(token string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
}

func CreateRememberToken(username string, db *sql.DB, context context.Context, writer http.ResponseWriter) error {
	query := "SELECT id FROM users WHERE username=?"

	row := db.QueryRowContext(context, query, username)
	var user_id int
	err := row.Scan(&user_id)
	if err != nil {
		return err
	}

	query = "INSERT INTO remember_me (user_id, token_hash, expires_at) VALUES (?, ?, ?)"

	var token string

	for {
		token = rand.Text()
		token_hash := HashToken(token)

		_, err = db.ExecContext(context, query, user_id, token_hash, time.Now().Add(30*24*time.Hour))
		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				continue
			}
			return err
		}
		break
	}

	SetRememberCookie(writer, token, 30*24*time.Hour)
	return nil
}

func SetRememberCookie(writer http.ResponseWriter, token string, maxAge time.Duration) {
	http.SetCookie(writer, &http.Cookie{
		Name:     "rmt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(maxAge),
	})
}
