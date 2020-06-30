package main

import (
	"fmt"
	"github.com/blasrodri/frown/ui"
	"os"
	"runtime"
)

func main() {
	if !isSupportedOs() {
		fmt.Fprintln(os.Stderr, "Frown is only available for Linux. Sorry!")
		os.Exit(-1)
	}
	var uiConfig = &ui.UIConfig{
		FilterSecurityLevel: 1,
	}
	manageState(uiConfig, ui.TerminalUI)
}

func isSupportedOs() bool {
	os := runtime.GOOS
	switch os {
	case "linux":
		return true
	default:
		return false
	}
}
