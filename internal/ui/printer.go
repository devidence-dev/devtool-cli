package ui

import "github.com/pterm/pterm"

func init() {
	pterm.Info.Prefix = pterm.Prefix{
		Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack, pterm.Bold),
		Text:  "INFO",
	}
	pterm.Success.Prefix = pterm.Prefix{
		Style: pterm.NewStyle(pterm.BgGreen, pterm.FgBlack, pterm.Bold),
		Text:  "SUCCESS",
	}
	pterm.Warning.Prefix = pterm.Prefix{
		Style: pterm.NewStyle(pterm.BgYellow, pterm.FgBlack, pterm.Bold),
		Text:  "WARNING",
	}
	pterm.Error.Prefix = pterm.Prefix{
		Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite, pterm.Bold),
		Text:  "ERROR",
	}
}

func Info(msg string)    { pterm.Info.Println(msg) }
func Success(msg string) { pterm.Success.Println(msg) }
func Warning(msg string) { pterm.Warning.Println(msg) }
func Error(msg string)   { pterm.Error.Println(msg) }
