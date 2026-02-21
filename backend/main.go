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

	_ "github.com/alexedwards/argon2id"
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

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/api/game", gameHandler)

	log.Println("listening on :5000")
	log.Fatal(http.ListenAndServe(":5000", mux))
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

	log_elo(Aelo, Belo, Adelta, Aelo+Adelta, Belo-Adelta)

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

// DB logging

func log_elo(a float64, b float64, delta float64, a_new float64, b_new float64) error {
	db, err := DB_connect()

	if err != nil {
		log.Println("failed to connect to database: ", err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO transactions (a, b, delta, a_new, b_new) VALUES (?, ?, ?, ?, ?)"

	_, err = db.ExecContext(ctx, query, a, b, delta, a_new, b_new)

	if err != nil {
		log.Println("error while inserting: ", err.Error())
		return err
	}

	return nil
}

//--DB stuff

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
		log.Fatalf("error: %v", err)
	}
	defer db.Close()

	hash := "$argon2id$v=19$m=65536,t=1,p=4$kaoO+n/58AUyNzG/17SVkg$fB0iMWEU1Zwro8QPTexI4YfkJ9FQ6nUh7TQSSR0LDPI"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id string
	var username string
	var email string
	var password_hash string
	var created_at string
	err = db.QueryRowContext(ctx, "SELECT * FROM users").Scan(&id, &username, &email, &password_hash, &created_at)
	if err == sql.ErrNoRows {
		db.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
		db.Exec(fmt.Sprintf("INSERT INTO users (id, username, email, password_hash) values (1, 'admin', 'trocheniefajnie@gmail.com', '%s')", hash))
	}
	err = db.QueryRowContext(ctx, "SELECT * FROM users WHERE id=1").Scan(&id, &username, &email, &password_hash, &created_at)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	println("| ", id, " | ", username, " | ", email, " | ", password_hash, " | ", created_at, " |")
}
