package repository

import "github.com/pkg/errors"

var UpdateStaleObjectErr = errors.New("attempted to update a stale object")
