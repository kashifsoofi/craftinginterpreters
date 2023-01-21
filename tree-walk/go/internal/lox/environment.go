package lox

import "fmt"

type environment struct {
	enclosing *environment
	values    map[string]any
}

func newEnvironment(enclosing *environment) *environment {
	return &environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (e *environment) define(name string, value any) {
	e.values[name] = value
}

func (e *environment) assign(name *Token, value any) {
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

func (e *environment) get(name *Token) any {
	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	panic(newRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme)))
}

func (e *environment) getAt(distance int, name string) any {
	return e.ancestor(distance).values[name]
}

func (e *environment) assignAt(distance int, name *Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}

func (e *environment) ancestor(distance int) *environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}

	return environment
}
