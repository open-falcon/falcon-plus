package utils

import (
	"regexp"
	"strings"
)

func IsUsernameValid(name string) bool {
	match, err := regexp.Match("^[a-zA-Z0-9\\-\\_\\.]+$", []byte(name))
	if err != nil {
		return false
	}

	return match
}

func HasDangerousCharacters(str string) bool {
	if strings.Contains(str, "<") {
		return true
	}

	if strings.Contains(str, ">") {
		return true
	}

	if strings.Contains(str, "&") {
		return true
	}

	if strings.Contains(str, "'") {
		return true
	}

	if strings.Contains(str, "\"") {
		return true
	}

	return false
}
