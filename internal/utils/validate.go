package utils

import "regexp"

func ValidatePasswordRoom(roomPassword string) bool {
	regex := `^\d{6}$`
	passwordLen := len(roomPassword)

	match, _ := regexp.MatchString(regex, roomPassword)

	if passwordLen == 6 && match {
		return true
	}

	return false
}
