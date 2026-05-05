package cli

import "errors"

var (
	ErrLoadFileIsDirectory  = errors.New("load file is a directory")
	ErrLoadFileDoesNotExist = errors.New("load file does not exist")
	ErrInvalidFlag          = errors.New("invalid flag")
	ErrInvalidParserState   = errors.New("invalid parser state")
)
