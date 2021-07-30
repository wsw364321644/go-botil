package botil

import "strings"

func StringInSlice(a string, list []string,ignorecase bool) bool {
	for _, b := range list {
		if ignorecase{
			return strings.EqualFold(a,b)
		}else {
			if b == a {
				return true
			}
		}
	}
	return false
}
