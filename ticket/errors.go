package ticket

import "errors"

var ErrTypeAssertionAuth = errors.New("Could not parse credential data")
var ErrNotUnauthorized = errors.New("received 401")
var ErrNotFound = errors.New("received 404")
var ErrCredentialSaving = errors.New("Could not save credentials to system keyring")
