package lox

import "fmt"

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

	if p.match(TokenTypeClass) {
		return p.classDeclaration()
	}
	if p.match(TokenTypeFun) {
		return p.function("function")
	}
	if p.match(TokenTypeVar) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) classDeclaration() Stmt {
	name := p.consume(TokenTypeIdentifier, "Expect class name.")
	p.consume(TokenTypeLeftBrace, "Expect '{' before class body.")

	methods := make([]*Function, 0)
	for !p.check(TokenTypeRightBrace) && !p.isAtEnd() {
		method := p.function("method").(*Function)
		methods = append(methods, method)
	}

	p.consume(TokenTypeRightBrace, "Expect '}' after class body.")

	return NewClass(name, nil, methods)
}

func (p *Parser) function(kind string) Stmt {
	name := p.consume(TokenTypeIdentifier, fmt.Sprintf("Expect %s name.", kind))
	p.consume(TokenTypeLeftParen, fmt.Sprintf("Expect '(' after %s name.", kind))
	parameters := make([]*Token, 0)
	if !p.check(TokenTypeRightParen) {
		for {
			if len(parameters) >= 255 {
				newParseError(p.peek(), "Can't have more than 255 parameters.")
			}

			parameters = append(parameters, p.consume(TokenTypeIdentifier, "Expect parameter name."))
			if !p.match(TokenTypeComma) {
				break
			}
		}
	}
	p.consume(TokenTypeRightParen, "Expect ')' after parameters.")

	p.consume(TokenTypeLeftBrace, fmt.Sprintf("Expect '{' before %s body.", kind))
	body := p.block()
	return NewFunction(name, parameters, body)
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
	if p.match(TokenTypeFor) {
		return p.forStatement()
	}
	if p.match(TokenTypeIf) {
		return p.ifStatement()
	}
	if p.match(TokenTypePrint) {
		return p.printStatement()
	}
	if p.match(TokenTypeReturn) {
		return p.returnStatement()
	}
	if p.match(TokenTypeWhile) {
		return p.whileStatement()
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

func (p *Parser) forStatement() Stmt {
	p.consume(TokenTypeLeftParen, "Expect '(' after 'for'.")

	var initializer Stmt
	if p.match(TokenTypeSemicolon) {
		initializer = nil
	} else if p.match(TokenTypeVar) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr = nil
	if !p.check(TokenTypeSemicolon) {
		condition = p.expression()
	}
	p.consume(TokenTypeSemicolon, "Expect ';' after loop condition.")

	var increment Expr = nil
	if !p.check(TokenTypeRightParen) {
		increment = p.expression()
	}
	p.consume(TokenTypeRightParen, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = NewBlock([]Stmt{body, NewExpression(increment)})
	}

	if condition == nil {
		condition = NewLiteral(true)
	}
	body = NewWhile(condition, body)

	if initializer != nil {
		body = NewBlock([]Stmt{initializer, body})
	}

	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(TokenTypeLeftParen, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(TokenTypeRightParen, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt = nil
	if p.match(TokenTypeElse) {
		elseBranch = p.statement()
	}

	return NewIf(condition, thenBranch, elseBranch)
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(TokenTypeSemicolon, "Expect ';' after value.")
	return NewPrint(value)
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr = nil
	if !p.check(TokenTypeSemicolon) {
		value = p.expression()
	}

	p.consume(TokenTypeSemicolon, "Expect ';' after return value.")
	return NewReturn(keyword, value)
}

func (p *Parser) whileStatement() Stmt {
	p.consume(TokenTypeLeftParen, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(TokenTypeRightParen, "Expect ')' after condition.")
	body := p.statement()

	return NewWhile(condition, body)
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(TokenTypeEqual) {
		equals := p.previous()
		value := p.assignment()

		if variable, ok := expr.(*Variable); ok {
			name := variable.Name
			return NewAssign(name, value)
		} else if get, ok := expr.(*Get); ok {
			return NewSet(get, get.Name, value)
		}

		newParseError(equals, "Invalid assignment target.")
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(TokenTypeOr) {
		operator := p.previous()
		right := p.and()
		expr = NewLogical(expr, operator, right)
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(TokenTypeAnd) {
		operator := p.previous()
		right := p.equality()
		expr = NewLogical(expr, operator, right)
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

	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(TokenTypeLeftParen) {
			expr = p.finishCall(expr)
		} else if p.match(TokenTypeDot) {
			name := p.consume(TokenTypeIdentifier, "Expect property name after '.'.")
			expr = NewGet(expr, name)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := make([]Expr, 0)
	if !p.check(TokenTypeRightParen) {
		for {
			if len(arguments) >= 255 {
				newParseError(p.peek(), "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, p.expression())
			if !p.match(TokenTypeComma) {
				break
			}
		}
	}

	paren := p.consume(TokenTypeRightParen, "Expect ')' after arguments.")

	return NewCall(callee, paren, arguments)
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

	if p.match(TokenTypeThis) {
		return NewThis(p.previous())
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
