package lox

import (
	"fmt"
	"strings"
)

type Interpreter struct {
	globals     *environment
	environment *environment
	locals      map[Expr]int
}

func NewInterpreter() *Interpreter {
	g := newEnvironment(nil)
	g.define("clock", newClockNativeFunction())
	return &Interpreter{
		globals:     g,
		environment: g,
		locals:      make(map[Expr]int, 0),
	}
}

func (i *Interpreter) Interpret(statements []Stmt) {
	defer func() {
		if err := recover(); err != nil {
			if runtimeErr, ok := err.(runtimeError); ok {
				reportRuntimeError(runtimeErr)
				return
			}
			panic(err)
		}
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) VisitAssignExpr(expr *Assign) any {
	value := i.evaluate(expr.Value)

	distance, ok := i.locals[expr]
	if ok {
		i.environment.assignAt(distance, expr.Name, value)
	} else {
		i.globals.assign(expr.Name, value)
	}
	return value
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case TokenTypeGreater:
		checkNumberOperands(expr.Operator, left, right)
		left, _ := left.(float64)
		right, _ := right.(float64)
		return left > right
	case TokenTypeGreaterEqual:
		checkNumberOperands(expr.Operator, left, right)
		left, _ := left.(float64)
		right, _ := right.(float64)
		return left >= right
	case TokenTypeLess:
		checkNumberOperands(expr.Operator, left, right)
		left, _ := left.(float64)
		right, _ := right.(float64)
		return left < right
	case TokenTypeLessEqual:
		checkNumberOperands(expr.Operator, left, right)
		left, _ := left.(float64)
		right, _ := right.(float64)
		return left <= right
	case TokenTypeBangEqual:
		return !i.isEqual(left, right)
	case TokenTypeEqualEqual:
		return i.isEqual(left, right)
	case TokenTypeMinus:
		checkNumberOperands(expr.Operator, left, right)
		left, _ := left.(float64)
		right, _ := right.(float64)
		return left - right
	case TokenTypePlus:
		lf, lok := left.(float64)
		rf, rok := right.(float64)
		if lok && rok {
			return lf + rf
		}

		ls, lok := left.(string)
		rs, rok := right.(string)
		if lok && rok {
			return ls + rs
		}

		panic(newRuntimeError(expr.Operator, "Operands must be two numbers or two strings."))
	case TokenTypeSlash:
		checkNumberOperands(expr.Operator, left, right)
		left, _ := left.(float64)
		right, _ := right.(float64)
		return left / right
	case TokenTypeStar:
		checkNumberOperands(expr.Operator, left, right)
		left, _ := left.(float64)
		right, _ := right.(float64)
		return left * right
	}
	return nil
}

func (i *Interpreter) VisitCallExpr(expr *Call) any {
	callee := i.evaluate(expr.Callee)

	arguments := make([]any, 0)
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		panic(newRuntimeError(expr.Paren, "Can only call functions and classes."))
	}

	if len(arguments) != function.arity() {
		panic(newRuntimeError(expr.Paren, fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(arguments))))
	}

	return function.call(i, arguments)
}

func (i *Interpreter) VisitGetExpr(expr *Get) any {
	object := i.evaluate(expr.Object)
	if instance, ok := object.(*loxInstance); ok {
		return instance.get(expr.Name)
	}

	panic(newRuntimeError(expr.Name, "Only instances have properties."))
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) any {
	return expr.Value
}

func (i *Interpreter) VisitLogicalExpr(expr *Logical) any {
	left := i.evaluate(expr.Left)

	if expr.Operator.Type == TokenTypeOr {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitSetExpr(expr *Set) any {
	object := i.evaluate(expr.Object)
	instance, ok := object.(*loxInstance)
	if !ok {
		panic(newRuntimeError(expr.Name, "Only instances have fields."))
	}

	value := i.evaluate(expr.Value)
	instance.set(expr.Name, value)
	return value
}

func (i *Interpreter) VisitSuperExpr(expr *Super) any {
	distance := i.locals[expr]
	superclass, _ := i.environment.getAt(distance, "super").(*loxClass)

	instance, _ := i.environment.getAt(distance-1, "this").(*loxInstance)

	method := superclass.findMethod(expr.Method.Lexeme)
	if method == nil {
		panic(newRuntimeError(expr.Method, fmt.Sprintf("Undefined property '%s'.", expr.Method.Lexeme)))
	}

	return method.bind(instance)
}

func (i *Interpreter) VisitThisExpr(expr *This) any {
	return i.lookupVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) any {
	right := i.evaluate(expr.Right)
	switch expr.Operator.Type {
	case TokenTypeBang:
		return !i.isTruthy(right)
	case TokenTypeMinus:
		checkNumberOperand(expr.Operator, right)
		f, _ := right.(float64)
		return -f
	}

	// Unreachable.
	return nil
}

func (i *Interpreter) VisitVariableExpr(expr *Variable) any {
	return i.lookupVariable(expr.Name, expr)
}

func (i *Interpreter) VisitBlockStmt(stmt *Block) any {
	i.executeBlock(stmt.Statements, newEnvironment(i.environment))
	return nil
}

func (i *Interpreter) VisitClassStmt(stmt *Class) any {
	var superclass *loxClass = nil
	if stmt.Superclass != nil {
		object := i.evaluate(stmt.Superclass)
		c, ok := object.(*loxClass)
		if !ok {
			panic(newRuntimeError(stmt.Superclass.Name, "Superclass must be a class."))
		}
		superclass = c
	}
	i.environment.define(stmt.Name.Lexeme, nil)

	if superclass != nil {
		i.environment = newEnvironment(i.environment)
		i.environment.define("super", superclass)
	}

	methods := map[string]*loxFunction{}
	for _, method := range stmt.Methods {
		function := newLoxFunction(method, i.environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = function
	}

	class := newLoxClass(stmt.Name.Lexeme, superclass, methods)

	if superclass != nil {
		i.environment = i.environment.enclosing
	}

	i.environment.assign(stmt.Name, class)

	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *Expression) any {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitFunctionStmt(stmt *Function) any {
	function := newLoxFunction(stmt, i.environment, false)
	i.environment.define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *If) any {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) any {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringify(value))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *Return) any {
	var value any = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}

	panic(newReturnControl(value))
}

func (i *Interpreter) VisitVarStmt(stmt *Var) any {
	var value any = nil
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.environment.define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *While) any {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *environment) {
	previousEnvironment := i.environment

	defer func() {
		i.environment = previousEnvironment
	}()

	i.environment = environment

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) isTruthy(object any) bool {
	if object == nil {
		return false
	}

	if b, ok := object.(bool); ok {
		return b
	}

	return true
}

func (i *Interpreter) isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) lookupVariable(name *Token, expr Expr) any {
	distance, ok := i.locals[expr]
	if ok {
		return i.environment.getAt(distance, name.Lexeme)
	} else {
		return i.globals.get(name)
	}
}

func checkNumberOperand(token *Token, operand any) {
	if _, ok := operand.(float64); ok {
		return
	}
	panic(newRuntimeError(token, "Operand must be a number."))
}

func checkNumberOperands(token *Token, left, right any) {
	_, lok := left.(float64)
	_, rok := right.(float64)
	if lok && rok {
		return
	}

	panic(newRuntimeError(token, "Operands must be numbers."))
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}

	if f, ok := object.(float64); ok {
		text := fmt.Sprintf("%v", f)
		if strings.HasSuffix(text, ".0") {
			text = strings.TrimRight(text, ".0")
		}
		return text
	}

	return fmt.Sprintf("%v", object)
}
