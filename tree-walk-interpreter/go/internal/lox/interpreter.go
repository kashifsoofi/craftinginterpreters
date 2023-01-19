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

func (i *Interpreter) VisitAssignExpr(expr *Assign) interface{} {
	value := i.evaluate(expr.Value)

	distance, ok := i.locals[expr]
	if ok {
		i.environment.assignAt(distance, expr.Name, value)
	} else {
		i.globals.assign(expr.Name, value)
	}
	return value
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) interface{} {
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

func (i *Interpreter) VisitCallExpr(expr *Call) interface{} {
	callee := i.evaluate(expr.Callee)

	arguments := make([]interface{}, 0)
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

func (i *Interpreter) VisitGetExpr(expr *Get) interface{} {
	object := i.evaluate(expr.Object)
	if instance, ok := object.(*loxInstance); ok {
		return instance.get(expr.Name)
	}

	panic(newRuntimeError(expr.Name, "Only instances have properties."))
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitLogicalExpr(expr *Logical) interface{} {
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

func (i *Interpreter) VisitSetExpr(expr *Set) interface{} {
	object := i.evaluate(expr.Object)
	instance, ok := object.(*loxInstance)
	if !ok {
		panic(newRuntimeError(expr.Name, "Only instances have fields."))
	}

	value := i.evaluate(expr.Value)
	instance.set(expr.Name, value)
	return nil
}

func (i *Interpreter) VisitSuperExpr(expr *Super) interface{} {
	return nil
}

func (i *Interpreter) VisitThisExpr(expr *This) interface{} {
	return i.lookupVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) interface{} {
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

func (i *Interpreter) VisitVariableExpr(expr *Variable) interface{} {
	return i.lookupVariable(expr.Name, expr)
}

func (i *Interpreter) VisitBlockStmt(stmt *Block) interface{} {
	i.executeBlock(stmt.Statements, newEnvironment(i.environment))
	return nil
}

func (i *Interpreter) VisitClassStmt(stmt *Class) interface{} {
	i.environment.define(stmt.Name.Lexeme, nil)

	methods := map[string]*loxFunction{}
	for _, method := range stmt.Methods {
		function := newLoxFunction(method, i.environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = function
	}

	class := newLoxClass(stmt.Name.Lexeme, methods)
	i.environment.assign(stmt.Name, class)

	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *Expression) interface{} {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitFunctionStmt(stmt *Function) interface{} {
	function := newLoxFunction(stmt, i.environment, false)
	i.environment.define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *If) interface{} {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringify(value))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *Return) interface{} {
	var value interface{} = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}

	panic(newReturnControl(value))
}

func (i *Interpreter) VisitVarStmt(stmt *Var) interface{} {
	var value interface{} = nil
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.environment.define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *While) interface{} {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
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

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}

	if b, ok := object.(bool); ok {
		return b
	}

	return true
}

func (i *Interpreter) isEqual(a, b interface{}) bool {
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

func (i *Interpreter) lookupVariable(name *Token, expr Expr) interface{} {
	distance, ok := i.locals[expr]
	if ok {
		return i.environment.getAt(distance, name.Lexeme)
	} else {
		return i.globals.get(name)
	}
}

func checkNumberOperand(token *Token, operand interface{}) {
	if _, ok := operand.(float64); ok {
		return
	}
	panic(newRuntimeError(token, "Operand must be a number."))
}

func checkNumberOperands(token *Token, left, right interface{}) {
	_, lok := left.(float64)
	_, rok := right.(float64)
	if lok && rok {
		return
	}

	panic(newRuntimeError(token, "Operands must be numbers."))
}

func stringify(object interface{}) string {
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
