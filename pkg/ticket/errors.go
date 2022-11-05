package ticket

import "errors"

var ErrTypeAssertionAuth = errors.New("Could not parse credential data")
var ErrNotUnauthorized = errors.New("Could not authorize with ticket system, are your credentials still valid?")
var ErrNotFound = errors.New("Could not find issue/ticket")
var ErrCredentialSaving = errors.New("Could not save credentials to system keyring")
