package git

import (
	"fmt"
	"regexp"
	"strings"
)

// FormatAsValidRef transforms `s` to a string that is a valid refname
//
// A reference is used in git to specify branches and tags. The following rules
// must be followed when is comes to naming references:
// https://git-scm.com/docs/git-check-ref-format
func FormatAsValidRef(s string) string {
	// Remove everything inside () and [] to remove tags.
	innerValRe, _ := regexp.Compile(`([\(\[]).*?([\)\]])`)
	s = innerValRe.ReplaceAllString(s, "")

	// Remove all non alpha numeric chars expect whitespace.
	special, _ := regexp.Compile(`[^a-zA-Z\d\s]`)
	s = special.ReplaceAllString(s, " ")

	// Remove leading and trailing whitespace.
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	// Remove spaces with single space.
	d, _ := regexp.Compile(`\s\s+`)
	s = d.ReplaceAllString(s, " ")

	// Split the string on the spaces.
	parts := strings.Split(s, " ")
	cutoff := GetLengthWithUpperbound(parts, 12)

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

func GetBranchName(base, key, title string) string {
	//TODO: should we check key and base here to also valid for refname?
	ref := FormatAsValidRef(title)

	return fmt.Sprintf("%s/%s/%s", base, key, ref)
}
