package lox

type returnControl struct {
	value any
}

func newReturnControl(value any) returnControl {
	return returnControl{
		value: value,
	}
}

func (e returnControl) Error() string {
	return "return value"
}
