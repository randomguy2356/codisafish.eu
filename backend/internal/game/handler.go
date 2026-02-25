package game

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"codisafish.eu/app/internal/httpx"
)

type Handler struct {
	DB *sql.DB
}

func Register(mux *http.ServeMux, db *sql.DB) {
	handler := &Handler{DB: db}

	mux.Handle("/api/game", handler)
}

type GameRequest struct {
	K      string `json:"k"`
	Aelo   string `json:"a_elo"`
	Belo   string `json:"b_elo"`
	ScoreA string `json:"score_a"`
}

type GameResponse struct {
	Aelo  string `json:"a_elo"`
	Belo  string `json:"b_elo"`
	Error string `json:"error,omitempty"`
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer request.Body.Close()

	request.Body = http.MaxBytesReader(writer, request.Body, 1<<20)

	var requestData GameRequest

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&requestData); err != nil {
		httpx.WriteJSON(writer, http.StatusBadRequest, GameResponse{
			Error: "invalid JSON request: " + err.Error(),
		})
		return
	}

	if requestData.K == "" ||
		requestData.Aelo == "" ||
		requestData.Belo == "" ||
		requestData.ScoreA == "" {

		httpx.WriteJSON(writer, http.StatusBadRequest, GameResponse{
			Error: "missing one or more required fields: k, a_elo, b_elo, score_a",
		})

		return
	}

	K, kerr := strconv.ParseFloat(requestData.K, 64)

	a_elo, aerr := strconv.ParseFloat(requestData.Aelo, 64)
	b_elo, berr := strconv.ParseFloat(requestData.Belo, 64)

	score_a, serr := strconv.ParseFloat(requestData.ScoreA, 64)

	if kerr != nil || aerr != nil || berr != nil || serr != nil {
		httpx.WriteJSON(writer, http.StatusBadRequest, GameResponse{
			Error: "the fields: k, a_elo, b_elo, score; must be valid numbers",
		})
		return
	}

	a_elo_new, b_elo_new := calculate_elo(K, a_elo, b_elo, score_a, handler.DB)

	httpx.WriteJSON(writer, http.StatusOK, GameResponse{
		Aelo: strconv.FormatFloat(a_elo_new, 'f', 2, 64),
		Belo: strconv.FormatFloat(b_elo_new, 'f', 2, 64),
	})
}

func calculate_elo(K float64, Aelo float64, Belo float64, scoreA float64, db *sql.DB) (float64, float64) {
	ExpectedScoreA := 1 / (1 + (math.Pow(10, (Belo-Aelo)/400.0)))

	Adelta := K * (scoreA - ExpectedScoreA)

	maxDelta := Belo - 10
	minDelta := -(Aelo - 10)

	if Adelta > maxDelta {
		Adelta = maxDelta
	}
	if Adelta < minDelta {
		Adelta = minDelta
	}

	Adelta = math.Round(Adelta*100) / 100

	log_elo(Aelo, Belo, Adelta, Aelo+Adelta, Belo-Adelta, db)

	return Aelo + Adelta, Belo - Adelta
}

// DB logging

func log_elo(a float64, b float64, delta float64, a_new float64, b_new float64, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	query := "INSERT INTO transactions (a, b, delta, a_new, b_new) VALUES (?, ?, ?, ?, ?)"

	_, err := db.ExecContext(ctx, query, a, b, delta, a_new, b_new)

	if err != nil {
		log.Println("error while inserting: ", err.Error())
		return err
	}

	return nil
}
