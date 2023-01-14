package lox

type TokenType int

const (
	// Singe-character tokens.
	TokenTypeLeftParen TokenType = iota
	TokenTypeRightParen
	TokenTypeLeftBrace
	TokenTypeRightBrace
	TokenTypeComma
	TokenTypeDot
	TokenTypeMinus
	TokenTypePlus
	TokenTypeSemicolon
	TokenTypeSlash
	TokenTypeStar

	// One or two character tokens.
	TokenTypeBang
	TokenTypeBangEqual
	TokenTypeEqual
	TokenTypeEqualEqual
	TokenTypeGreater
	TokenTypeGreaterEqual
	TokenTypeLess
	TokenTypeLessEqual

	// Literals.
	TokenTypeIdentifier
	TokenTypeString
	TokenTypeNumber

	// Keywords
	TokenTypeAnd
	TokenTypeClass
	TokenTypeElse
	TokenTypeFalse
	TokenTypeFun
	TokenTypeFor
	TokenTypeIf
	TokenTypeNil
	TokenTypeOr
	TokenTypePrint
	TokenTypeReturn
	TokenTypeSuper
	TokenTypeThis
	TokenTypeTrue
	TokenTypeVar
	TokenTypeWhile

	TokenTypeEOF
)

var tokenTypeNames = map[TokenType]string{
	TokenTypeLeftParen:    "LEFT_PAREN",
	TokenTypeRightParen:   "RIGHT_PAREN",
	TokenTypeLeftBrace:    "LEFT_BRACE",
	TokenTypeRightBrace:   "RIGHT_BRACE",
	TokenTypeComma:        "COMMA",
	TokenTypeDot:          "DOT",
	TokenTypeMinus:        "MINUS",
	TokenTypePlus:         "PLUS",
	TokenTypeSemicolon:    "SEMICOLON",
	TokenTypeSlash:        "SLASH",
	TokenTypeStar:         "STAR",
	TokenTypeBang:         "BANG",
	TokenTypeBangEqual:    "BANG_EQUAL",
	TokenTypeEqual:        "EQUAL",
	TokenTypeEqualEqual:   "EQUAL_EQUAL",
	TokenTypeGreater:      "GREATER",
	TokenTypeGreaterEqual: "GREATER_EQUAL",
	TokenTypeLess:         "LESS",
	TokenTypeLessEqual:    "LESS_EQUAL",
	TokenTypeIdentifier:   "IDENTIFIER",
	TokenTypeString:       "STRING",
	TokenTypeNumber:       "NUMBER",
	TokenTypeAnd:          "AND",
	TokenTypeClass:        "CLASS",
	TokenTypeElse:         "ELSE",
	TokenTypeFalse:        "FALSE",
	TokenTypeFun:          "FUN",
	TokenTypeFor:          "FOR",
	TokenTypeIf:           "IF",
	TokenTypeNil:          "NIL",
	TokenTypeOr:           "OR",
	TokenTypePrint:        "PRINT",
	TokenTypeReturn:       "RETURN",
	TokenTypeSuper:        "SUPER",
	TokenTypeThis:         "THIS",
	TokenTypeTrue:         "TRUE",
	TokenTypeVar:          "VAR",
	TokenTypeWhile:        "WHILE",
	TokenTypeEOF:          "EOF",
}

func (t TokenType) String() string {
	s, ok := tokenTypeNames[t]
	if ok {
		return s
	}

	return "Unknown"
}
