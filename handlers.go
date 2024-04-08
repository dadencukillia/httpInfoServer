package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// -- WebSocket --

// When minecraft agent wants to connect
func WsClientConnect(w http.ResponseWriter, r *http.Request) {
	con, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		con.Close()
		return
	}

	// Sending message with data collection command
	err = con.WriteMessage(websocket.TextMessage, []byte("collectData"))
	if err != nil {
		con.Close()
		return
	}

	// Receiving collected data
	_, content, err := con.ReadMessage()
	if err != nil {
		apiLog.Println(err)
		con.Close()
		return
	}

	// Turning JSON data to map
	var a map[string]any

	err = json.Unmarshal(content, &a)
	if err != nil {
		apiLog.Println(err)
		con.Close()
		return
	}

	// Getting player's nickname from json
	nickname, ok := a["nickname"]
	if !ok {
		con.Close()
		return
	}

	if s, ok := nickname.(string); ok {
		wsLog.Println(s + " connected")

		// Removes client from list when disconnected from
		con.SetCloseHandler(func(code int, text string) error {
			wsLog.Println(s + " disconnected")
			delete(clients, s)
			return nil
		})

		// Add websocket client to list
		clients[s] = &Connection{
			ws:    con,
			mutex: sync.Mutex{},
		}
		return
	}
	con.Close()
}

// -- API HTTP Server --

// Logging API HTTP server activity
func apiLogMiddleware(r *http.Request) {
	apiLog.Printf("(%s): %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
}

// When client wants to get nicknames list of connected players
func ApiGetPlayersList(w http.ResponseWriter, r *http.Request) {
	apiLogMiddleware(r)

	// Nicknames list of connected players = keys of map connected players
	keys := []string{}
	for k := range clients {
		keys = append(keys, k)
	}

	// Turning list in JSON string
	jsonCont, err := json.Marshal(keys)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("JSON generating error"))
		return
	}

	// Sending
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonCont)
}

// When client wants to get info about player by him nickname
func ApiGetPlayerInfo(w http.ResponseWriter, r *http.Request) {
	apiLogMiddleware(r)

	// Receiving nickname from path
	nickname := r.PathValue("nickname")

	// Getting WebSocket client instance from the map
	con, ok := clients[nickname]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found connected user with this nickname"))
		return
	}

	// Locking and unlocking mutex to not allow use WebSocket client twice at one time
	con.mutex.Lock()
	defer con.mutex.Unlock()

	// Sending data collection request
	err := con.ws.WriteMessage(websocket.TextMessage, []byte("collectData"))
	if err != nil {
		delete(clients, nickname)
		con.ws.Close()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("User with this nickname disconnected"))
		return
	}

	// Receiving collected data
	_, data, err := con.ws.ReadMessage()
	if err != nil {
		delete(clients, nickname)
		con.ws.Close()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("User with this nickname disconnected"))
		return
	}

	// Checking if received data is correct
	var a interface{}
	err = json.Unmarshal(data, &a)
	if err != nil {
		delete(clients, nickname)
		con.ws.Close()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Incorrect data received from user"))
		return
	}

	// Sending data about player to client
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

// When page not found
func ApiNotFound(w http.ResponseWriter, r *http.Request) {
	apiLogMiddleware(r)

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("This page is not exists"))
}
