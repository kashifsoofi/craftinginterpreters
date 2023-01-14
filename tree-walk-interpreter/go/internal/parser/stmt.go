package parser

import (
	"github.com/kashifsoofi/go-lox/internal/scanner"
)

type StmtVisitor interface {
	VisitBlockStmt(expr *Block) interface{}
	VisitClassStmt(expr *Class) interface{}
	VisitExpressionStmt(expr *Expression) interface{}
	VisitFunctionStmt(expr *Function) interface{}
	VisitIfStmt(expr *If) interface{}
	VisitPrintStmt(expr *Print) interface{}
	VisitReturnStmt(expr *Return) interface{}
	VisitVarStmt(expr *Var) interface{}
	VisitWhileStmt(expr *While) interface{}
}

type Stmt interface {
	Accept(v StmtVisitor) interface{}
}

type Block struct {
	Statements []Stmt
}

func NewBlock(statements []Stmt) *Block {
	return &Block{
		Statements: statements,
	}
}

func (stmt *Block) Accept(v StmtVisitor) interface{} {
	return v.VisitBlockStmt(stmt)
}

type Class struct {
	Name       *scanner.Token
	Superclass *Variable
	Methods    []*Function
}

func NewClass(name *scanner.Token, superclass *Variable, methods []*Function) *Class {
	return &Class{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}
}

func (stmt *Class) Accept(v StmtVisitor) interface{} {
	return v.VisitClassStmt(stmt)
}

type Expression struct {
	Expression Expr
}

func NewExpression(expression Expr) *Expression {
	return &Expression{
		Expression: expression,
	}
}

func (stmt *Expression) Accept(v StmtVisitor) interface{} {
	return v.VisitExpressionStmt(stmt)
}

type Function struct {
	Name       *scanner.Token
	Parameters []*scanner.Token
	Body       []Stmt
}

func NewFunction(name *scanner.Token, parameters []*scanner.Token, body []Stmt) *Function {
	return &Function{
		Name:       name,
		Parameters: parameters,
		Body:       body,
	}
}

func (stmt *Function) Accept(v StmtVisitor) interface{} {
	return v.VisitFunctionStmt(stmt)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIf(condition Expr, thenbranch Stmt, elsebranch Stmt) *If {
	return &If{
		Condition:  condition,
		ThenBranch: thenbranch,
		ElseBranch: elsebranch,
	}
}

func (stmt *If) Accept(v StmtVisitor) interface{} {
	return v.VisitIfStmt(stmt)
}

type Print struct {
	Expression Expr
}

func NewPrint(expression Expr) *Print {
	return &Print{
		Expression: expression,
	}
}

func (stmt *Print) Accept(v StmtVisitor) interface{} {
	return v.VisitPrintStmt(stmt)
}

type Return struct {
	Keyword *scanner.Token
	Value   Expr
}

func NewReturn(keyword *scanner.Token, value Expr) *Return {
	return &Return{
		Keyword: keyword,
		Value:   value,
	}
}

func (stmt *Return) Accept(v StmtVisitor) interface{} {
	return v.VisitReturnStmt(stmt)
}

type Var struct {
	Name        *scanner.Token
	Initializer Expr
}

func NewVar(name *scanner.Token, initializer Expr) *Var {
	return &Var{
		Name:        name,
		Initializer: initializer,
	}
}

func (stmt *Var) Accept(v StmtVisitor) interface{} {
	return v.VisitVarStmt(stmt)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func NewWhile(condition Expr, body Stmt) *While {
	return &While{
		Condition: condition,
		Body:      body,
	}
}

func (stmt *While) Accept(v StmtVisitor) interface{} {
	return v.VisitWhileStmt(stmt)
}
