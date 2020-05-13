package main

import (
	"github.com/blasrodri/frown/ui"
)

func main() {

	var uiConfig = &ui.UIConfig{
		FilterSecurityLevel: 1,
	}
	manageState(uiConfig, ui.TerminalUI)
}
