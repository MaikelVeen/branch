package prompt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MaikelVeen/branch/printer"
	"github.com/fatih/color"
)

type Prompt struct {
	InfoLines   []string
	Label       string
	LabelColor  color.Attribute
	HideEntered bool
	Validator   ValidateFunc

	Valid   string
	Invalid string

	Retries int
}

func (p *Prompt) Run() (string, error) {
	// If no Retries has been set, set to high value
	if p.Retries == 0 {
		p.Retries = 999
	}

	for _, line := range p.InfoLines {
		printer.Print(line)
	}
	printer.NewLine()

	r := bufio.NewReader(os.Stdin)

	for ; p.Retries > 0; p.Retries-- {
		fmt.Printf("%s: ", p.Label)

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}

		r := strings.TrimSuffix(res, "\n")

		if p.Validator != nil {
			err := p.Validator(r)

			if err != nil {
				printer.Warning(p.Invalid)
				continue
			}
		}

		return res, nil
	}

	return "", nil
}
