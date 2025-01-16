package wnet

import "errors"

var (
	ErrConnectionNotFound = errors.New("connection not found")
	ErrStreamNotFound     = errors.New("stream not found")
)
