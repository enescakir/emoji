package strutil

import (
	"regexp"
	"strings"
)

var (
	nonAlphaNumRegex = regexp.MustCompile(`[^\w\d]+`)
	whitespaceRegex  = regexp.MustCompile(`\s+`)

	changes = map[string]string{
		"*":    "asterisk",
		"#":    "hash",
		"1st":  "first",
		"2nd":  "second",
		"3rd":  "third",
		"&":    "and",
		"U.S.": "US",
		"Š":    "S",
		"š":    "s",
		"Đ":    "Dj",
		"đ":    "dj",
		"Ž":    "Z",
		"ž":    "z",
		"Č":    "C",
		"č":    "c",
		"Ć":    "C",
		"ć":    "c",
		"À":    "A",
		"Á":    "A",
		"Â":    "A",
		"Ã":    "A",
		"Ä":    "A",
		"Å":    "A",
		"Æ":    "A",
		"Ç":    "C",
		"È":    "E",
		"É":    "E",
		"Ê":    "E",
		"Ë":    "E",
		"Ì":    "I",
		"Í":    "I",
		"Î":    "I",
		"Ï":    "I",
		"Ñ":    "N",
		"Ò":    "O",
		"Ó":    "O",
		"Ô":    "O",
		"Õ":    "O",
		"Ö":    "O",
		"Ø":    "O",
		"Ù":    "U",
		"Ú":    "U",
		"Û":    "U",
		"Ü":    "U",
		"Ý":    "Y",
		"Þ":    "B",
		"ß":    "Ss",
		"à":    "a",
		"á":    "a",
		"â":    "a",
		"ã":    "a",
		"ä":    "a",
		"å":    "a",
		"æ":    "a",
		"ç":    "c",
		"è":    "e",
		"é":    "e",
		"ê":    "e",
		"ë":    "e",
		"ì":    "i",
		"í":    "i",
		"î":    "i",
		"ï":    "i",
		"ð":    "o",
		"ñ":    "n",
		"ò":    "o",
		"ó":    "o",
		"ô":    "o",
		"õ":    "o",
		"ö":    "o",
		"ø":    "o",
		"ù":    "u",
		"ú":    "u",
		"û":    "u",
		"ý":    "y",
		"þ":    "b",
		"ÿ":    "y",
		"Ŕ":    "R",
		"ŕ":    "r",
	}
)

func Clean(str string) string {
	for o, n := range changes {
		str = strings.ReplaceAll(str, o, n)
	}

	str = nonAlphaNumRegex.ReplaceAllString(str, " ")

	return str
}

func RemoveSpaces(str string) string {
	return whitespaceRegex.ReplaceAllString(str, "")
}
