package router

import (
	"database/sql"
	"net/http"

	"codisafish.eu/app/internal/auth"
	"codisafish.eu/app/internal/game"
	"codisafish.eu/app/internal/transactions"
)

type Deps struct {
	DB *sql.DB
}

func NewRouter(deps Deps) http.Handler {
	mux := http.NewServeMux()

	auth.Register(mux, deps.DB)
	game.Register(mux, deps.DB)
	transactions.Register(mux, deps.DB)

	return mux
}
