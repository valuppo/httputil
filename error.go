package httputil

import "errors"

var ErrDecodeRequest = errors.New("Wrong request params format, see example in data")
var ErrInternalServerError = errors.New("Internal server error")
