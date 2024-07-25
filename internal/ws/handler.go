package ws

import (
	"fmt"

	// "errors"
	"log"
	"net/http"

	"github.com/darkphotonKN/collabradoc/internal/utils/commprotocol"
	"github.com/gorilla/websocket"
)

// channel to track websocket payloads
var wsChan = make(chan WebSocketPayload)

// track connected clients
var clients = make(map[WebSocketConnection]string)

// response of payload sent back to clients via websocket
type WebSocketResponse[T any] struct {
	MessageType string `json:"message_type"`
	Action      string `json:"action"`
	Value       T      `json:"value"`
}

// Paylaod for sending back Websocket Information to client
type WebSocketPayload struct {
	Action string `json:"action"`
	Value  string `json:"value"`
	Conn   WebSocketConnection
}

type WebSocketConnection struct {
	*websocket.Conn
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

	response.MessageType = "ServerMessage"
	response.Value = "Game server connected."

	err = ws.WriteJSON(response)

	// client connected

	// create new connection type and add them to list of connections
	clientConnection := WebSocketConnection{
		Conn: ws,
	}
	clients[clientConnection] = ""

	// start goroutine thread to listen to all future incoming payloads
	go ListenForWS(&clientConnection)

	if err != nil {
		log.Println("Error when writing back to client:", err)
	}
}

func ListenForWS(conn *WebSocketConnection) {
	// logs error when the function stops and recovers
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error:", r)
		}
	}()

	log.Println("Listening for websocket connection. Current clients", clients)

	var payload WebSocketPayload

	// infinite loop to listen to incoming payloads

	for {
		err := conn.ReadJSON(&payload)

		if err != nil {
			// do nothing if there is an error
		} else {
			// handle new connection
			payload.Conn = *conn

			// send payload back to websocket channel
			wsChan <- payload
		}

	}
}

// Listen to the WebSocket CHANNEL
func ListenForWSChannel() {
	log.Println("Started listening concurrently for websocket connections on a goroutine.")
	var genericResponse WebSocketResponse[any]

	for {
		// storing websocket payload coming from wsChan
		event := <-wsChan

		// based on action we do something different

		log.Println("event received:", event)
		// TODO: Decode to figure out action and value

		// Match based on Action
		switch event.Action {
		case "client_list":
			log.Printf("client %s joining game.", event.Value)

			genericResponse.Action = "join_game"

			// add user to list of client connections
			clients[event.Conn] = event.Value // add user to map

			// get list of client for user
			genericResponse = getclientsListRes()
			broadcastToAll(genericResponse)

			// skip rest of function and continue listening for further websocket messages
			continue
		case "join_game":
			// add them to the connection

			clients[event.Conn] = event.Value

			// test encoding

			// encode message to binary
			encodedMsg, err := commprotocol.EncodeMessage(commprotocol.JOIN, event.Value)

			if err != nil {
				fmt.Errorf("Error occured when attempting to encode message %s, err was %s", event.Value, err)
			}

			// send binary message back to user
			event.Conn.WriteMessage(websocket.BinaryMessage, encodedMsg)

			continue
		default:
			return

		}
		// responds events sent to the channel to all users

		// not matching anything, we send back generic response
		genericResponse.Action = "event"
		genericResponse.Value = fmt.Sprintf("Message received, action was %s. Value was: %s", event.Action, event.Value)
		genericResponse.MessageType = "server_message"
		// broadcast to all clients
		broadcastToAll(genericResponse)
	}
}

// get the clients list and package it to fit action and message
func getclientsListRes() WebSocketResponse[any] {
	var clientsNameList []string

	// convert clients map into slice of names
	for _, name := range clients {
		log.Println("client", name)
		clientsNameList = append(clientsNameList, name)
	}

	clientListRes := WebSocketResponse[any]{
		Action:      "join_game",
		Value:       clientsNameList,
		MessageType: "client_list",
	}

	return clientListRes
}

// Broadcast to all users
func broadcastToAll(message WebSocketResponse[any]) {
	// loop through all connected clients and broadcast to them
	for clientWS := range clients {

		err := clientWS.WriteJSON(message)

		// handle if client errored / disconnected
		if err != nil {
			log.Println("Websocket errored")

			// close their WS connection
			_ = clientWS.Close()

			// remove the client that errored
			delete(clients, clientWS)
		}
	}
}
