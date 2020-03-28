package text

import (
	"strings"
)

var mapping = map[string]string{
	"А": "A",
	"а": "a",
	"Б": "B",
	"б": "b",
	"В": "V",
	"в": "v",
	"Г": "H",
	"г": "h",
	"Ґ": "G",
	"ґ": "g",
	"Д": "D",
	"д": "d",
	"Е": "E",
	"е": "e",
	"Є": "Ye",
	"є": "ie",
	"Ж": "Zh",
	"ж": "zh",
	"З": "Z",
	"з": "z",
	"И": "Y",
	"и": "y",
	"І": "I",
	"і": "i",
	"Ї": "Yi",
	"ї": "i",
	"Й": "Y",
	"й": "i",
	"К": "K",
	"к": "k",
	"Л": "L",
	"л": "l",
	"М": "M",
	"м": "m",
	"Н": "N",
	"н": "n",
	"О": "O",
	"о": "o",
	"П": "P",
	"п": "p",
	"Р": "R",
	"р": "r",
	"С": "S",
	"с": "s",
	"Т": "T",
	"т": "t",
	"У": "U",
	"у": "u",
	"Ф": "F",
	"ф": "f",
	"Х": "Kh",
	"х": "kh",
	"Ц": "Ts",
	"ц": "ts",
	"Ч": "Ch",
	"ч": "ch",
	"Ш": "Sh",
	"ш": "sh",
	"Щ": "Shch",
	"щ": "shch",
	"Ю": "Yu",
	"ю": "iu",
	"Я": "Ya",
	"я": "ia",
	"ь": "",
	"'": "",
}

func TransliterateUkr(text string) string {
	// `зг` should be translated as `zgh` instead of `zh`.
	text = strings.ReplaceAll(text, "Зг", "Zgh")
	//works for small and big letters
	text = strings.ReplaceAll(text, "зг", "zgh")
	for key, value := range mapping {
		text = strings.ReplaceAll(text, key, value)
	}
	return text
}
