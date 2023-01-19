package lox

type loxClass struct {
	name    string
	methods map[string]*loxFunction
}

func newLoxClass(name string, methods map[string]*loxFunction) *loxClass {
	return &loxClass{
		name:    name,
		methods: methods,
	}
}

func (c *loxClass) arity() int {
	return 0
}

func (c *loxClass) call(interpreter *Interpreter, arguments []interface{}) (returnVal interface{}) {
	instance := newLoxInstance(c)
	return instance
}

func (c *loxClass) findMethod(name string) *loxFunction {
	method, ok := c.methods[name]
	if ok {
		return method
	}

	return nil
}

func (c *loxClass) String() string {
	return c.name
}
