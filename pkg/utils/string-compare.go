package utils

import (
	"strings"
)

const CutSet = " \t\n\r"
const NotPrefix = "!"

func CaseInsensitiveCompareStrings(a string, b string) bool {
	a = strings.Trim(strings.ToLower(a), CutSet)
	b = strings.Trim(strings.ToLower(b), CutSet)
	return a == b
}

func PrefixedCaseInsensitiveCompareStrings(a string, b string) bool {
	b, invert := strings.CutPrefix(b, NotPrefix)
	return CaseInsensitiveCompareStrings(a, b) != invert
}
