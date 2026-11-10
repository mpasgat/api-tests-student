package utils

import "strings"

// MixCase делает регистр букв вперемешку: "abcd" -> "AbCd"
func MixCase(s string) string {
	var sb strings.Builder
	for i, r := range []rune(s) {
		if i%2 == 0 {
			sb.WriteString(strings.ToUpper(string(r)))
		} else {
			sb.WriteString(strings.ToLower(string(r)))
		}
	}

	return sb.String()
}
