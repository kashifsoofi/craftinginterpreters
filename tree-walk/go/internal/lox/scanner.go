package lox

import (
	"strconv"
)

var keywordsTokenTypeMap = map[string]TokenType{
	"and":    TokenTypeAnd,
	"class":  TokenTypeClass,
	"else":   TokenTypeElse,
	"false":  TokenTypeFalse,
	"for":    TokenTypeFor,
	"fun":    TokenTypeFun,
	"if":     TokenTypeIf,
	"nil":    TokenTypeNil,
	"or":     TokenTypeOr,
	"print":  TokenTypePrint,
	"return": TokenTypeReturn,
	"super":  TokenTypeSuper,
	"this":   TokenTypeThis,
	"true":   TokenTypeTrue,
	"var":    TokenTypeVar,
	"while":  TokenTypeWhile,
}

type Scanner struct {
	source  []rune
	tokens  []*Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  []rune(source),
		tokens:  make([]*Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []*Token {
	for !s.isAtEnd() {
		// We are at the start of the next lexeme.
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, NewToken(TokenTypeEOF, "", nil, s.line))

	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	r := s.source[s.current]
	s.current++
	return r
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}

	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return rune(0)
	}

	return s.source[s.current+1]
}

func (s *Scanner) match(r rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != r {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (s *Scanner) isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r == '_'
}

func (s *Scanner) isAlphaNumeric(r rune) bool {
	return s.isAlpha(r) || s.isDigit(r)
}

func (s *Scanner) scanToken() {
	r := s.advance()
	switch r {
	case '(':
		s.addToken(TokenTypeLeftParen)
	case ')':
		s.addToken(TokenTypeRightParen)
	case '{':
		s.addToken(TokenTypeLeftBrace)
	case '}':
		s.addToken(TokenTypeRightBrace)
	case ',':
		s.addToken(TokenTypeComma)
	case '.':
		s.addToken(TokenTypeDot)
	case '-':
		s.addToken(TokenTypeMinus)
	case '+':
		s.addToken(TokenTypePlus)
	case ';':
		s.addToken(TokenTypeSemicolon)
	case '*':
		s.addToken(TokenTypeStar)
	case '!':
		if s.match('=') {
			s.addToken(TokenTypeBangEqual)
		} else {
			s.addToken(TokenTypeBang)
		}
	case '=':
		if s.match('=') {
			s.addToken(TokenTypeEqualEqual)
		} else {
			s.addToken(TokenTypeEqual)
		}
	case '<':
		if s.match('=') {
			s.addToken(TokenTypeLessEqual)
		} else {
			s.addToken(TokenTypeLess)
		}
	case '>':
		if s.match('=') {
			s.addToken(TokenTypeGreaterEqual)
		} else {
			s.addToken(TokenTypeGreater)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(TokenTypeSlash)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.scanString()
	default:
		if s.isDigit(r) {
			s.scanNumber()
		} else if s.isAlpha(r) {
			s.scanIdentifier()
		} else {
			error(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		error(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(TokenTypeString, string(value))
}

func (s *Scanner) scanNumber() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	literal := string(s.source[s.start:s.current])
	n, _ := strconv.ParseFloat(literal, 64)
	s.addTokenWithLiteral(TokenTypeNumber, n)
}

func (s *Scanner) scanIdentifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.source[s.start:s.current])
	tokenType, ok := keywordsTokenTypeMap[text]
	if !ok {
		tokenType = TokenTypeIdentifier
	}

	s.addToken(tokenType)
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	text := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}
