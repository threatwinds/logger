package rerror

import "strings"

func Is(err error, kargs ...string) bool {
	for _, e := range kargs {
		if strings.Contains(err.Error(), e) {
			return true
		}
	}
	return false
}
