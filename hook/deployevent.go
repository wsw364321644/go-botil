package hook
var preDeployMysqlCBs []HookFunc
func AddPreDeployMysqlCallback(infunc HookFunc){
	preDeployMysqlCBs= append(preDeployMysqlCBs, infunc)
}
func CallPreDeployMysqlCallback(){
	for _,f:=range preDeployMysqlCBs{
		f()
	}
	preDeployMysqlCBs=nil
}