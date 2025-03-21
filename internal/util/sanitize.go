package util

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func RemoveAccents(value string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, value)
	return result
}

var (
	nonAlphaRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)
)

func RemoveNonAlpha(value string) string {
	return nonAlphaRegex.ReplaceAllString(value, "-")
}

func SanitizeAlphaLower(value string) string {
	value = RemoveAccents(value)
	value = RemoveNonAlpha(value)
	value = strings.ToLower(value)
	value = strings.Trim(value, "-")
	return value
}
