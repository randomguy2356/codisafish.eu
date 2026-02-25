package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codisafish.eu/app/internal/router"
	_ "github.com/alexedwards/argon2id"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := DB_connect()

	if err != nil {
		log.Fatal("error while connecting to database: ", err)
	}

	//	mux := http.NewServeMux()
	//	mux.HandleFunc("/api/game", gameHandler)
	//	mux.HandleFunc("/api/transactions", fetch_transactions)

	handler := router.NewRouter(router.Deps{DB: db})

	server := http.Server{
		Addr:              ":5000",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = server.Shutdown(shutCtx)
		_ = db.Close()
	}()

	log.Println("listening on :5000")
	log.Fatal(server.ListenAndServe())
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
