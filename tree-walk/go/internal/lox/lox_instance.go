package lox

import (
	"fmt"
)

type loxInstance struct {
	class  *loxClass
	fields map[string]any
}

func newLoxInstance(class *loxClass) *loxInstance {
	return &loxInstance{
		class:  class,
		fields: make(map[string]any),
	}
}

func (i *loxInstance) get(name *Token) any {
	if value, ok := i.fields[name.Lexeme]; ok {
		return value
	}

	method := i.class.findMethod(name.Lexeme)
	if method != nil {
		return method.bind(i)
	}

	panic(newRuntimeError(name, fmt.Sprintf("Undefined property '%s'.", name.Lexeme)))
}

func (i *loxInstance) set(name *Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *loxInstance) String() string {
	return i.class.name + " instance"
}
