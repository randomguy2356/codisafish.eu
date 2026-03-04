package auth

import (
	"context"
	"database/sql"
	"net/http"

	"codisafish.eu/app/internal/httpx"
)

type existsResponse struct {
	Username *bool `json:"username,omitempty"`
	Email    *bool `json:"email,omitempty"`
}

func (handler *ExistsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	queries := request.URL.Query()

	username := queries.Get("username")
	email := queries.Get("email")

	response := existsResponse{}

	userexists, err := UserExists(username, handler.DB, request.Context())
	if err != nil {
		http.Error(writer, "internal friggin error. some db stuff", http.StatusInternalServerError)
		return
	}

	emailexists, err := EmailExists(email, handler.DB, request.Context())
	if err != nil {
		http.Error(writer, "internal friggin error. some db stuff", http.StatusInternalServerError)
		return
	}

	response.Username = userexists
	response.Email = emailexists
	httpx.WriteJSON(writer, http.StatusOK, response)
}

func UserExists(username string, db *sql.DB, contex context.Context) (*bool, error) {
	if username == "" {
		return nil, nil
	}
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE username=?)"

	row := db.QueryRowContext(contex, query, username)

	var exist bool

	err := row.Scan(&exist)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return &exist, err
}

func EmailExists(email string, db *sql.DB, contex context.Context) (*bool, error) {
	if email == "" {
		return nil, nil
	}
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email=?)"

	row := db.QueryRowContext(contex, query, email)

	var exist bool

	err := row.Scan(&exist)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return &exist, err
}
