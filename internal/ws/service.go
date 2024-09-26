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

type WebSocketService struct {
	// channels for transporting payloads
	wsChan          chan WebSocketInfo
	wsCommunityChan chan WebSocketCommunityInfo

	// maps to hold connected websocket user instances
	clientConnections    map[string]map[types.WebSocketConnection]string
	communityClientConns map[uint]map[types.WebSocketConnection]string

	mu sync.Mutex
}

func NewWebSocketService() *WebSocketService {
	// channels to track websocket payloads
	var wsChan = make(chan WebSocketInfo)                   // private channel
	var wsCommunityChan = make(chan WebSocketCommunityInfo) // public channel

	// map of sessionIds that map to maps of websocket connections to client names
	var clientConnections = make(map[string]map[types.WebSocketConnection]string)

	// map of documentId that map to maps of websocket connections to client names
	var communityClientConns = make(map[uint]map[types.WebSocketConnection]string)

	return &WebSocketService{
		wsChan:               wsChan,
		wsCommunityChan:      wsCommunityChan,
		clientConnections:    clientConnections,
		communityClientConns: communityClientConns,
	}
}

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
* -- ListenForWS Private Connection Version --
* Handles each unique client via websocket
*
* This function is ran concurrently for each unique client that connects
* to a live session (to edit documents).
**/
func (w *WebSocketService) ListenForWS(conn *types.WebSocketConnection, sessionId string) {
	defer func() {
		conn.Close()
		fmt.Println("Connection closed.")
	}()

	log.Println("Listening for websocket connection. All session clients", w.clientConnections)

	var payload WebSocketPayload

	for {
		// TODO: Update to decode via comm protocol
		err := conn.ReadJSON(&payload)

		if err != nil {
			// handle unexpected client errors
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Unexpected Close Error: %v\n", err)

				// remove client from connection
				delete(w.clientConnections[sessionId], *conn)
			} else {
				fmt.Printf("JSON Error: %v\n", err)

				// remove client from connection
				delete(w.clientConnections[sessionId], *conn)
			}

			break // only exits the loop, not entire function, allows for graceful exit

		} else {
			// no error - handle new connection

			// handling payload each subsequent request
			// create websocket information packet for handling each unique connectionws service
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
			w.wsChan <- wsInfo
		}
	}
}

