package test

import (
	"bytes"
	"testing"

	"github.com/darkphotonKN/collabradoc/internal/utils/commprotocol"
)

// User Joining Encoding Test
func TestJoinEncoding(t *testing.T) {

	username := "John Doe"

	encodedMsg, err := commprotocol.EncodeMessage(commprotocol.JOIN, username)

	if err != nil {
		t.Error("Error occured when attemping to encode message")
	}

	expected := []byte{
		0x12,                                   // Action
		8,                                      // Length of Username
		'J', 'o', 'h', 'n', ' ', 'D', 'o', 'e', // User name
	}

	if !bytes.Equal(encodedMsg, expected) {
		t.Errorf("EncodeMessage with %v, result %v; expected %v", username, encodedMsg, expected)
	}
}

// User Joining Encoding Test
func TestEditorListEncoding(t *testing.T) {

	users := "username1,username2,mary poppins,batman,asohka,anakin"

	encodedMsg, err := commprotocol.EncodeMessage(commprotocol.EDITORLIST, users)

	if err != nil {
		t.Error("Error occured when attemping to encode message")
	}

	expected := []byte{
		0x14, // Action
		6,    // Length of Users
	}

	expected = append(expected, []byte(users)...)

	if !bytes.Equal(encodedMsg, expected) {
		t.Errorf("EncodeMessage with %v, result %v; expected %v", users, encodedMsg, expected)
	}
}
