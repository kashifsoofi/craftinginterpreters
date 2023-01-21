package lox

import "time"

type clockNativeFunction struct{}

func newClockNativeFunction() clockNativeFunction {
	return clockNativeFunction{}
}

func (c clockNativeFunction) arity() int {
	return 0
}

func (c clockNativeFunction) call(interpreter *Interpreter, arguments []any) any {
	return time.Now().UnixMilli()
}

func (c clockNativeFunction) String() string {
	return "<native fn>"
}
