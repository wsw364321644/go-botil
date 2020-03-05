package botil

import (
	"runtime/debug"
)
var  fstack=false
func OpenErrorStack(){
	fstack=true
}
func CloseErrorStack(){
	fstack=false
}
func NewError(text string,code int32) *LLError {
	return &LLError{&text,code,[]byte{}}
}
func ConvertError(err error) error {
	if(err==nil){
		return nil
	}
	str:=err.Error()
	return &LLError{&str,0x0fffffff,[]byte{}}
}
// errorString is a trivial implementation of error.
type LLError struct {
	s *string
	code int32
	stack []byte
}
func (e *LLError) Record() error{
	//_,file,line,_:=runtime.Caller(1)
	if fstack {
		return &LLError{e.s, e.code, debug.Stack()}
	}else{
		return e
	}
}
func (e *LLError) Error() string {
	return *e.s
}
func (e *LLError) Code() int32 {
	return e.code
}
func (e *LLError) Message() string {
	return *e.s+"\n"+string(e.stack)
}