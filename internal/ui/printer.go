package ui

import "github.com/pterm/pterm"

func Info(msg string) {
	pterm.Info.Println(msg)
}

func Success(msg string) {
	pterm.Success.Println(msg)
}

func Warning(msg string) {
	pterm.Warning.Println(msg)
}

func Error(msg string) {
	pterm.Error.Println(msg)
}
