package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/kashifsoofi/go-lox/internal/lox"
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
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("File not found.")
		return
	}

	run(string(bytes))

	if lox.HadError {
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
		lox.HadError = false
	}
}

func run(source string) {
	scanner := lox.NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := lox.NewParser(tokens)
	expression := parser.Parse()

	// Stop if there was a syntax error.
	if lox.HadError {
		return
	}

	astPrinter := lox.AstPrinter{}
	fmt.Println(astPrinter.Print(expression))
}