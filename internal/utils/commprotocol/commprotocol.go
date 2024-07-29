package commprotocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

/**
* Binary Communication Protocol
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
	var buf bytes.Buffer

	// convert action to binary
	err := binary.Write(&buf, binary.BigEndian, uint8(action))

	if err != nil {
		return nil, err
	}

	// convert value based on action
	switch action {
	case JOIN:
		// new user joined, value will be user id
		if str, ok := any(value).(string); ok {

			// write length of string in buffer
			err := binary.Write(&buf, binary.BigEndian, uint8(len(str)))
			if err != nil {
				return nil, err
			}

			// write string at the end of buffer as value
			_, err = buf.WriteString(str)

			if err != nil {
				return nil, err
			}
		} else {
			// handle the case its not a string
			return nil, fmt.Errorf("expected value to be a string for action: %d", action)
		}
		break
	case LEAVE:
		log.Println("left")
		break
	default:
		break
	}

	return buf.Bytes(), nil
}

func Testing() {
	buf, err := EncodeMessage(JOIN, "JohnDoe")

	if err != nil {
		log.Println("Err:", err)
	}

	log.Printf("Encoded buffer: %x\n", buf)

}
