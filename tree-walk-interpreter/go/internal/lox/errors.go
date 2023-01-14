package lox

import (
	"fmt"
	"os"
)

var (
	HadError bool
)

func error(line int, message string) {
	report(line, "", message)
}

func errorWithToken(token *Token, message string) {
	if token.Type == TokenTypeEOF {
		report(token.Line, "at end", message)
	} else {
		report(token.Line, "at '"+token.Lexeme+"'", message)
	}
}

func report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
	HadError = true
}
