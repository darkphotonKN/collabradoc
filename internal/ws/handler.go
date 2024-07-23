package ws

import (
	"log"
	"net/http"
)

/**
* Binary Protocol
*
* Basic Structure
* ---
* message type - action type - value
*
* So it composes of two enumes and a value
**/

func WsHandler(w http.ResponseWriter, r *http.Request) {
	// setup websocket

}

func Testing() {
	log.Println("Websockets Testing")

}
