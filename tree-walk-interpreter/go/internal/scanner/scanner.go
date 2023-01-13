package scanner

import (
	"strconv"

	loxError "github.com/kashifsoofi/go-lox/internal/error"
)

var (
	source  []rune
	tokens  []*Token
	start   int
	current int
	line    int
)

var keywordsTokenTypeMap = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

func ScanTokens(src string) []*Token {
	source = []rune(src)
	tokens = make([]*Token, 0)
	start = 0
	current = 0
	line = 1

	scanTokens()

	return tokens
}

func scanTokens() {
	for !isAtEnd() {
		// We are at the start of the next lexeme.
		start = current
		scanToken()
	}

	tokens = append(tokens, NewToken(EOF, "", nil, line))
}

func isAtEnd() bool {
	return current >= len(source)
}

func advance() rune {
	r := source[current]
	current++
	return r
}

func peek() rune {
	if isAtEnd() {
		return rune(0)
	}

	return source[current]
}

func peekNext() rune {
	if current+1 >= len(source) {
		return rune(0)
	}

	return source[current+1]
}

func match(r rune) bool {
	if isAtEnd() {
		return false
	}

	if source[current] != r {
		return false
	}

	current++
	return true
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}

func scanToken() {
	r := advance()
	switch r {
	case '(':
		addToken(LeftParen)
	case ')':
		addToken(RightParen)
	case '{':
		addToken(LeftBrace)
	case '}':
		addToken(RightBrace)
	case ',':
		addToken(Comma)
	case '.':
		addToken(Dot)
	case '-':
		addToken(Minus)
	case '+':
		addToken(Plus)
	case ';':
		addToken(Semicolon)
	case '*':
		addToken(Star)
	case '!':
		if match('=') {
			addToken(BangEqual)
		} else {
			addToken(Bang)
		}
	case '=':
		if match('=') {
			addToken(EqualEqual)
		} else {
			addToken(Equal)
		}
	case '<':
		if match('=') {
			addToken(LessEqual)
		} else {
			addToken(Less)
		}
	case '>':
		if match('=') {
			addToken(GreaterEqual)
		} else {
			addToken(Greater)
		}
	case '/':
		if match('/') {
			for peek() != '\n' && !isAtEnd() {
				advance()
			}
		} else {
			addToken(Slash)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		line++
	case '"':
		scanString()
	default:
		if isDigit(r) {
			scanNumber()
		} else if isAlpha(r) {
			scanIdentifier()
		} else {
			loxError.Error(line, "Unexpected character.")
		}
	}
}

func scanString() {
	for peek() != '"' && !isAtEnd() {
		if peek() == '\n' {
			line++
		}
		advance()
	}

	if isAtEnd() {
		loxError.Error(line, "Unterminated string.")
		return
	}

	// The closing ".
	advance()

	value := source[start+1 : current-1]
	addTokenWithLiteral(String, string(value))
}

func scanNumber() {
	for isDigit(peek()) {
		advance()
	}

	if peek() == '.' && isDigit(peekNext()) {
		// Consume the "."
		advance()

		for isDigit(peek()) {
			advance()
		}
	}

	s := string(source[start:current])
	n, _ := strconv.ParseFloat(s, 64)
	addTokenWithLiteral(Number, n)
}

func scanIdentifier() {
	for isAlphaNumeric(peek()) {
		advance()
	}

	text := string(source[start:current])
	tokenType, ok := keywordsTokenTypeMap[text]
	if !ok {
		tokenType = Identifier
	}

	addToken(tokenType)
}

func addToken(tokenType TokenType) {
	addTokenWithLiteral(tokenType, nil)
}

func addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	text := string(source[start:current])
	tokens = append(tokens, NewToken(tokenType, text, literal, line))
}
