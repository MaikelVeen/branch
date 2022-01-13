package git

import (
	"fmt"
	"regexp"
	"strings"
)

// filter will clean `s` to be a valid part of a git branch
func filter(s string) string {
	// Remove everything inside () and [] to remove tags.
	innerValRe, _ := regexp.Compile(`([\(\[]).*?([\)\]])`)
	s = innerValRe.ReplaceAllString(s, "")

	// Remove all non alpha numeric chars expect whitespace.
	special, _ := regexp.Compile(`[^a-zA-Z\d\s:]`)
	s = special.ReplaceAllString(s, "")

	// Remove leading and trailing whitespace.
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	// Remove double spaces
	d, _ := regexp.Compile(`\s\s+`)
	s = d.ReplaceAllString(s, "")

	// Split the string on the spaces.
	parts := strings.Split(s, " ")
	cutoff := getMax(parts)

	// Return hyphenated string
	return strings.Join(parts[:cutoff], "-")
}

func getMax(p []string) int {
	l := len(p)

	if l > 12 {
		return 12
	}

	return l
}

func GetBranchName(key, base, title string) string {
	cleanTitle := filter(title)

	return fmt.Sprintf("%s/%s/%s", base, key, cleanTitle)
}
