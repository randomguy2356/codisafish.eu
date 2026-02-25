package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"codisafish.eu/app/internal/httpx"
	"github.com/alexedwards/argon2id"
)

type Handler struct {
	db *sql.DB
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Msg   string `json:"msg"`
	Error string `json:"error,omitempty"`
}

func Register(mux *http.ServeMux, db *sql.DB) {
	handler := &Handler{db: db}

	mux.Handle("/api/auth", handler)
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer request.Body.Close()

	request.Body = http.MaxBytesReader(writer, request.Body, 1<<20)

	var requestData authRequest

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&requestData); err != nil {
		httpx.WriteJSON(writer, http.StatusBadRequest, authResponse{
			Error: "invalid JSON request: " + err.Error(),
		})
		return
	}

	if requestData.Username == "" || requestData.Password == "" {
		httpx.WriteJSON(writer, http.StatusBadRequest, authResponse{
			Error: "missing required fields",
		})
		return
	}

	query := "SELECT password_hash FROM users WHERE username = ?"

	var password_hash string

	row := handler.db.QueryRowContext(request.Context(), query, requestData.Username)

	err := row.Scan(&password_hash)
	if err == sql.ErrNoRows {
		httpx.WriteJSON(writer, http.StatusOK, authResponse{
			Error: "no user with this username",
		})
		return
	}

	matches, err := argon2id.ComparePasswordAndHash(requestData.Password, password_hash)
	if err != nil {
		httpx.WriteJSON(writer, http.StatusInternalServerError, authResponse{
			Error: "error while comparing passwords: " + err.Error(),
		})
		return
	}

	httpx.WriteJSON(writer, http.StatusOK, authResponse{
		Msg: "matches: " + strconv.FormatBool(matches),
	})
}
