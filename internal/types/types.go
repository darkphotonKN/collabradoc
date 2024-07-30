package types

import "github.com/gorilla/websocket"

/* -- Types --
* This file stores all the shared types, mainly for DRY code and preventing cycle imports.
 */

// -- WebSockets --

type WebSocketConnection struct {
	*websocket.Conn
}
