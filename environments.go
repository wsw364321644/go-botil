package botil

import (
	"runtime"
	"strings"
)

var Platform=[]string{
	"ps4",
	"pc",
	"xbox",
}

var EnvironmentType=[]string{
	"Development",
	"Certification",
	"Production",
	"Unknown",
}

func IsWindows()bool{
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}
func GetPathSep()byte{
	if(IsWindows()){
		return ';'
	}else{
		return ':'
	}
}
func AddEnv(env []string,key string,values []string)[]string{
	index:=-1
	var finalvalue string
	for i,str := range env{
		str=strings.Trim(str," ")
		if(!strings.EqualFold(str[:len(key)], key)){
			continue
		}
		index=i
		finalvalue=str
		for _,value:=range values{
			if(!(str[len(str)-1]==GetPathSep())){
				finalvalue+=string(GetPathSep())
			}
			finalvalue+=value
		}
		break
	}
	env[index]=finalvalue
	return env
}