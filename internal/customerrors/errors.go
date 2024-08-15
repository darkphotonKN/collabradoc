package customerrors

import "errors"

// --- Users & Auth Errors ---

var PasswordIncorrectErr = errors.New("Password entered was incorrect.")
var UserExistsErr = errors.New("User already exists.")

// --- Live Session Errors ---
var LiveSessionUnauthorized = errors.New("This user or userId is unauthorized to use this live session.")
