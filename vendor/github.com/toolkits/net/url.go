package net

import (
	"net/url"
)

// url encode string, is + not %20
func UrlEncode(str string) string {
	return url.QueryEscape(str)
}

// url decode string
func UrlDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}
