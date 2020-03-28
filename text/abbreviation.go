package text

import "strings"

func Abbreviate(text string) string {
	abbreviation := ""
	tokens := strings.Split(text, "_")
	for _, token := range tokens {
		cleanToken := strings.TrimSpace(token)
		if len(cleanToken) > 0 {
			abbreviation = abbreviation + cleanToken[0:1]
		}
	}
	return abbreviation
}
