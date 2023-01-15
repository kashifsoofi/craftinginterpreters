package lox

type Parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() []Stmt {
	statements := make([]Stmt, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(parseError); ok {
				p.synchronize()
				return
			}
			panic(err)
		}
	}()

	if p.match(TokenTypeVar) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(TokenTypeIdentifier, "Expect variable name.")

	var initializer Expr = nil
	if p.match(TokenTypeEqual) {
		initializer = p.expression()
	}

	p.consume(TokenTypeSemicolon, "Expect ';' after variable declaration.")
	return NewVar(name, initializer)
}

func (p *Parser) statement() Stmt {
	if p.match(TokenTypePrint) {
		return p.printStatement()
	}
	if p.match(TokenTypeLeftBrace) {
		return NewBlock(p.block())
	}

	return p.expressionStatement()
}

func (p *Parser) block() []Stmt {
	statements := make([]Stmt, 0)

	for !p.check(TokenTypeRightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(TokenTypeRightBrace, "Expect '}' after block.")
	return statements
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(TokenTypeSemicolon, "Expect ';' after expression.")
	return NewExpression(expr)
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(TokenTypeSemicolon, "Expect ';' after value.")
	return NewPrint(value)
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.equality()

	if p.match(TokenTypeEqual) {
		equals := p.previous()
		value := p.assignment()

		if variable, ok := expr.(*Variable); ok {
			name := variable.Name
			return NewAssign(name, value)
		}

		newParseError(equals, "Invalid assignment target.")
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(TokenTypeBangEqual, TokenTypeEqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(TokenTypeGreater, TokenTypeGreaterEqual, TokenTypeLess, TokenTypeLessEqual) {
		operator := p.previous()
		right := p.term()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(TokenTypeMinus, TokenTypePlus) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(TokenTypeSlash, TokenTypeStar) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(TokenTypeBang, TokenTypeMinus) {
		operator := p.previous()
		right := p.unary()
		return NewUnary(operator, right)
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(TokenTypeFalse) {
		return NewLiteral(false)
	}
	if p.match(TokenTypeTrue) {
		return NewLiteral(true)
	}
	if p.match(TokenTypeNil) {
		return NewLiteral(nil)
	}

	if p.match(TokenTypeNumber, TokenTypeString) {
		return NewLiteral(p.previous().Literal)
	}

	if p.match(TokenTypeIdentifier) {
		return NewVariable(p.previous())
	}

	if p.match(TokenTypeLeftParen) {
		expr := p.expression()
		p.consume(TokenTypeRightParen, "Expect ')' after expression.")
		return NewGrouping(expr)
	}

	panic(newParseError(p.peek(), "Expect expression."))
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == TokenTypeEOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType TokenType, message string) *Token {
	if p.check(tokenType) {
		return p.advance()
	}

	panic(newParseError(p.peek(), message))
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == TokenTypeSemicolon {
			return
		}

		switch p.peek().Type {
		case TokenTypeClass:
			return
		case TokenTypeFun:
			return
		case TokenTypeVar:
			return
		case TokenTypeFor:
			return
		case TokenTypeIf:
			return
		case TokenTypeWhile:
			return
		case TokenTypePrint:
			return
		case TokenTypeReturn:
			return
		}

		p.advance()
	}
}
