package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: go-lox [script]")
		os.Exit(64)
	}
	outputDir := os.Args[1]

	exprTypeNames := []string{
		"Assign",
		"Binary",
		"Call",
		"Get",
		"Grouping",
		"Literal",
		"Logical",
		"Set",
		"Super",
		"This",
		"Unary",
		"Variable",
	}
	exprTypes := map[string][]string{
		"Assign":   {"Name *Token", "Value Expr"},
		"Binary":   {"Left Expr", "Operator *Token", "Right Expr"},
		"Call":     {"Callee Expr", "Paren *Token", "Arguments []Expr"},
		"Get":      {"Object Expr", "Name *Token"},
		"Grouping": {"Expression Expr"},
		"Literal":  {"Value interface{}"},
		"Logical":  {"Left Expr", "Operator *Token", "Right Expr"},
		"Set":      {"Object Expr", "Name *Token", "Value Expr"},
		"Super":    {"Keyword *Token", "Method *Token"},
		"This":     {"Keyword *Token"},
		"Unary":    {"Operator *Token", "Right Expr"},
		"Variable": {"Name *Token"},
	}
	generateAst(outputDir, "Expr", exprTypeNames, exprTypes)

	stmtTypeNames := []string{
		"Block",
		"Class",
		"Expression",
		"Function",
		"If",
		"Print",
		"Return",
		"Var",
		"While",
	}
	stmtTypes := map[string][]string{
		"Block":      {"Statements []Stmt"},
		"Class":      {"Name *Token", "Superclass *Variable", "Methods []*Function"},
		"Expression": {"Expression Expr"},
		"Function":   {"Name *Token", "Parameters []*Token", "Body []Stmt"},
		"If":         {"Condition Expr", "ThenBranch Stmt", "ElseBranch Stmt"},
		"Print":      {"Expression Expr"},
		"Return":     {"Keyword *Token", "Value Expr"},
		"Var":        {"Name *Token", "Initializer Expr"},
		"While":      {"Condition Expr", "Body Stmt"},
	}
	generateAst(outputDir, "Stmt", stmtTypeNames, stmtTypes)
}

func generateAst(outputDir, baseName string, typeNames []string, types map[string][]string) {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "package lox")
	fmt.Fprintln(f, "")

	generateVisitor(f, baseName, typeNames)

	fmt.Fprintf(f, "type %s interface {\n", baseName)
	fmt.Fprintf(f, "\tAccept(v %sVisitor) interface{}\n", baseName)
	fmt.Fprintf(f, "}\n")
	fmt.Fprintln(f, "")

	for _, typeName := range typeNames {
		generateType(f, baseName, typeName, types[typeName])
	}
}

func generateVisitor(f *os.File, baseName string, typeNames []string) {
	fmt.Fprintf(f, "type %sVisitor interface {\n", baseName)
	for _, typeName := range typeNames {
		fmt.Fprintf(f, "\tVisit%s%s(expr *%s) interface{}\n", typeName, baseName, typeName)
	}
	fmt.Fprintf(f, "}\n")
	fmt.Fprintln(f, "")
}

func generateType(f *os.File, baseName, typeName string, fields []string) {
	fmt.Fprintf(f, "type %s struct {\n", typeName)
	for _, field := range fields {
		fmt.Fprintf(f, "\t%s\n", field)
	}
	fmt.Fprintf(f, "}\n")
	fmt.Fprintln(f, "")
	// New func
	fmt.Fprintf(f, "func New%s(", typeName)
	for i, field := range fields {
		fieldName, fieldType, _ := strings.Cut(field, " ")
		fmt.Fprintf(f, "%s %s", strings.ToLower(fieldName), fieldType)
		if i+1 < len(fields) {
			fmt.Fprintf(f, ", ")
		}
	}
	fmt.Fprintf(f, ") *%s {\n", typeName)
	fmt.Fprintf(f, "\treturn &%s{\n", typeName)
	for _, field := range fields {
		fieldName, _, _ := strings.Cut(field, " ")
		fmt.Fprintf(f, "\t\t%s: %s,\n", fieldName, strings.ToLower(fieldName))
	}
	fmt.Fprintln(f, "\t}")
	fmt.Fprintln(f, "}")
	fmt.Fprintln(f, "")

	// Assign
	fmt.Fprintf(f, "func (%s *%s) Accept(v %sVisitor) interface{} {\n", strings.ToLower(baseName), typeName, baseName)
	fmt.Fprintf(f, "\treturn v.Visit%s%s(%s)\n", typeName, baseName, strings.ToLower(baseName))
	fmt.Fprintln(f, "}")
	fmt.Fprintln(f, "")
}
