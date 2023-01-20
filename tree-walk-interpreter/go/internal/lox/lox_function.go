package lox

import (
	"fmt"
)

type loxFunction struct {
	declaration *Function
	closure     *environment
	initializer bool
}

func newLoxFunction(declaration *Function, closure *environment, initializer bool) *loxFunction {
	return &loxFunction{
		declaration: declaration,
		closure:     closure,
		initializer: initializer,
	}
}

func (f *loxFunction) arity() int {
	return len(f.declaration.Parameters)
}

func (f *loxFunction) call(interpreter *Interpreter, arguments []interface{}) (returnVal interface{}) {
	defer func() {
		if r := recover(); r != nil {
			if returnValue, ok := r.(returnControl); ok {
				if f.initializer {
					returnVal = f.closure.getAt(0, "this")
					return
				}

				returnVal = returnValue.value
				return
			}

			panic(r)
		}
	}()

	environment := newEnvironment(f.closure)
	for i, paramter := range f.declaration.Parameters {
		environment.define(paramter.Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.declaration.Body, environment)
	if f.initializer {
		returnVal = f.closure.getAt(0, "this")
	}
	return
}

func (f *loxFunction) bind(instance *loxInstance) *loxFunction {
	environment := newEnvironment(f.closure)
	environment.define("this", instance)
	return newLoxFunction(f.declaration, environment, f.initializer)
}

func (f *loxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}
