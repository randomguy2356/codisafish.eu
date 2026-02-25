package transactions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"codisafish.eu/app/internal/httpx"
)

type Handler struct {
	DB *sql.DB
}

func Register(mux *http.ServeMux, db *sql.DB) {
	handler := &Handler{DB: db}

	mux.Handle("/api/transactions", handler)
}

type TransactionsRequest struct {
	LastIndex int `json:"last_index"`
	Limit     int `json:"limit"`
}

type TransactionsResponse struct {
	LastIndex int      `json:"last_index"`
	Entries   []string `json:"entries"`
	Error     string   `json:"error,omitempty"`
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer request.Body.Close()

	request.Body = http.MaxBytesReader(writer, request.Body, 1<<20)

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	var requestData TransactionsRequest

	if err := decoder.Decode(&requestData); err != nil {
		httpx.WriteJSON(writer, http.StatusBadRequest, TransactionsResponse{
			Error: "bad request",
		})
		return
	}

	query := "SELECT * FROM transactions WHERE id > ? ORDER BY id DESC LIMIT ?"

	ctx := request.Context()

	rows, err := handler.DB.QueryContext(ctx, query, requestData.LastIndex, requestData.Limit)

	if err != nil {
		httpx.WriteJSON(writer, http.StatusInternalServerError, TransactionsResponse{
			Error: "database querry failed",
		})
		return
	}
	defer rows.Close()

	var entries []string

	for rows.Next() {
		var id int
		var a_elo float64
		var b_elo float64
		var delta float64
		var a_elo_new float64
		var b_elo_new float64

		if err := rows.Scan(&id, &a_elo, &b_elo, &delta, &a_elo_new, &b_elo_new); err != nil {
			httpx.WriteJSON(writer, http.StatusInternalServerError, TransactionsResponse{
				Error: "unexpected row format",
			})
			return
		}

		entries = append(entries, fmt.Sprintf("%v:  | a: %v | b: %v | d: %v | a_new: %v | b_new: %v |", id, a_elo, b_elo, delta, a_elo_new, b_elo_new))
	}

	if err := rows.Err(); err != nil {
		httpx.WriteJSON(writer, http.StatusInternalServerError, TransactionsResponse{
			Error: "error while reading database",
		})
		return
	}

	response := TransactionsResponse{
		Entries: entries,
	}

	httpx.WriteJSON(writer, http.StatusOK, response)
}
