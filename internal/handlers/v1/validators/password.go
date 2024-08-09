package validators

import "regexp"

func IsValidPassword(password string) string {
	var invalid = ""

	if len(password) < 8 {
		invalid = "Password must be at least 8 characters long"
	}

	if len(password) > 64 {
		invalid = "Password must be at most 64 characters long"
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		invalid = "Password must contain at least one uppercase letter"
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		invalid = "Password must contain at least one lowercase letter"
	}

	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		invalid = "Password must contain at least one digit"
	}

	return invalid
}
