package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	hadError bool
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: go-lox [script]")
		os.Exit(64)
	}

	if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}

}

func runFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("File not found.")
		return
	}

	run(string(bytes))

	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	var reader = bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, _, err := reader.ReadLine()
		if err == io.EOF {
			os.Exit(0)
		} else if err != nil {
			os.Exit(0)
		}

		if len(line) == 0 {
			continue
		}

		if string(line) == "exit" {
			break
		}

		run(string(line))
		hadError = false
	}
}

func run(source string) {
	fmt.Print(source)
}

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
	hadError = true
}
