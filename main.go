package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ws    *websocket.Conn
	mutex sync.Mutex // Allows use websocket client only one thread at one time
}

var (
	// gorilla/websocket configuration
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize: 4096 * 4,
	}

	// List of connected WebSocket clients (player's nickname = key)
	clients = map[string]*Connection{}

	// Loggers
	wsLog  = log.New(os.Stdout, "[WebSocket] ", log.Ltime|log.Lmsgprefix)
	apiLog = log.New(os.Stdout, "[API] ", log.Ltime|log.Lmsgprefix)
)

// Returns WebSocket Server routers
func websocketServerRouters() *http.ServeMux {
	ws := http.NewServeMux()

	// Registering handler functions
	ws.HandleFunc("/connectClient", WsClientConnect)

	return ws
}

// Returns HTTP API server routers
func apiServerRouters() *http.ServeMux {
	api := http.NewServeMux()

	// Registering handler functions
	api.HandleFunc("GET /users", ApiGetPlayersList)
	api.HandleFunc("GET /user/{nickname}", ApiGetPlayerInfo)
	api.HandleFunc("/", ApiNotFound)

	return api
}

// Entrance function
func main() {
	// Starting WebSocket server in another thread
	go func() {
		wsLog.Println("WebSocket server started!")
		wsLog.Fatalln(http.ListenAndServe("127.0.0.1:9922", websocketServerRouters()))
	}()

	// Parsing port from command-line argument
	port := 8000
	if len(os.Args) > 1 {
		i, err := strconv.Atoi(os.Args[1])
		if err == nil {
			if i > 21 && i <= 65535 {
				port = i
			}
		}
	}

	// Starting API server
	apiLog.Println("HTTP Server started on " + fmt.Sprint(port) + " port (URL: http://127.0.0.1:" + fmt.Sprint(port) + ")")
	apiLog.Fatalln(http.ListenAndServe(":"+fmt.Sprint(port), apiServerRouters()))
}
