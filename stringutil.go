package botil

import "strings"

func StringInSlice(a string, list []string,ignorecase bool) bool {
	for _, b := range list {
		if ignorecase{
			if strings.EqualFold(a,b){
				return true
			}
		}else {
			if b == a {
				return true
			}
		}
	}
	return false
}
