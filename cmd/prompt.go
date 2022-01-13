package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// PromptConfirmation displays a prompt `s` to the user and returns a bool indicating yes / no
// If the lowercased, trimmed input begins with anything other than 'y', it returns false
// It accepts an int `tries` representing the number of attempts before returning false
func PromptConfirmation(s string, tries int) bool {
	r := bufio.NewReader(os.Stdin)

	for ; tries > 0; tries-- {
		fmt.Printf("%s [y/n]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}

		return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
	}

	return false
}

// Temp add here:
// StringInSliceCaseInsensitive checks whether a string exists in a slice.
func StringInSliceCaseInsensitive(str string, slice []string) bool {
	for _, sliceItem := range slice {
		if strings.EqualFold(str, sliceItem) {
			return true
		}
	}

	return false
}
