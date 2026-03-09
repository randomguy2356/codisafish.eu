package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/alexedwards/argon2id"
)

type Session struct {
	ExpiresAtAbs time.Time
	ExpiresAt    time.Time
	User         string
}

type MutexSessions struct {
	Mutex sync.RWMutex
	Map   map[string]Session
}

var Sessions = &MutexSessions{
	Map: make(map[string]Session),
}

func CreateSession(username string, writer http.ResponseWriter) (sid string) {
	session := Session{
		ExpiresAtAbs: time.Now().Add(24 * time.Hour),
		ExpiresAt:    time.Now().Add(time.Hour),
		User:         username,
	}
	sid = NewUSID()
	Sessions.Mutex.Lock()
	Sessions.Map[sid] = session
	Sessions.Mutex.Unlock()
	SetSIDCookie(writer, sid, time.Hour)
	return
}

func NewUSID() string {
	Sessions.Mutex.RLock()
	defer Sessions.Mutex.RUnlock()
	var sid string
	for {
		sid = rand.Text()
		if _, exists := Sessions.Map[sid]; !exists {
			break
		}
	}
	return sid
}

func ValidateUser(username string, password string, ctx context.Context, db *sql.DB) (bool, error) {

	query := "SELECT password_hash FROM users WHERE username = ?"

	var password_hash string

	row := db.QueryRowContext(ctx, query, username)

	err := row.Scan(&password_hash)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	mathes, err := argon2id.ComparePasswordAndHash(password, password_hash)
	if err != nil {
		return false, err
	}

	return mathes, nil
}

func ValidateSID(sid string, writer http.ResponseWriter) *Session {
	Sessions.Mutex.Lock()
	defer Sessions.Mutex.Unlock()

	session, exists := Sessions.Map[sid]

	if !exists {
		SetSIDCookie(writer, "", -1)
		return nil
	}

	if session.ExpiresAt.Before(time.Now()) {
		delete(Sessions.Map, sid)
		SetSIDCookie(writer, "", -1)
		return nil
	}

	session.ExpiresAt = time.Now().Add(time.Hour)

	Sessions.Map[sid] = session

	return &session
}

func RotateSid(sid string, writer http.ResponseWriter) *string {
	Sessions.Mutex.Lock()
	defer Sessions.Mutex.Unlock()
	session, exists := Sessions.Map[sid]

	if !exists {
		return nil
	}

	sid = NewUSID()
	Sessions.Map[sid] = session

	return &sid
}

func InvalidateSID(sid string, writer http.ResponseWriter) {
	Sessions.Mutex.Lock()
	defer Sessions.Mutex.Unlock()

	delete(Sessions.Map, sid)
	SetSIDCookie(writer, "", -1)
}

func SetSIDCookie(writer http.ResponseWriter, sid string, maxAge time.Duration) {
	http.SetCookie(writer, &http.Cookie{
		Name:     "sid",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(maxAge),
	})
}

func CheckSIDCookie(writer http.ResponseWriter, request *http.Request) (bool, error) {
	sid, err := request.Cookie("sid")
	if err != http.ErrNoCookie {
		if err != nil {
			return false, err
		}
		session := ValidateSID(sid.Value, writer)
		if session != nil {
			writer.WriteHeader(http.StatusOK)
			return true, nil
		}
	}
	return false, nil
}
