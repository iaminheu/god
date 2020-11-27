package util

import "strings"

// 首字母大写
func Title(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

// 首字母小写
func UnTitle(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToLower(s[:1]) + s[1:]
}
