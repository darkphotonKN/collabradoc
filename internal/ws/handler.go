package ws

import (
	"log"
	"net/http"

	"github.com/darkphotonKN/collabradoc/internal/types"
	"github.com/gorilla/websocket"
)

// channel to track websocket payloads
var wsChan = make(chan WebSocketInfo)

// track connected clients
var clients = make(map[types.WebSocketConnection]string)

// response of payload sent back to clients via websocket
type WebSocketResponse[T any] struct {
	Action string `json:"action"`
	Value  T      `json:"value"`
}

// Payload for sending / recieving Websocket Information
type WebSocketPayload struct {
	Action string `json:"action"`
	Value  string `json:"value"`
}

// For internal websocket tracking
type WebSocketInfo struct {
	WebSocketPayload
	Conn types.WebSocketConnection
}

// for upgrading response writer / request connections
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *(http.Request)) bool { return true },
}

// update strandard response writer, request and header to a websocket connection
func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error creating websocket connection:", err)
	}

	log.Println("Connected to websocket server.")

	var response WebSocketResponse[string]

	response.Action = "ServerMessage"
	response.Value = "Game server connected."

	err = ws.WriteJSON(response)

	// client connected

	// create new connection type and add them to list of connections
	clientConnection := types.WebSocketConnection{
		Conn: ws,
	}

	// initialize connection for client but don't allocate username until acquired
	// clients[clientConnection] = ""

	// start goroutine thread to listen to all future incoming payloads
	go ListenForWS(&clientConnection)

	if err != nil {
		log.Println("Error when writing back to client:", err)
	}
}
