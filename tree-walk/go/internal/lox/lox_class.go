package lox

type loxClass struct {
	name       string
	superclass *loxClass
	methods    map[string]*loxFunction
}

func newLoxClass(name string, superclass *loxClass, methods map[string]*loxFunction) *loxClass {
	return &loxClass{
		name:       name,
		superclass: superclass,
		methods:    methods,
	}
}

func (c *loxClass) arity() int {
	initializer := c.findMethod("init")
	if initializer == nil {
		return 0
	}

	return initializer.arity()
}

func (c *loxClass) call(interpreter *Interpreter, arguments []any) (returnVal any) {
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

	if c.superclass != nil {
		return c.superclass.findMethod(name)
	}

	return nil
}

func (c *loxClass) String() string {
	return c.name
}
