package hook

var serverStartedCBs []HookFunc
func AddServerstartedCallback(infunc HookFunc){
	serverStartedCBs= append(serverStartedCBs, infunc)
}
func CallServerstartedCallback(){
	for _,f:=range serverStartedCBs{
		f()
	}
	serverStartedCBs=nil
}