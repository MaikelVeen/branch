package ticket

import "errors"

var ErrTypeAssertionAuth = errors.New("could not parse credential data")
var ErrNotUnauthorized = errors.New("could not authorize with ticket system, are your credentials still valid?")
var ErrNotFound = errors.New("could not find issue/ticket")
var ErrCredentialSaving = errors.New("could not save credentials to system keyring")
