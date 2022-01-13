package git

import (
	"fmt"
	"regexp"
	"strings"
)

// filter will apply the following filter to `s`:
// It will remove trailing whitespace.
// It will make all chars lowercase.
// It will first remove any values that are between () and [].
// It it will filter all the special characters.
func filter(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	// Remove everything inside () and [] to remove tags
	innerValRe, _ := regexp.Compile(`([\(\[]).*?([\)\]])`)
	s = innerValRe.ReplaceAllString(s, "")

	// Remove all non alpha numeric chars expect whitespace
	special, _ := regexp.Compile(`[^a-zA-Z\d\s:]`)
	s = special.ReplaceAllString(s, "")

	// Split the string on the spaces.
	parts := strings.Split(s, " ")
	cutoff := getMax(parts)

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
