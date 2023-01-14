package parser

import (
	"fmt"

	"github.com/kashifsoofi/go-lox/internal/scanner"
)

type ParseError struct {
	Token   *scanner.Token
	Message string
}

func NewParseError(token *scanner.Token, message string) ParseError {
	return ParseError{
		Token:   token,
		Message: message,
	}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%v: %s", e.Token, e.Message)
}
