package botil

func NewError(text string,code int32 ) error {
	return &LLError{text,code}
}

// errorString is a trivial implementation of error.
type LLError struct {
	s string
	code int32
}

func (e *LLError) Error() string {
	return e.s
}
func (e *LLError) Code() int32 {
	return e.code
}