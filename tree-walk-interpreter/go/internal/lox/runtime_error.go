package lox

import "fmt"

type runtimeError struct {
	token   *Token
	message string
}

func newRuntimeError(token *Token, message string) runtimeError {
	return runtimeError{
		token:   token,
		message: message,
	}
}

func (e *runtimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", e.message, e.token.Line)
}
