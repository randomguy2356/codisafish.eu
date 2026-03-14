package game

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"codisafish.eu/app/internal/auth"
)

type sseHandler struct {
	DB *sql.DB
}

type pingHandler struct {
	DB *sql.DB
}

type connectedUsers struct {
	Mutex sync.RWMutex
	Map   map[string]chan string
}

var ConnectedUsers = &connectedUsers{
	Map: make(map[string]chan string),
}

func (connected_users *connectedUsers) Register(userName string) chan string {
	channel := make(chan string)
	connected_users.Mutex.Lock()
	connected_users.Map[userName] = channel
	connected_users.Mutex.Unlock()
	return channel
}

func (connected_users *connectedUsers) Unregister(userName string) {
	connected_users.Mutex.Lock()
	if ch, ok := connected_users.Map[userName]; ok {
		close(ch)
		delete(connected_users.Map, userName)
	}
	connected_users.Mutex.Unlock()
}

func (connected_users *connectedUsers) Send(targetName, message string) bool {
	connected_users.Mutex.RLock()
	channel, ok := connected_users.Map[targetName]
	connected_users.Mutex.RUnlock()
	if ok {
		channel <- message
	}
	return ok
}

func RegisterSSE(mux *http.ServeMux, db *sql.DB) {
	ssehandler := sseHandler{DB: db}

	println("register sse")

	mux.Handle("/api/events", ssehandler)
}

func (sseHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	println("reached backend")
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sid_cookie, err := request.Cookie("sid")

	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(writer, "not logged in", http.StatusUnauthorized)
			return
		}
		http.Error(writer, "bad cookies", http.StatusBadRequest)
		return
	}

	session := auth.ValidateSID(sid_cookie.Value, writer)

	if session == nil {
		http.Error(writer, "not logged in", http.StatusUnauthorized)
		return
	}

	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")

	ch := ConnectedUsers.Register(session.User)
	defer ConnectedUsers.Unregister(session.User)

	flusher, ok := writer.(http.Flusher)

	if !ok {
		http.Error(writer, "SSE not supported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case msg, open := <-ch:
			if !open {
				return
			}
			fmt.Fprintf(writer, "data: %s\n\n", msg)
			flusher.Flush()
		case <-request.Context().Done():
			return
		}
	}
}

func RegisterPing(mux *http.ServeMux, db *sql.DB) {
	pinghandler := pingHandler{
		DB: db,
	}
	println("register ping")

	mux.Handle("/api/ping", pinghandler)
}

func (pingHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	println("reached backend")
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sid_cookie, err := request.Cookie("sid")

	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(writer, "not logged in", http.StatusUnauthorized)
			return
		}
		http.Error(writer, "bad cookies", http.StatusBadRequest)
		return
	}

	session := auth.ValidateSID(sid_cookie.Value, writer)

	senderName := session.User
	targetName := request.URL.Query().Get("target")

	online := ConnectedUsers.Send(targetName, senderName)

	if online {
		fmt.Fprintln(writer, "delivered")
	} else {
		http.Error(writer, "target not online or doesn't exist", http.StatusBadRequest)
	}
}
