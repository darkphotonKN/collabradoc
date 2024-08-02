package commprotocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
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
	JOIN        Action = 0x12
	LEAVE       Action = 0x13
	EDITOR_LIST Action = 0x14
	EDIT        Action = 0x15
	SYSTEM_MSG  Action = 0x16
)

func EncodeMessage[T any](action Action, value T) ([]byte, error) {
	var buf bytes.Buffer

	// [action]: convert action to binary
	err := binary.Write(&buf, binary.BigEndian, uint8(action))

	if err != nil {
		return nil, err
	}

	// convert value based on action
	switch action {

	// user joins or leaves document to collaborate in editing
	// encoding structure: [action] - [string length] - [username string]
	case JOIN, LEAVE, SYSTEM_MSG:
		// will be user id so will be type string
		if str, ok := any(value).(string); ok {

			// write length of string in buffer for decode function
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

	// provides current users that can edit the document
	// encoding structure: [action] - [byte length] - [users []byte separated by ","]
	case EDITOR_LIST:
		// so value will be a slice of current users
		if users, ok := any(value).([]string); ok {

			fmt.Println("joining:", users)
			userString := strings.Join(users, ",")
			fmt.Println("userString:", userString)

			usersBytes := make([][]byte, len(users))

			for i, user := range users {
				// convert user to byte and store it
				usersBytes[i] = []byte(user)
			}

			usersBytesJoined := bytes.Join(usersBytes, []byte(","))

			// [byte length]: write the length to the next spot after action
			length := uint32(len(usersBytesJoined))
			binary.Write(&buf, binary.BigEndian, length)

			// [users []byte separated by ","]: write the users in bytes with delimiter
			buf.Write(usersBytesJoined)
		} else {
			// handle the case its not a map of the correct type
			return nil, fmt.Errorf("expected value to be a slice of type \"string\"  for action: %d", action)
		}

		break
	default:
		break
	}

	return buf.Bytes(), nil
}
