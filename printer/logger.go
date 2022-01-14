package printer

// TODO: function documentation

import (
	"fmt"

	"github.com/fatih/color"
)

var greetColor = color.New(color.FgHiGreen, color.Bold, color.Underline)

func Error(msg *string, err error) {
	var message string

	if msg == nil {
		message = "Error encountered"
	} else {
		message = *msg
	}

	color.Red("%s: %s", message, err)
}

func Success(msg string) {
	color.Green(msg)
}

func Warning(msg string) {
	color.Yellow(msg)
}

func Print(msg string) {
	color.White(msg)
}

func Greet(msg string) {
	greetColor.Println(msg)
}

func NewLine() {
	fmt.Println()
}

func Fatal(err string) {

}
