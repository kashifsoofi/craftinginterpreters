package lox

type StmtVisitor interface {
	VisitBlockStmt(stmt *Block) any
	VisitClassStmt(stmt *Class) any
	VisitExpressionStmt(stmt *Expression) any
	VisitFunctionStmt(stmt *Function) any
	VisitIfStmt(stmt *If) any
	VisitPrintStmt(stmt *Print) any
	VisitReturnStmt(stmt *Return) any
	VisitVarStmt(stmt *Var) any
	VisitWhileStmt(stmt *While) any
}

type Stmt interface {
	Accept(v StmtVisitor) any
}

type Block struct {
	Statements []Stmt
}

func NewBlock(statements []Stmt) *Block {
	return &Block{
		Statements: statements,
	}
}

func (stmt *Block) Accept(v StmtVisitor) any {
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

func (stmt *Class) Accept(v StmtVisitor) any {
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

func (stmt *Expression) Accept(v StmtVisitor) any {
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

func (stmt *Function) Accept(v StmtVisitor) any {
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

func (stmt *If) Accept(v StmtVisitor) any {
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

func (stmt *Print) Accept(v StmtVisitor) any {
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

func (stmt *Return) Accept(v StmtVisitor) any {
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

func (stmt *Var) Accept(v StmtVisitor) any {
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

func (stmt *While) Accept(v StmtVisitor) any {
	return v.VisitWhileStmt(stmt)
}
