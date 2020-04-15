package botil

import (
	"github.com/bitly/go-simplejson"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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

func HttpJsonStrRespond(w http.ResponseWriter,rootstr string ,httpcode int,err error){
	if(!gjson.Parse(rootstr).IsArray()) {
		rootstr, _ = sjson.Set(rootstr, "code", httpcode)
		llerr, ok := err.(*LLError)
		var errcode int32 = 1
		var msg string
		if err != nil {
			if (ok) {
				errcode = llerr.Code()
				msg = llerr.Message()
			} else {
				msg = err.Error()
			}
		}
		switch httpcode {
		case 200:
			rootstr, _ = sjson.Set(rootstr, "status", "ok")
		case 400:
			rootstr, _ = sjson.Set(rootstr, "status", "badrequest")
		case 500:
			rootstr, _ = sjson.Set(rootstr, "status", "servererror")
		}
		w.Header().Set("content-type", "application/json; charset=utf-8")
		if (httpcode != 200) {
			w.WriteHeader(httpcode)
			rootstr, _ = sjson.Set(rootstr, "errorCode", errcode)
			rootstr, _ = sjson.Set(rootstr, "errorMessage", msg)
		}
	}
	w.Write([]byte(rootstr))
}

func HttpLLErrorRespond(w http.ResponseWriter ,httpcode int,err error){
	HttpJsonStrRespond(w,"{}",httpcode,err)
}