package rerror

import "strings"

func EIs(err error, kargs ...string) bool {
	for _, e := range kargs {
		if strings.Contains(err.Error(), e) {
			return true
		}
	}
	return false
}


func RIs(err *Error, kargs ...string) bool {
	for _, e := range kargs {
		if strings.Contains(err.Message, e) {
			return true
		}
	}
	return false
}
