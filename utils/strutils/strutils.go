// strutils/strutils.go
package strutils

import "strings"

// Reverse возвращает обратную строку.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ToUpper возвращает строку в верхнем регистре.
func ToUpper(s string) string {
	return strings.ToUpper(s)
}