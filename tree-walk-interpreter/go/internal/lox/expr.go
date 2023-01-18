package lox

type ExprVisitor interface {
	VisitAssignExpr(expr *Assign) interface{}
	VisitBinaryExpr(expr *Binary) interface{}
	VisitCallExpr(expr *Call) interface{}
	VisitGetExpr(expr *Get) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitLogicalExpr(expr *Logical) interface{}
	VisitSetExpr(expr *Set) interface{}
	VisitSuperExpr(expr *Super) interface{}
	VisitThisExpr(expr *This) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitVariableExpr(expr *Variable) interface{}
}

type Expr interface {
	Accept(v ExprVisitor) interface{}
}

type Assign struct {
	Name *Token
	Value Expr
}

func NewAssign(name *Token, value Expr) *Assign {
	return &Assign{
		Name: name,
		Value: value,
	}
}

func (expr *Assign) Accept(v ExprVisitor) interface{} {
	return v.VisitAssignExpr(expr)
}

type Binary struct {
	Left Expr
	Operator *Token
	Right Expr
}

func NewBinary(left Expr, operator *Token, right Expr) *Binary {
	return &Binary{
		Left: left,
		Operator: operator,
		Right: right,
	}
}

func (expr *Binary) Accept(v ExprVisitor) interface{} {
	return v.VisitBinaryExpr(expr)
}

type Call struct {
	Callee Expr
	Paren *Token
	Arguments []Expr
}

func NewCall(callee Expr, paren *Token, arguments []Expr) *Call {
	return &Call{
		Callee: callee,
		Paren: paren,
		Arguments: arguments,
	}
}

func (expr *Call) Accept(v ExprVisitor) interface{} {
	return v.VisitCallExpr(expr)
}

type Get struct {
	Object Expr
	Name *Token
}

func NewGet(object Expr, name *Token) *Get {
	return &Get{
		Object: object,
		Name: name,
	}
}

func (expr *Get) Accept(v ExprVisitor) interface{} {
	return v.VisitGetExpr(expr)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{
		Expression: expression,
	}
}

func (expr *Grouping) Accept(v ExprVisitor) interface{} {
	return v.VisitGroupingExpr(expr)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(value interface{}) *Literal {
	return &Literal{
		Value: value,
	}
}

func (expr *Literal) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteralExpr(expr)
}

type Logical struct {
	Left Expr
	Operator *Token
	Right Expr
}

func NewLogical(left Expr, operator *Token, right Expr) *Logical {
	return &Logical{
		Left: left,
		Operator: operator,
		Right: right,
	}
}

func (expr *Logical) Accept(v ExprVisitor) interface{} {
	return v.VisitLogicalExpr(expr)
}

type Set struct {
	Object Expr
	Name *Token
	Value Expr
}

func NewSet(object Expr, name *Token, value Expr) *Set {
	return &Set{
		Object: object,
		Name: name,
		Value: value,
	}
}

func (expr *Set) Accept(v ExprVisitor) interface{} {
	return v.VisitSetExpr(expr)
}

type Super struct {
	Keyword *Token
	Method *Token
}

func NewSuper(keyword *Token, method *Token) *Super {
	return &Super{
		Keyword: keyword,
		Method: method,
	}
}

func (expr *Super) Accept(v ExprVisitor) interface{} {
	return v.VisitSuperExpr(expr)
}

type This struct {
	Keyword *Token
}

func NewThis(keyword *Token) *This {
	return &This{
		Keyword: keyword,
	}
}

func (expr *This) Accept(v ExprVisitor) interface{} {
	return v.VisitThisExpr(expr)
}

type Unary struct {
	Operator *Token
	Right Expr
}

func NewUnary(operator *Token, right Expr) *Unary {
	return &Unary{
		Operator: operator,
		Right: right,
	}
}

func (expr *Unary) Accept(v ExprVisitor) interface{} {
	return v.VisitUnaryExpr(expr)
}

type Variable struct {
	Name *Token
}

func NewVariable(name *Token) *Variable {
	return &Variable{
		Name: name,
	}
}

func (expr *Variable) Accept(v ExprVisitor) interface{} {
	return v.VisitVariableExpr(expr)
}

