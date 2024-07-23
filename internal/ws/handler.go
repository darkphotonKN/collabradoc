package ws

import (
	"bytes"
	"encoding/binary"
	"errors"
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
* Incoming / Outgoing message composes of two "enums" and a value
**/

type Action int

const (
	JOIN     Action = 0x12
	LEAVE    Action = 0x13
	USERLIST Action = 0x14
	EDIT     Action = 0x15
	DOCUMENT Action = 0x16
)

func EncodeMessage[T any](action Action, value T) ([]byte, error) {
	buf := new(bytes.Buffer)

	// convert action to binary
	err := binary.Write(buf, binary.BigEndian, uint8(action))

	if err != nil {
		log.Printf("Error occured when converting action %d to binary.\n", action)
	}

	// convert value based on action
	switch action {
	case JOIN:
		// new user joined, value will be user id
		log.Println("New User Joined Editor")

		// type assertion
		strVal, ok := value.(string)
		if !ok {
			return nil, errors.New("Action Join requires value of type string.")
		}
		// strBytes := []byte(value)

		break
	case LEAVE:
		log.Println("left")
		break
	default:
		break
	}

	return buf.Bytes()
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	// setup websocket

}

func Testing() {
	log.Println("Websockets Testing")

	// encodedMsg := EncodeMessage(LEAVE, 1213121391090)
	// secondEncodedMsg := EncodeMessage(JOIN, 1213121391090)

	// var message bytes.Buffer
	//
	// message.Write(encodedMsg)
	//
	// message.Write(secondEncodedMsg)
	//
	// log.Println("Encoded Message:", message)
}
