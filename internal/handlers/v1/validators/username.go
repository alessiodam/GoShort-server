package validators

import "regexp"

func IsValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._]+$`)
	return re.MatchString(username)
}
