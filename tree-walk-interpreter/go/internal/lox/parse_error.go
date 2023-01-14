package lox

type parseError struct {
}

func newParseError(token *Token, message string) parseError {
	errorWithToken(token, message)
	return parseError{}
}

func (e *parseError) Error() string {
	return "parse error"
}
