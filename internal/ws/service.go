package ws

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/darkphotonKN/collabradoc/internal/types"
	"github.com/darkphotonKN/collabradoc/internal/user"
	"github.com/darkphotonKN/collabradoc/internal/utils/commprotocol"
	"github.com/gorilla/websocket"
)

var (
	mu       sync.Mutex
	shutdown = make(chan struct{})
)

func Shutdown() {
	close(shutdown)
	// Close all client connections
	mu.Lock()
	defer mu.Unlock()
	for conn := range clients {
		conn.Conn.Close()
		delete(clients, conn)
	}
}

func ListenForWS(conn *types.WebSocketConnection) {
	// logs error when the function stops and recovers
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Println("Recovered, error was:", r)
	// 	}
	// }()

	defer func() {
		conn.Close()
		fmt.Println("Connection closed.")
	}()

	log.Println("Listening for websocket connection. Current clients", clients)

	var payload WebSocketPayload

	// infinite loop to listen to incoming payloads

	for {
		err := conn.ReadJSON(&payload)

		// Handle errors
		if err != nil {
			// handle unexpected client errors
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Unexpected Close Error: %v\n", err)

				// remove client from connection
				delete(clients, *conn)
			} else {
				// json error
				fmt.Printf("JSON Error: %v\n", err)

				// remove client from connection
				delete(clients, *conn)
			}

			break // only exits the loop, not entire function, allows for graceful exit

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

	for {

		// handle different types of channels
		select {

		// handle incoming websocket events
		case event := <-wsChan:

			// storing websocket payload coming from wsChan

			log.Println("event value received:", event.Value)
			// TODO: Decode to figure out action and value
			// Match based on Action
			switch event.Action {
			case "editor_list":

				// get list of client for user
				list, err := getEditorList()

				if err != nil {
					encodedErrMsg, _ := commprotocol.EncodeMessage(commprotocol.SYSTEM_MSG, fmt.Sprintf("Error when retrieving list of users: %v", err))

					event.Conn.WriteMessage(websocket.BinaryMessage, encodedErrMsg)
				}

				broadcastToAllClients(list)

				// skip rest of function and continue listening for further websocket messages
				continue

			case "join_doc":
				fmt.Printf("ADDING CLIENT with CONNECTION %v and NAME %s\n", event.Conn, event.Value)

				// get user from db

				userId, err := strconv.ParseUint(event.Value, 10, 0)
				if err != nil {
					fmt.Printf("Error when attempting to parse uint from user id:\n", userId)
				}

				user, err := user.FindUserById(uint(userId))

				if err != nil {
					encodedErrMsg, _ := commprotocol.EncodeMessage(commprotocol.SYSTEM_MSG, fmt.Sprintf("Error when retrieving user with connected id: %v", err))
					event.Conn.WriteMessage(websocket.BinaryMessage, encodedErrMsg)
					continue // continue to next message read
				}

				// add them to the current pool of live doc editors
				clients[event.Conn] = user.Name

				// encode message to binary
				encodedMsg, err := commprotocol.EncodeMessage(commprotocol.JOIN, event.Value)
				log.Println("encodedMessage:", encodedMsg)

				if err != nil {
					log.Printf("Error occured when attempting to encode message %s, err was %s", event.Value, err)
				}

				// get current editor client list in binary
				encodedEditors, err := getEditorList()
				fmt.Printf("encoded editor list %v\n", encodedEditors)

				if err != nil {
					encodedErrMsg, _ := commprotocol.EncodeMessage(commprotocol.SYSTEM_MSG, fmt.Sprintf("Error when retrieving list of users: %v", err))

					event.Conn.WriteMessage(websocket.BinaryMessage, encodedErrMsg)
				}

				// send binary message to all users saying who has joined
				broadcastToAllClients(encodedMsg)

				// send list of current editors back to all clients
				broadcastToAllClients(encodedEditors)

				continue

			// when user disconnects
			case "disconnected":
				disconnectedUser := clients[event.Conn]

				fmt.Printf("User %s disconnected\n", disconnectedUser)

				// close the channel
				event.Conn.Close()

				// delete that user
				delete(clients, event.Conn)

			default:
				// not matching anything, we send back generic response

				continue
			}

		case <-shutdown:
			log.Println("Stopped listening to websocket channel.")
		}
	}
}

// get the clients list and package it to fit action and message
func getEditorList() ([]byte, error) {
	editorUsernames := make([]string, 0, len(clients)) // pre-allocate the size but initialize at 0

	fmt.Printf("editorUsernames length: %d\n", len(editorUsernames))
	fmt.Printf("clients before decoding: %v\n", clients)

	// convert clients map to a slice of strings (usernames)
	for _, username := range clients { // forgo key which is the WS connection
		fmt.Printf("debug username in client list was %s \n", username)
		editorUsernames = append(editorUsernames, username)
	}

	fmt.Printf("editorUsernames after: %v\n", editorUsernames)
	fmt.Printf("editorUsernames length after: %v\n", len(editorUsernames))

	// encode slice of usernames
	encodedEditorUsernames, err := commprotocol.EncodeMessage(commprotocol.EDITOR_LIST, editorUsernames)

	if err != nil {
		fmt.Printf("Error when attempting to encode %v, error was %v", editorUsernames, err)
		return nil, err
	}

	// return encoded
	return encodedEditorUsernames, nil
}

func broadcastToAllClients(encodedMsg []byte) {
	for wsConn := range clients {
		err := wsConn.Conn.WriteMessage(websocket.BinaryMessage, encodedMsg)

		if err != nil {
			// close connection if that write failed
			wsConn.Conn.Close()

			// delete from client list
			delete(clients, wsConn)
		}
	}
}
