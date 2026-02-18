package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

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

func writeJSON(writer http.ResponseWriter, status int, response any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(response)
}

func calculate_elo(K float64, Aelo float64, Belo float64, scoreA float64) (float64, float64) {
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

	return Aelo + Adelta, Belo - Adelta
}

func gameHandler(writer http.ResponseWriter, request *http.Request) {
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
		writeJSON(writer, http.StatusBadRequest, GameResponse{
			Error: "invalid JSON request: " + err.Error(),
		})
		return
	}

	if requestData.K == "" ||
		requestData.Aelo == "" ||
		requestData.Belo == "" ||
		requestData.ScoreA == "" {

		writeJSON(writer, http.StatusBadRequest, GameResponse{
			Error: "missing one or more required fields: k, a_elo, b_elo, score_a",
		})

		return
	}

	K, kerr := strconv.ParseFloat(requestData.K, 64)

	a_elo, aerr := strconv.ParseFloat(requestData.Aelo, 64)
	b_elo, berr := strconv.ParseFloat(requestData.Belo, 64)

	score_a, serr := strconv.ParseFloat(requestData.ScoreA, 64)

	if kerr != nil || aerr != nil || berr != nil || serr != nil {
		writeJSON(writer, http.StatusBadRequest, GameResponse{
			Error: "the fields: k, a_elo, b_elo, score; must be valid numbers",
		})
		return
	}

	a_elo_new, b_elo_new := calculate_elo(K, a_elo, b_elo, score_a)

	writeJSON(writer, http.StatusOK, GameResponse{
		Aelo: strconv.FormatFloat(a_elo_new, 'f', 2, 64),
		Belo: strconv.FormatFloat(b_elo_new, 'f', 2, 64),
	})
}

func DB_connect() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")

	if host == "" || port == "" || name == "" || user == "" || pass == "" {
		return nil, fmt.Errorf("missing required DB environment variables")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, pass, host, port, name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func dbtest() {
	db, err := DB_connect()

	if err != nil {
		log.Fatalf("we fucked: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var version string
	err = db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatalf("we fucked up againg: %v", err)
	}

	println("Connected to: ", version)
}

func main() {

	dbtest()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/game", gameHandler)

	log.Println("listening on :5000")
	log.Fatal(http.ListenAndServe(":5000", mux))
}
