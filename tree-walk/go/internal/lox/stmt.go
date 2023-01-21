package lox

type StmtVisitor interface {
	VisitBlockStmt(stmt *Block) interface{}
	VisitClassStmt(stmt *Class) interface{}
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitFunctionStmt(stmt *Function) interface{}
	VisitIfStmt(stmt *If) interface{}
	VisitPrintStmt(stmt *Print) interface{}
	VisitReturnStmt(stmt *Return) interface{}
	VisitVarStmt(stmt *Var) interface{}
	VisitWhileStmt(stmt *While) interface{}
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
	Name       *Token
	Superclass *Variable
	Methods    []*Function
}

func NewClass(name *Token, superclass *Variable, methods []*Function) *Class {
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
	Name       *Token
	Parameters []*Token
	Body       []Stmt
}

func NewFunction(name *Token, parameters []*Token, body []Stmt) *Function {
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
	Keyword *Token
	Value   Expr
}

func NewReturn(keyword *Token, value Expr) *Return {
	return &Return{
		Keyword: keyword,
		Value:   value,
	}
}

func (stmt *Return) Accept(v StmtVisitor) interface{} {
	return v.VisitReturnStmt(stmt)
}

type Var struct {
	Name        *Token
	Initializer Expr
}

func NewVar(name *Token, initializer Expr) *Var {
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
