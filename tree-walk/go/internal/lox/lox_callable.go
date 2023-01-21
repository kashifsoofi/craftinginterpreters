package lox

type LoxCallable interface {
	arity() int
	call(interpreter *Interpreter, arguments []interface{}) interface{}
}
