package httpx

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(writer http.ResponseWriter, status int, response any) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	err := json.NewEncoder(writer).Encode(response)
	return err
}
