package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alexedwards/argon2id"
)

type registerRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (handler *RegisterHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed, use Post", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	var requestData registerRequest

	if err := decoder.Decode(&requestData); err != nil {
		http.Error(writer, "invalid json request", http.StatusBadRequest)
		return
	}

	if requestData.Username == "" || requestData.Email == "" || requestData.Password == "" {
		http.Error(writer, "missing required fields", http.StatusBadRequest)
		return
	}

	user_exists, err := UserExists(requestData.Username, handler.DB, request.Context())
	if err != nil || user_exists == nil {
		http.Error(writer, "internal error", http.StatusInternalServerError)
		return
	}

	email_exists, err := EmailExists(requestData.Email, handler.DB, request.Context())
	if err != nil || email_exists == nil {
		http.Error(writer, "internal error", http.StatusInternalServerError)
		return
	}

	if *user_exists || *email_exists {
		http.Error(writer, "username or email taken", http.StatusConflict)
		return
	}

	err = CreateUser(requestData.Username, requestData.Email, requestData.Password, handler.DB, writer, request.Context())
	if err != nil {
		log.Default().Println(fmt.Errorf("CreateUser: %w", err))
		http.Error(writer, "internal error", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func CreateUser(username string, email string, password string, db *sql.DB, writer http.ResponseWriter, contex context.Context) error {
	query := "INSERT INTO users (username, email, password_hash) values (?, ?, ?)"

	password_hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return fmt.Errorf("argon2id: %w", err)
	}

	if _, err := db.ExecContext(contex, query, username, email, password_hash); err != nil {
		return fmt.Errorf("db exec: %w", err)
	}

	CreateSession(username, writer)
	return nil
}
