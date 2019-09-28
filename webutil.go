package botil

import (
	"github.com/bitly/go-simplejson"
	"github.com/urfave/negroni"
	"net/http"
)


type ApiDetail struct {
	Function func(http.ResponseWriter, *http.Request)
	ApiPath  string
	Method   []string
	Middlewares  []negroni.Handler
}

func HttpJsonRespond(w http.ResponseWriter,json *simplejson.Json ,httpcode int,errcode int32,errmsg string){
	json.Set("code",httpcode)
	switch httpcode {
	case 200:
		json.Set("status","ok")
	case 400:
		json.Set("status","badrequest")
	case 500:
		json.Set("status","servererror")
	}
	w.Header().Set("content-type","application/json; charset=utf-8")
	if(httpcode!=200){
		w.WriteHeader(httpcode)
		json.Set("errorCode",errcode)
		json.Set("errorMessage",errmsg)
	}
	b,_:=json.Encode()
	w.Write(b)
}

func HttpLLErrorRespond(w http.ResponseWriter ,httpcode int,err error){
	json:=simplejson.New()
	llerr,ok:=err.(*LLError)
	var errcode int32 = 1
	if(ok) {
		errcode=llerr.Code()
	}
	json.Set("code",httpcode)
	switch httpcode {
	//200
	case  http.StatusOK:
		json.Set("status","ok")
		//400
	case  http.StatusBadRequest:
		json.Set("status","badrequest")
	case http.StatusForbidden:
		json.Set("status","forbidden")
		//500
	case  http.StatusInternalServerError:
		json.Set("status","servererror")
	}
	w.Header().Set("content-type","application/json; charset=utf-8")
	if(httpcode!=200){
		w.WriteHeader(httpcode)
		json.Set("errorCode",errcode)
		json.Set("errorMessage",err.Error())
	}
	b,_:=json.Encode()
	w.Write(b)

}