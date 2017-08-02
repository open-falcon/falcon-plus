package util

import "encoding/json"

func ParseCookieFromResp(body string) (string, string) {
	var paresToken map[string]string
	json.Unmarshal([]byte(body), &paresToken)
	sname, ssig := paresToken["name"], paresToken["sig"]
	return sname, ssig
}
