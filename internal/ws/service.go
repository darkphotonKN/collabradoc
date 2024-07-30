package ws

import (
	"fmt"
	"log"

	"github.com/darkphotonKN/collabradoc/internal/types"
	"github.com/darkphotonKN/collabradoc/internal/utils/commprotocol"
	"github.com/gorilla/websocket"
)

func ListenForWS(conn *types.WebSocketConnection) {
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
			log.Printf("Error occured when reading payload from websocket %s", err)
		} else {

			// handle new connection

			// create websocket information packet for handling each unique connection
			wsInfo := WebSocketInfo{
				WebSocketPayload: WebSocketPayload{
					Action: payload.Action,
					Value:  payload.Value,
				},
				Conn: *conn,
			}

			log.Printf("payload: %+v", payload)

			// send payload back to websocket channel
			wsChan <- wsInfo
		}
	}
}

// Listen to the WebSocket CHANNEL
func ListenForWSChannel() {
	log.Println("Started listening concurrently for websocket connections on a goroutine")
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

			// encode message to binary
			encodedMsg, err := commprotocol.EncodeMessage(commprotocol.JOIN, event.Value)
			log.Println("encodedMessage:", encodedMsg)

			if err != nil {
				log.Printf("Error occured when attempting to encode message %s, err was %s", event.Value, err)
			}

			// send binary message back to user
			event.Conn.WriteMessage(websocket.BinaryMessage, encodedMsg)

			continue

		default:
			// responds events sent to the channel to all users

			// not matching anything, we send back generic response
			genericResponse.Value = fmt.Sprintf("Message received, action was %s. Value was: %s", event.Action, event.Value)
			genericResponse.Action = "server_message"
			// broadcast to all clients
			broadcastToAll(genericResponse)
			continue
		}
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
		Action: "client_list",
		Value:  clientsNameList,
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
