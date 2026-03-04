package auth

import (
	_ "crypto/rand"
	"encoding/json"
	"net/http"

	"codisafish.eu/app/internal/httpx"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Error string `json:"error,omitempty"`
}

func (handler *LoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer request.Body.Close()

	request.Body = http.MaxBytesReader(writer, request.Body, 1<<20)

	var requestData loginRequest

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&requestData); err != nil {
		httpx.WriteJSON(writer, http.StatusBadRequest, loginResponse{
			Error: "invalid JSON request: " + err.Error(),
		})
		return
	}

	if requestData.Username == "" || requestData.Password == "" {
		httpx.WriteJSON(writer, http.StatusBadRequest, loginResponse{
			Error: "missing required fields",
		})
		return
	}

	wrong := loginResponse{
		Error: "wrong username or password",
	}

	valid, err := ValidateUser(requestData.Username, requestData.Password, request.Context(), handler.DB)

	if err != nil {
		httpx.WriteJSON(writer, http.StatusInternalServerError, loginResponse{
			Error: "internal server error: " + err.Error(),
		})
		return
	}

	if !valid {
		httpx.WriteJSON(writer, http.StatusOK, wrong)
		return
	}

	CreateSession(requestData.Username, writer)

	writer.WriteHeader(http.StatusOK)
}
