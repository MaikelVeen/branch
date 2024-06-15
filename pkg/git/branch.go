package git

import (
	"regexp"
	"strings"
)

const (
	MaxParts = 12
)

// FormatAsValidRef transforms `s` to a string that is a valid refname
//
// A reference is used in git to specify branches and tags. The following rules
// must be followed when is comes to naming references:
// https://git-scm.com/docs/git-check-ref-format
func FormatAsValidRef(s string) string {
	// Remove everything inside () and [] to remove tags.
	innerValRe := regexp.MustCompile(`([\(\[]).*?([\)\]])`)
	s = innerValRe.ReplaceAllString(s, "")

	// Remove all non alpha numeric chars expect whitespace.
	special := regexp.MustCompile(`[^a-zA-Z\d\s]`)
	s = special.ReplaceAllString(s, " ")

	// Remove leading and trailing whitespace.
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	// Remove spaces with single space.
	d := regexp.MustCompile(`\s\s+`)
	s = d.ReplaceAllString(s, " ")

	// Split the string on the spaces.
	parts := strings.Split(s, " ")
	cutoff := GetLengthWithUpperbound(parts, MaxParts)

	// Return hyphenated string
	return strings.Join(parts[:cutoff], "-")
}

func GetLengthWithUpperbound(s []string, m int) int {
	l := len(s)

	if l > m {
		return m
	}

	return l
}
