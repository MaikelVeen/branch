package logger

import "github.com/fatih/color"

func Error(err error, msg string, debug bool) {
	color.Red(msg)
}

func Success(msg string) {
	color.Green(msg)
}

func Warning(msg string) {
	color.Yellow(msg)
}
