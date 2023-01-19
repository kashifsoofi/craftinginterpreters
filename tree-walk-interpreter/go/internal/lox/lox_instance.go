package lox

type loxInstance struct {
	class *loxClass
}

func newLoxInstance(class *loxClass) *loxInstance {
	return &loxInstance{
		class: class,
	}
}

func (i *loxInstance) String() string {
	return i.class.name + " instance"
}
