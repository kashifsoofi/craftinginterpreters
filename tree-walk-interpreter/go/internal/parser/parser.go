package parser

import (
	loxError "github.com/kashifsoofi/go-lox/internal/error"
	"github.com/kashifsoofi/go-lox/internal/scanner"
)

var (
	tokens  []*scanner.Token
	current int
)

func Parse(scannedTokens []*scanner.Token) Expr {
	tokens = scannedTokens
	current = 0

	return parse()
}

func parse() Expr {
	var expr Expr
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(ParseError); ok {
				expr = nil
				return
			}
			panic(err)
		}
	}()

	expr = expression()
	return expr
}

func expression() Expr {
	return equality()
}

func equality() Expr {
	expr := comparison()

	for match(scanner.BangEqual, scanner.EqualEqual) {
		operator := previous()
		right := comparison()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func comparison() Expr {
	expr := term()

	for match(scanner.Greater, scanner.GreaterEqual, scanner.Less, scanner.LessEqual) {
		operator := previous()
		right := term()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func term() Expr {
	expr := factor()

	for match(scanner.Minus, scanner.Plus) {
		operator := previous()
		right := factor()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func factor() Expr {
	expr := unary()

	for match(scanner.Slash, scanner.Star) {
		operator := previous()
		right := unary()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func unary() Expr {
	if match(scanner.Bang, scanner.Minus) {
		operator := previous()
		right := unary()
		return NewUnary(operator, right)
	}

	return primary()
}

func primary() Expr {
	if match(scanner.False) {
		return NewLiteral(false)
	}
	if match(scanner.True) {
		return NewLiteral(true)
	}
	if match(scanner.Nil) {
		return NewLiteral(nil)
	}

	if match(scanner.Number, scanner.String) {
		return NewLiteral(previous().Literal)
	}

	if match(scanner.LeftParen) {
		expr := expression()
		consume(scanner.RightParen, "Expect ')' after expression.")
		return NewGrouping(expr)
	}

	panic(error(peek(), "Expect expression."))
}

func match(tokenTypes ...scanner.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if check(tokenType) {
			advance()
			return true
		}
	}

	return false
}

func check(tokenType scanner.TokenType) bool {
	if isAtEnd() {
		return false
	}

	return peek().Type == tokenType
}

func advance() *scanner.Token {
	if !isAtEnd() {
		current++
	}
	return previous()
}

func isAtEnd() bool {
	return peek().Type == scanner.EOF
}

func peek() *scanner.Token {
	return tokens[current]
}

func previous() *scanner.Token {
	return tokens[current-1]
}

func consume(tokenType scanner.TokenType, message string) *scanner.Token {
	if check(tokenType) {
		return advance()
	}

	panic(error(previous(), message))
}

func error(token *scanner.Token, message string) ParseError {
	if token.Type == scanner.EOF {
		loxError.ErrorWithWhere(token.Line, "at end", message)
	} else {
		loxError.ErrorWithWhere(token.Line, "at '"+token.Lexeme+"'", message)
	}
	return NewParseError(token, message)
}

func synchronize() {
	advance()

	for !isAtEnd() {
		if previous().Type == scanner.Semicolon {
			return
		}

		switch peek().Type {
		case scanner.Class:
			return
		case scanner.Fun:
			return
		case scanner.Var:
			return
		case scanner.For:
			return
		case scanner.If:
			return
		case scanner.While:
			return
		case scanner.Print:
			return
		case scanner.Return:
			return
		}

		advance()
	}
}
