package lox

import (
	"fmt"
	"strings"
)

type Interpreter struct {
	environment *environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: newEnvironment(nil),
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
	i.environment.assign(expr.Name, value)
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
	return nil
}

func (i *Interpreter) VisitGetExpr(expr *Get) interface{} {
	return nil
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitLogicalExpr(expr *Logical) interface{} {
	return nil
}

func (i *Interpreter) VisitSetExpr(expr *Set) interface{} {
	return nil
}

func (i *Interpreter) VisitSuperExpr(expr *Super) interface{} {
	return nil
}

func (i *Interpreter) VisitThisExpr(expr *This) interface{} {
	return nil
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
	return i.environment.get(expr.Name)
}

func (i *Interpreter) VisitBlockStmt(stmt *Block) interface{} {
	i.executeBlock(stmt.Statements, newEnvironment(i.environment))
	return nil
}

func (i *Interpreter) VisitClassStmt(stmt *Class) interface{} {
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *Expression) interface{} {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitFunctionStmt(stmt *Function) interface{} {
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *If) interface{} {
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringify(value))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *Return) interface{} {
	return nil
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
