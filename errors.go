package service

import "errors"

var (
	errStop         = errors.New("stop serve signals")
	errWouldBlock   = errors.New("resource temporarily unavailable")
	errNotSupported = errors.New("only posix is supported")
	errWritePID     = errors.New("failed to write pid")
	errParsePID     = errors.New("failed to parse process")
)
