package lox

type loxClass struct {
	name string
}

func newLoxClass(name string) *loxClass {
	return &loxClass{
		name: name,
	}
}
func (c *loxClass) arity() int {
	return 0
}

func (c *loxClass) call(interpreter *Interpreter, arguments []interface{}) (returnVal interface{}) {
	instance := newLoxInstance(c)
	return instance
}

func (c *loxClass) String() string {
	return c.name
}
