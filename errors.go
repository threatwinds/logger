package rerror

import "strings"

const (
	GORM_DUPLICATE_KEY string = "duplicate key value violates unique constraint"
)

func Is(err error, kargs ...string) bool {
	for _, e := range kargs {
		if strings.Contains(err.Error(), e) {
			return true
		}
	}
	return false
}
