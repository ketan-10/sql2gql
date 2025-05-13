package utils

import (
	"strings"
	"unicode"
)

func Filter[T any](list []T, filter func(item T) bool) []T {
	var newList []T
	for _, v := range list {
		if filter(v) {
			newList = append(newList, v)
		}
	}
	return newList
}

// SnakeToCamel converts snake_case to CamelCase.
func SnakeToCamel(name string) string {
	var result string
	for _, word := range strings.Split(name, "_") {
		if word == "" {
			continue
		}
		result += strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
	}
	return result
}

// lcFirst converts the first character of a string to lowercase.
func LowCaseFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
