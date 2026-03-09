package auth

import (
	"context"
	"database/sql"
	"net/http"
)

func (handler *checkRMTHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	rmt, err := request.Cookie("rmt")
	if err != nil {
		if err == http.ErrNoCookie {
			writer.WriteHeader(http.StatusOK)
			return
		}
		http.Error(writer, "cookies bad", http.StatusBadRequest)
	}
	username, error := ValidateRMT(rmt.Value, handler.DB, request.Context())
	if error != nil {
		http.Error(writer, "internal server error", http.StatusInternalServerError)
		return
	}
	if username == nil {
		writer.WriteHeader(http.StatusOK)
		return
	}
	CreateSession(*username, writer)
	writer.WriteHeader(http.StatusOK)
}

func ValidateRMT(rmt string, db *sql.DB, context context.Context) (*string, error) {
	token_hash := HashToken(rmt)
	query := "SELECT users.username FROM remember_me JOIN users ON remember_me.user_id=users.id WHERE token_hash=?"

	row := db.QueryRowContext(context, query, token_hash)

	var username *string = nil
	if err := row.Scan(username); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return username, nil
}

func CheckRMTCookie(db *sql.DB, writer http.ResponseWriter, request *http.Request) (string, error) {
	rmt, err := request.Cookie("rmt")
	if err != http.ErrNoCookie {
		if err != nil {
			return "", err
		}
		username, err := ValidateRMT(rmt.Value, db, request.Context())
		if err != nil {
			return "", err
		}
		if username != nil {
			sid := CreateSession(*username, writer)
			return sid, nil
		}
	}
	return "", nil
}
