package storage

import "errors"

var ErrURLNotFound = errors.New("URL not found")
var ErrURLExist = errors.New("URL with the same alias already exists")
