package lox

type returnControl struct {
	value interface{}
}

func newReturnControl(value interface{}) returnControl {
	return returnControl{
		value: value,
	}
}

func (e returnControl) Error() string {
	return "return value"
}
