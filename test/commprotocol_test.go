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

// Returning Editor List Encoding Test
func TestEditorListEncoding(t *testing.T) {
	usersList := []string{
		"username1",
		"username2",
		"mary poppins",
		"batman",
		"ahsoka",
		"anakin",
	}

	users := "username1,username2,mary poppins,batman,ahsoka,anakin"

	encodedMsg, err := commprotocol.EncodeMessage(commprotocol.EDITOR_LIST, usersList)

	if err != nil {
		t.Error("Error occured when attemping to encode message")
	}

	expected := []byte{
		0x14,        // Action
		0, 0, 0, 53, // Length of Users in 4 bytes, matching the uint32 used in the encoding function
	}

	expected = append(expected, []byte(users)...)

	if !bytes.Equal(encodedMsg, expected) {
		t.Errorf("EncodeMessage with %v, result %v; expected %v", users, encodedMsg, expected)
	}
}
