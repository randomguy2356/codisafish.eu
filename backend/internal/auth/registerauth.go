package auth

import (
	"database/sql"
	"net/http"
)

type LoginHandler struct {
	DB *sql.DB
}

type LogoutHandler struct {
	DB *sql.DB
}

type UserinfoHandler struct {
	DB *sql.DB
}

type RegisterHandler struct {
	DB *sql.DB
}

type ExistsHandler struct {
	DB *sql.DB
}

type checkRMTHandler struct {
	DB *sql.DB
}

func Register(mux *http.ServeMux, db *sql.DB) {
	loginhandler := &LoginHandler{DB: db}
	logouthandler := &LogoutHandler{DB: db}
	userinfohandler := &UserinfoHandler{DB: db}
	registerhandler := &RegisterHandler{DB: db}
	existshandler := &ExistsHandler{DB: db}
	checkrmthandler := &checkRMTHandler{DB: db}

	authMux := http.NewServeMux()

	authMux.Handle("/api/auth/login", loginhandler)
	authMux.Handle("/api/auth/logout", logouthandler)
	authMux.Handle("/api/auth/userinfo", userinfohandler)
	authMux.Handle("/api/auth/register", registerhandler)
	authMux.Handle("/api/auth/exists", existshandler)
	authMux.Handle("/api/auth/checkrmt", checkrmthandler)

	mux.Handle("/api/auth/", authMux)
}
