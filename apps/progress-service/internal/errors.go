package internal

import "errors"

var ErrChildProfileForbidden = errors.New("child profile access is not allowed for current principal")
