package botil

import "github.com/bitly/go-simplejson"

func GetRPCResult(result interface{},id int)[]byte{
	json:=simplejson.New()
	json.Set("result",result)

	json.Set("id",id)
	b,_:=json.Encode()
	return b
}

func GetRPCError(code int32,message string,id int)[]byte{
	json:=simplejson.New()
	error:=simplejson.New()
	error.Set("code",code)
	error.Set("message",message)
	json.Set("error",error)
	json.Set("id",id)
	b,_:=json.Encode()
	return b
}

func GetRPCLLError(err error,id int)[]byte{
	json:=simplejson.New()
	error:=simplejson.New()
	llerr,ok:=err.(*LLError)
	var errcode int32
	if(ok) {
		errcode=llerr.Code()
	}
	error.Set("code", errcode)
	error.Set("message",err.Error())
	json.Set("error",error)
	json.Set("id",id)
	b,_:=json.Encode()
	return b
}
type RPCparam struct{
	Key string
	Value interface{}
}
func GetRPC(method string,id int,RPCparams  ...*RPCparam)[]byte{
	json:=simplejson.New()

	json.Set("method",method)
	json.Set("id",id)

	params := simplejson.New()
	for _,rpcparem:=range RPCparams{
		if(rpcparem!=nil) {
			params.Set(rpcparem.Key, rpcparem.Value)
		}
	}
	json.Set("params",params)
	b,_:=json.Encode()
	return b
}

func GetRPCfromJson(method string,id int,params interface{})[]byte{
	json:=simplejson.New()

	json.Set("method",method)
	json.Set("id",id)

	json.Set("params",params)
	b,_:=json.Encode()
	return b
}