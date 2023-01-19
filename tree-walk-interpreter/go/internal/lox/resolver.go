package lox

type functionType int

const (
	functionTypeNone functionType = iota
	functionTypeFunction
)

type Resolver struct {
	interpreter         *Interpreter
	scopes              *stack
	currentFunctionType functionType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:         interpreter,
		scopes:              newStack(),
		currentFunctionType: functionTypeNone,
	}
}

func (r *Resolver) Resolve(statements []Stmt) {
	r.resolveStatements(statements)
}

func (r *Resolver) VisitAssignExpr(expr *Assign) interface{} {
	r.resolveExpression(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *Binary) interface{} {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *Call) interface{} {
	r.resolveExpression(expr.Callee)

	for _, argument := range expr.Arguments {
		r.resolveExpression(argument)
	}

	return nil
}

func (r *Resolver) VisitGetExpr(expr *Get) interface{} {
	r.resolveExpression(expr.Object)
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *Grouping) interface{} {
	r.resolveExpression(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *Literal) interface{} {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *Logical) interface{} {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil
}

func (r *Resolver) VisitSetExpr(expr *Set) interface{} {
	r.resolveExpression(expr.Value)
	r.resolveExpression(expr.Object)
	return nil
}

func (r *Resolver) VisitSuperExpr(expr *Super) interface{} {
	return nil
}

func (r *Resolver) VisitThisExpr(expr *This) interface{} {
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *Unary) interface{} {
	r.resolveExpression(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *Variable) interface{} {
	if !r.scopes.empty() {
		if declared, ok := r.scopes.peek()[expr.Name.Lexeme]; ok && !declared {
			newParseError(expr.Name, "Can't read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBlockStmt(stmt *Block) interface{} {
	r.beginScope()
	r.resolveStatements(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitClassStmt(stmt *Class) interface{} {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *Expression) interface{} {
	r.resolveExpression(stmt.Expression)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *Function) interface{} {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, functionTypeFunction)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *If) interface{} {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStatement(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *Print) interface{} {
	r.resolveExpression(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *Return) interface{} {
	if r.currentFunctionType == functionTypeNone {
		newParseError(stmt.Keyword, "Can't return from top-level code.")
	}

	if stmt.Value != nil {
		r.resolveExpression(stmt.Value)
	}

	return nil
}

func (r *Resolver) VisitVarStmt(stmt *Var) interface{} {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpression(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *While) interface{} {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.Body)
	return nil
}

func (r *Resolver) resolveExpression(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveStatement(statement Stmt) {
	statement.Accept(r)
}

func (r *Resolver) resolveStatements(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStatement(statement)
	}
}

func (r *Resolver) resolveLocal(expr Expr, name *Token) {
	l := len(r.scopes.values)
	for i := l - 1; i >= 0; i-- {
		if _, ok := r.scopes.values[i][name.Lexeme]; ok {
			r.interpreter.resolve(expr, l-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function *Function, functionType functionType) {
	enclosingFunctionType := r.currentFunctionType
	r.currentFunctionType = functionType
	r.beginScope()
	for _, param := range function.Parameters {
		r.declare(param)
		r.define(param)
	}
	r.resolveStatements(function.Body)
	r.endScope()
	r.currentFunctionType = enclosingFunctionType
}

func (r *Resolver) beginScope() {
	r.scopes.push(newScope())
}

func (r *Resolver) endScope() {
	r.scopes.pop()
}

func (r *Resolver) declare(name *Token) {
	if r.scopes.empty() {
		return
	}

	scope := r.scopes.peek()
	if _, ok := scope[name.Lexeme]; ok {
		newParseError(name, "Already a variable with this name in this scope.")
	}

	scope[name.Lexeme] = false
}

func (r *Resolver) define(name *Token) {
	if r.scopes.empty() {
		return
	}
	r.scopes.peek()[name.Lexeme] = true
}
