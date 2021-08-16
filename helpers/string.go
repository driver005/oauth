package helper

import (
	"strings"
	"unicode"
)

// ToLowerInitial converts a string's first character to lower case.
func ToLowerInitial(s string) string {
	if s == "" {
		return ""
	}
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

// ToUpperInitial converts a string's first character to upper case.
func ToUpperInitial(s string) string {
	if s == "" {
		return ""
	}
	a := []rune(s)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}

// Splitx is a special case of strings.Split
// which returns an empty slice if the string is empty
func Splitx(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	return strings.Split(s, sep)
}

// Coalesce returns the first non-empty string value
func Coalesce(str ...string) string {
	for _, s := range str {
		if s != "" {
			return s
		}
	}
	return ""
}