/**
* Listens to the WebSocket Channel
*
* This function is running concurrently at the start of the application.
* Will handle all messages sent to the central channel and handle the
* messages accordingly.
**/
func (w *WebSocketService) ListenForWSChannel() {
	log.Println("Started listening concurrently for websocket connections on a goroutine")

	for {
		// handle different types of channels
		select {

		// handle incoming websocket events
		case event := <-w.wsChan:

			// storing websocket payload coming from wsChan

			log.Println("event value received:", event.Value)
			// TODO: Decode to figure out action and value
			// Match based on Action
			switch event.Action {
			case "editor_list":

				// get list of client for user
				list, err := w.getEditorList(event.SessionId)

				if err != nil {
					encodedErrMsg, _ := commprotocol.EncodeMessage(commprotocol.SYSTEM_MSG, fmt.Sprintf("Error when retrieving list of users: %v", err))

					event.Conn.WriteMessage(websocket.BinaryMessage, encodedErrMsg)
				}

				w.broadcastToAllClients(list, event.SessionId)

				// skip rest of function and continue listening for further websocket messages
				continue

			case "join_doc":
				fmt.Printf("ADDING CLIENT %s and for SESSIONID %s\n", event.Value, event.SessionId)

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

				existClientConnections, ok := w.clientConnections[event.SessionId]

				// initialize if map is empty, prevent nil pointer exceptions
				if !ok {
					existClientConnections = make(map[types.WebSocketConnection]string)
					w.clientConnections[event.SessionId] = existClientConnections
				}

				w.clientConnections[event.SessionId][event.Conn] = user.Name

				// encode message to binary
				encodedMsg, err := commprotocol.EncodeMessage(commprotocol.JOIN, user.Name)
				log.Println("encodedMessage:", encodedMsg)

				if err != nil {
					log.Printf("Error occured when attempting to encode message %s, err was %s", event.Value, err)
				}

				// get current editor client list in binary
				encodedEditors, err := w.getEditorList(event.SessionId)
				fmt.Printf("encoded editor list %v\n", encodedEditors)

				if err != nil {
					encodedErrMsg, _ := commprotocol.EncodeMessage(commprotocol.SYSTEM_MSG, fmt.Sprintf("Error when retrieving list of users: %v", err))

					event.Conn.WriteMessage(websocket.BinaryMessage, encodedErrMsg)
				}

				// send binary message to all users saying who has joined
				w.broadcastToAllClients(encodedMsg, event.SessionId)

				// send list of current editors back to all clients
				w.broadcastToAllClients(encodedEditors, event.SessionId)

				continue

			// when user disconnects
			case "disconnected":
				disconnectedUser := w.clientConnections[event.SessionId][event.Conn]

				fmt.Printf("User %s disconnected\n", disconnectedUser)

				// close the channel
				event.Conn.Close()

				// new delete that user
				delete(w.clientConnections[event.SessionId], event.Conn)

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
func (w *WebSocketService) getEditorList(sessionId string) ([]byte, error) {

	// retreive corresponding session map for users
	sessionUsers := w.clientConnections[sessionId]
	sessionUsernames := make([]string, len(sessionUsers))

	// convert clientConnections' current session connections (sessionUsers) to a slice of strings (names)
	for _, name := range sessionUsers {
		sessionUsernames = append(sessionUsernames, name)
	}

	// encode slice of usernames
	encodedSessionUsernames, err := commprotocol.EncodeMessage(commprotocol.EDITOR_LIST, sessionUsernames)

	if err != nil {
		fmt.Printf("Error when attempting to encode %v, error was %v", sessionUsernames, err)
		return nil, err
	}

	// return encoded
	return encodedSessionUsernames, nil
}

func (w *WebSocketService) broadcastToAllClients(encodedMsg []byte, sessionId string) {
	sessionClients := w.clientConnections[sessionId]
	for wsConn := range sessionClients {
		err := wsConn.Conn.WriteMessage(websocket.BinaryMessage, encodedMsg)

		if err != nil {
			// close connection if that write failed
			wsConn.Conn.Close()

			// delete from client list
			delete(w.clientConnections[sessionId], wsConn)
		}
	}
}

/**
* -- ListenForWS Community Public Connection Version --
* Handles each unique client via websocket
*
* This function is ran concurrently for each unique client that connects
* to a public community document.
**/

func (w *WebSocketService) ListenForWSCommunity(conn *types.WebSocketConnection, documentId uint) {
	defer func() {
		conn.Close()
		fmt.Println("Connection closed.")
	}()

	log.Println("Listening for websocket connection. All session clients", w.clientConnections)

	var payload WebSocketCommunityInfo

	for {
		err := conn.ReadJSON(&payload)

		if err != nil {
			// handle unexpected client errors
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Unexpected Close Error: %v\n", err)

				// unexepted connection errors, delete user from the community connection pool
				delete(w.communityClientConns[documentId], *conn)
			} else {
				fmt.Printf("JSON Error: %v\n", err)

				// remove client from connection
				delete(w.communityClientConns[documentId], *conn)

				break // only exits the loop, not entire function, allows for graceful exit

			}

		} else {
			// handling payload each subsequent request
			// create websocket information packet for handling each unique connectionws service
			message := WebSocketCommunityInfo{
				WebSocketPayload: WebSocketPayload{
					Action: payload.Action,
					Value:  payload.Value,
				},
				DocumentID: documentId,
				Conn:       *conn,
			}

			log.Printf("payload: %+v", payload)

			w.wsCommunityChan <- message
		}
	}
}
