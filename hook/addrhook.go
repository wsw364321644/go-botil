package hook

type  HookFunc func()
var hookMap map [interface{}]HookFunc
func init(){
	hookMap=make(map [interface{}]HookFunc)
}
func AddAddrHook(addr interface{},infunc HookFunc){
	hookMap[addr]=infunc
}
func CallAddrHook(addr interface{}){
	cbfunc,ok:=hookMap[addr]
	if ok{
		cbfunc()
	}
}