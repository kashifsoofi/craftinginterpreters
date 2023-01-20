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
	initializer := c.findMethod("init")
	if initializer == nil {
		return 0
	}

	return initializer.arity()
}

func (c *loxClass) call(interpreter *Interpreter, arguments []interface{}) (returnVal interface{}) {
	instance := newLoxInstance(c)
	initializer := c.findMethod("init")
	if initializer != nil {
		initializer.bind(instance).call(interpreter, arguments)
	}

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
