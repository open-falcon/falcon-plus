package str

import (
	"regexp"
	"strings"
)

func IsMatch(val, pattern string) bool {
	match, err := regexp.Match(pattern, []byte(val))
	if err != nil {
		return false
	}

	return match
}

func IsEnglishIdentifier(val string, pattern ...string) bool {
	defpattern := "^[a-zA-Z0-9\\-\\_\\.]+$"
	if len(pattern) > 0 {
		defpattern = pattern[0]
	}

	return IsMatch(val, defpattern)
}

func IsMail(val string) bool {
	return IsMatch(val, `\w[-._\w]*@\w[-._\w]*\.\w+`)
}

func IsPhone(val string) bool {
	if strings.HasPrefix(val, "+") {
		return IsMatch(val[1:], `\d{13}`)
	} else {
		return IsMatch(val, `\d{11}`)
	}
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

func Dangerous(str string) bool {
	return HasDangerousCharacters(str)
}
