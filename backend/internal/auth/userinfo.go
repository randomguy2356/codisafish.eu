package auth

import (
	"context"
	"database/sql"
	"net/http"

	"codisafish.eu/app/internal/httpx"
)

type userinfoResponse struct {
	Exists    bool    `json:"exists"`
	Username  string  `json:"username,omitempty"`
	Email     string  `json:"email,omitempty"`
	Elo       float64 `json:"elo,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
	Error     string  `json:"error,omitempty"`
}

func (handler *UserinfoHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}

	clientSID, err := request.Cookie("sid")

	if err != nil {
		if err == http.ErrNoCookie {
			httpx.WriteJSON(writer, http.StatusOK, userinfoResponse{
				Exists: false,
			})
			return
		}
		httpx.WriteJSON(writer, http.StatusBadRequest, userinfoResponse{
			Error: "cookie read error: " + err.Error(),
		})
		return
	}

	session := ValidateSID(clientSID.Value, writer)

	if session == nil {
		httpx.WriteJSON(writer, http.StatusOK, userinfoResponse{
			Exists: false,
		})
		return
	}

	userinfo, err := GetUserInfo(session.User, handler.DB, request.Context())

	if err != nil {
		httpx.WriteJSON(writer, http.StatusInternalServerError, userinfoResponse{
			Error: "failed while doing db stuff: " + err.Error(),
		})
		return
	}

	if userinfo == nil {
		httpx.WriteJSON(writer, http.StatusOK, userinfoResponse{
			Exists: false,
		})
	}

	httpx.WriteJSON(writer, http.StatusOK, userinfoResponse{
		Exists:    true,
		Username:  userinfo.Username,
		Email:     userinfo.Email,
		Elo:       userinfo.Elo,
		CreatedAt: userinfo.CreatedAt.Time.String(),
	})

}

type DBUserInfo struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	Elo          float64
	CreatedAt    sql.NullTime
}

func GetUserInfo(username string, db *sql.DB, context context.Context) (userinfo *DBUserInfo, err error) {
	query := "SELECT * FROM users WHERE username=? ORDER BY id"

	row := db.QueryRowContext(context, query, username)

	userinfo = &DBUserInfo{}

	err = row.Scan(&userinfo.ID, &userinfo.Username, &userinfo.Email, &userinfo.PasswordHash, &userinfo.Elo, &userinfo.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			userinfo = nil
			return
		}
		userinfo = nil
		return
	}

	return
}
