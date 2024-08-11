package ws

import (
	// "context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/darkphotonKN/collabradoc/internal/types"
	"github.com/darkphotonKN/collabradoc/internal/utils/auth"
	"github.com/golang-jwt/jwt/v5"
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
	envKey := os.Getenv("JWT_SECRET_KEY")
	jwtKey := []byte(envKey)

	// extract the token from the query parameters
	tokenString := r.URL.Query().Get("token")

	fmt.Println("tokenString", tokenString)

	if tokenString == "" {
		http.Error(w, "No token in connection.", http.StatusUnauthorized)
		return
	}

	// authenticate jwt token
	claims := &auth.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// jwt unauthorized
	if err != nil {
		fmt.Printf("Jwt unauthorized, error: %s\n", err.Error())
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Invalid token signature", http.StatusUnauthorized)
			return
		}
		http.Error(w, fmt.Sprintf("Error when parsing token: %s", err), http.StatusBadRequest)
		return
	}

	if !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	ws, err := upgradeConnection.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error creating websocket connection:", err)
	}

	log.Printf("User connected to websocket server: %d \n", claims.UserID)

	// send client with join_user action to websocket channel to add user
	// to current list of maintained editors
	joinUserAction := WebSocketInfo{
		WebSocketPayload: WebSocketPayload{
			Action: "join_doc",
			Value:  strconv.FormatUint(uint64(claims.UserID), 10),
		},
		Conn: types.WebSocketConnection{
			Conn: ws,
		},
	}
	wsChan <- joinUserAction

	var response WebSocketResponse[string]

	response.Action = "ServerMessage"
	response.Value = "Editor server connected."

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
