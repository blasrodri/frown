package ui

import (
	"log"
	"strconv"
	"github.com/blasrodri/frown/stats"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func TerminalUI(reportChan <- chan *stats.Report, closeChan chan <- bool) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	termWidth, termHeight := ui.TerminalDimensions()

	title := widgets.NewParagraph()
	title.Text = "Welcome to **Frown**"
	title.SetRect(0, 0, termWidth, 5)
	title.Border = false
	table1 := widgets.NewTable()
	for {
		report := <- reportChan
		table1.Rows = make([][]string, 0)
		for processName, mapConnRep := range report.ProcessInfo {
			for pid, socketIdMapConnReport := range mapConnRep {
				for _, connReport := range socketIdMapConnReport {
					table1.Rows = append(table1.Rows, []string{processName, strconv.Itoa(pid), connReport.DomainName})
				}
			}
		}
		table1.TextStyle = ui.NewStyle(ui.ColorWhite)
		table1.SetRect(0, 3, termWidth, termHeight)
		ui.Render(title, table1)
		go func (){
			uiEvents := ui.PollEvents()
			for {
				e := <-uiEvents
				switch e.ID {
				case "q", "<C-c>":
					closeChan <- true
					ui.Close()
				}
			}
		}()
	}
}
