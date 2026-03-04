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
	"github.com/alexedwards/argon2id"
	_ "github.com/alexedwards/argon2id"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jaswdr/faker/v2"
	_ "github.com/jaswdr/faker/v2"
)

func main() {
	db, err := DB_connect()

	if err != nil {
		log.Fatal("error while connecting to database: ", err)
	}

	//if err := add_more_users(db); err != nil {
	//	log.Fatal("failed to add example users: ", err.Error())
	//}

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

func add_more_users(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	stmt1, err := tx.PrepareContext(ctx, "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)")

	if err != nil {
		return err
	}

	defer stmt1.Close()

	stmt2, err := tx.PrepareContext(ctx, "INSERT INTO username_password (username, password) VALUES (?, ?)")

	if err != nil {
		return err
	}

	defer stmt2.Close()

	for range 100 {
		faker := faker.New()
		username := faker.Internet().User()
		email := faker.Internet().Email()
		password := faker.Internet().Password()
		password_hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		if err != nil {
			return err
		}

		if _, err = stmt1.ExecContext(ctx, username, email, password_hash); err != nil {
			return err
		}

		if _, err = stmt2.ExecContext(ctx, username, password); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
