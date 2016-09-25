package backend

import "errors"

// ErrWrongAPIKey is returned upon a wrong API key
var ErrWrongAPIKey = errors.New("wrong API key")
