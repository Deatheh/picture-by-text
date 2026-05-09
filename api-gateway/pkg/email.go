package pkg

import (
	"regexp"
)

func CheckIsEmailAllowed(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(email) {
		return false
	}
	return true
}
