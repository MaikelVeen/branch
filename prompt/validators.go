package prompt

import "net/mail"

// ValidateFunc is for any validation functions that validates a given input. It should return
// an err if the input is not valid.
type ValidateFunc func(string) error

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
