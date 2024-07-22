package customerrors

import "errors"

// --- Users & Auth Errors ---

var PasswordIncorrectErr = errors.New("Password entered was incorrect.")
var UserExistsErr = errors.New("User already exists.")
