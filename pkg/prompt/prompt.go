package prompt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MaikelVeen/branch/pkg/printer"
	"github.com/fatih/color"
	"golang.org/x/term"
)

// Prompt represents a single line input field with options for
// validation and input masks.
type Prompt struct {
	// InfoLines that are printed before the.
	InfoLines []string
	// The label of the prompt.
	Label string
	// The color of the label
	LabelColour color.Attribute
	// HideEntered will mask the input is true.
	HideEntered bool
	// The validation function for the input
	Validator ValidateFunc
	// Invalid represent the string to show if the input is invalid.
	Invalid string
	// The number of maximum allowed retries. If set to zero it will be automatically
	// assigned 999.
	Retries int
}

// Run executes the prompt. Its prints the informational lines first and then shows the user the prompt.
// The execution will stay alive until the max retries has been reached.
// It returns the value and an error if any occurred during the prompt's execution.
func (p *Prompt) Run() (string, error) {
	// If no Retries has been set, set to high value
	if p.Retries == 0 {
		p.Retries = 999
	}

	// Print all the info lines.
	for _, line := range p.InfoLines {
		printer.Print(line)
	}

	// Only print a new line when info lines were printed.
	if len(p.InfoLines) > 0 {
		printer.NewLine()
	}

	fd := os.Stdin.Fd()
	reader := bufio.NewReader(os.Stdin)

	// Loop until retries are zero.
	for ; p.Retries > 0; p.Retries-- {
		printLabel(p.Label, p.LabelColour)

		input, err := getInput(reader, int(fd), p.HideEntered)
		if err != nil {
			log.Fatal(err) // TODO: use printer
		}

		// Empty input (i.e. "\n")
		if len(input) < 2 {
			continue
		}

		result := strings.TrimSuffix(input, "\n")

		if p.Validator != nil {
			err := p.Validator(result)

			if err != nil {
				printer.Warning(p.Invalid)
				continue
			}
		}

		if p.HideEntered {
			fmt.Println()
		}
		return result, nil
	}

	return "", nil
}

func printLabel(label string, c color.Attribute) {
	p := color.New(c)
	p.Printf("%s: ", label)
}

// getInput retrieves a string from the standard input.
func getInput(rd *bufio.Reader, fileDescriptor int, hide bool) (string, error) {
	if !hide {
		return rd.ReadString('\n')
	}

	b, err := term.ReadPassword(fileDescriptor)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
