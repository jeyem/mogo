package mogo

import (
	"strings"
	"unicode"
)

func toSnake(s string) string {
	var (
		res  string
		last int
	)
	ls := []rune(s)
	for i, char := range ls {
		if (i == 0 || !unicode.IsUpper(char)) && i+1 != len(s) {
			continue
		}
		if i+1 != len(s) {
			res += strings.ToLower(s[last:i]) + "_"
		} else {
			res += strings.ToLower(s[last : i+1])
		}
		last = i
	}
	return res
}
