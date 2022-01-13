package prompt

import (
	"errors"
	"regexp"

	"github.com/fatih/color"
)

// GetConfirmationPrompt returns a Prompt that has been configured
// to show a [y/n] prompt to the user with `s` as label and `i` as info lines.
func GetConfirmationPrompt(s string, i []string) *Prompt {
	return &Prompt{
		InfoLines:   i,
		Label:       s,
		LabelColour: color.FgYellow,
		Invalid:     "Please enter [y/n]",
		Validator: func(s string) error {
			re, _ := regexp.Compile(`[yYnN]`)
			if !re.MatchString(s) {
				return errors.New("")
			}

			return nil
		},
		Retries: 3,
	}
}
