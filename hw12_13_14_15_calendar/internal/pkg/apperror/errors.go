package apperror

import (
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("Not found")
var ErrBadRequest = errors.New("Bad request")
var ErrInternal = errors.New("Internal error")
var ErrTokenHasExpired = errors.New("Token has expired")
