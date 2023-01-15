package lox

import "fmt"

type environment struct {
	enclosing *environment
	values    map[string]interface{}
}

func newEnvironment(enclosing *environment) *environment {
	return &environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (e *environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *environment) assign(name *Token, value interface{}) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}

	panic(newRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme)))
}

func (e *environment) get(name *Token) interface{} {
	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	panic(newRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme)))
}
