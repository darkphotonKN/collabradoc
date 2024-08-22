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
	// TODO: Close all client connections
	// mu.Lock()
	// defer mu.Unlock()
	// for conn := range clients {
	// 	conn.Conn.Close()
	// 	delete(clients, conn)
	// }
}

/**
* Handles each unique client via websocket
*
* This function is ran concurrently for each unique client that connects
* to a live session (to edit documents).
**/
func ListenForWS(conn *types.WebSocketConnection, sessionId string) {

	defer func() {
		conn.Close()
		fmt.Println("Connection closed.")
	}()

	log.Println("Listening for websocket connection. All session clients", clientConnections)

	var payload WebSocketPayload

	for {
		// TODO: Update to decode via comm protocol
		err := conn.ReadJSON(&payload)

		if err != nil {
			// handle unexpected client errors
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Unexpected Close Error: %v\n", err)

				// remove client from connection
				delete(clientConnections[sessionId], *conn)
			} else {
				fmt.Printf("JSON Error: %v\n", err)

				// remove client from connection
				delete(clientConnections[sessionId], *conn)
			}

			break // only exits the loop, not entire function, allows for graceful exit

		} else {

			// no error - handle new connection

			// handling payload each subsequent request
			// create websocket information packet for handling each unique connection
			wsInfo := WebSocketInfo{
				WebSocketPayload: WebSocketPayload{
					Action: payload.Action,
					Value:  payload.Value,
				},
				SessionId: sessionId,
				Conn:      *conn,
			}

			log.Printf("payload: %+v", payload)

			// send payload back to websocket channel
			wsChan <- wsInfo
		}
	}
}

/**
* Listens to the WebSocket CHANNEL
*
* This function is running concurrently at the start of the application.
* Will handle all messages sent to the central channel and handle the
* messages accordingly.
**/
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
				list, err := getEditorList(event.SessionId)

				if err != nil {
					encodedErrMsg, _ := commprotocol.EncodeMessage(commprotocol.SYSTEM_MSG, fmt.Sprintf("Error when retrieving list of users: %v", err))

					event.Conn.WriteMessage(websocket.BinaryMessage, encodedErrMsg)
				}

				broadcastToAllClients(list, event.SessionId)

				// skip rest of function and continue listening for further websocket messages
				continue

			case "join_doc":
				fmt.Printf("ADDING CLIENT with CONNECTION %v and NAME %s and for SESSIONID %s\n", event.Conn, event.Value, event.SessionId)

				// get user from db to store in current connection map
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

				// add them to their respective live sessions under their own connections
				temp := make(map[types.WebSocketConnection]string)
				temp[event.Conn] = user.Name
				clientConnections[event.SessionId] = temp

				// encode message to binary
				encodedMsg, err := commprotocol.EncodeMessage(commprotocol.JOIN, user.Name)
				log.Println("encodedMessage:", encodedMsg)

				if err != nil {
					log.Printf("Error occured when attempting to encode message %s, err was %s", event.Value, err)
				}

				// get current editor client list in binary
				encodedEditors, err := getEditorList(event.SessionId)
				fmt.Printf("encoded editor list %v\n", encodedEditors)

				if err != nil {
					encodedErrMsg, _ := commprotocol.EncodeMessage(commprotocol.SYSTEM_MSG, fmt.Sprintf("Error when retrieving list of users: %v", err))

					event.Conn.WriteMessage(websocket.BinaryMessage, encodedErrMsg)
				}

				// send binary message to all users saying who has joined
				broadcastToAllClients(encodedMsg, event.SessionId)

				// send list of current editors back to all clients
				broadcastToAllClients(encodedEditors, event.SessionId)

				continue

			// when user disconnects
			case "disconnected":
				disconnectedUser := clientConnections[event.SessionId][event.Conn]

				fmt.Printf("User %s disconnected\n", disconnectedUser)

				// close the channel
				event.Conn.Close()

				// new delete that user
				delete(clientConnections[event.SessionId], event.Conn)

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
func getEditorList(sessionId string) ([]byte, error) {

	// retreive corresponding session map for users
	sessionUsers := clientConnections[sessionId]
	sessionUsernames := make([]string, len(sessionUsers))

	// convert clientConnections' current session connections (sessionUsers) to a slice of strings (names)
	for _, name := range sessionUsers {
		sessionUsernames = append(sessionUsernames, name)
	}

	fmt.Printf("NEW sessionUsernames %v\n", sessionUsernames)
	fmt.Printf("NEW sessionUsernames length %d\n", len(sessionUsernames))

	// encode slice of usernames
	encodedSessionUsernames, err := commprotocol.EncodeMessage(commprotocol.EDITOR_LIST, sessionUsernames)

	if err != nil {
		fmt.Printf("Error when attempting to encode %v, error was %v", sessionUsernames, err)
		return nil, err
	}

	// return encoded
	return encodedSessionUsernames, nil
}

func broadcastToAllClients(encodedMsg []byte, sessionId string) {
	sessionClients := clientConnections[sessionId]
	for wsConn := range sessionClients {
		err := wsConn.Conn.WriteMessage(websocket.BinaryMessage, encodedMsg)

		if err != nil {
			// close connection if that write failed
			wsConn.Conn.Close()

			// delete from client list
			delete(clientConnections[sessionId], wsConn)
		}
	}
}
