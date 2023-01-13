package errors

import (
	"fmt"
	"os"
)

var (
	HadError bool
)

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
	HadError = true
}
