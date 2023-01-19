package lox

import "fmt"

type loxInstance struct {
	class  *loxClass
	fields map[string]interface{}
}

func newLoxInstance(class *loxClass) *loxInstance {
	return &loxInstance{
		class:  class,
		fields: make(map[string]interface{}),
	}
}

func (i *loxInstance) get(name *Token) interface{} {
	if value, ok := i.fields[name.Lexeme]; ok {
		return value
	}

	panic(newRuntimeError(name, fmt.Sprintf("Undefined property '%s'.", name.Lexeme)))
}

func (i *loxInstance) set(name *Token, value interface{}) {
	i.fields[name.Lexeme] = value
}

func (i *loxInstance) String() string {
	return i.class.name + " instance"
}
