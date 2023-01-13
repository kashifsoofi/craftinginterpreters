package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	loxError "github.com/kashifsoofi/go-lox/internal/error"
	"github.com/kashifsoofi/go-lox/internal/scanner"
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

	if loxError.HadError {
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
		loxError.HadError = false
	}
}

func run(source string) {
	tokens := scanner.ScanTokens(source)

	fmt.Printf("length: %d\n", len(tokens))

	for _, token := range tokens {
		fmt.Printf("%v\n", token)
	}
}
