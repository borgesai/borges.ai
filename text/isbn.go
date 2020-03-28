package text

import (
	"strconv"
	"strings"
)

func NormalizeISBN(text string) string {
	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(text, " ", ""), "-", ""));
}

func IsISBN(text string) bool {
	isValidLength := len(text) == 10 || len(text) == 13
	if !isValidLength {
		return false
	}
	// handle numbers like 000100039X
	if len(text) == 10 {
		text = strings.ToLower(text)
		text = strings.TrimRight(text, "x")
	}
	isbn, _ := strconv.ParseInt(text, 10, 64)
	return isbn > 0
}
