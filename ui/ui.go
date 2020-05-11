package ui

import (
	"github.com/blasrodri/frown/stats"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"sort"
	"strings"
	"strconv"
)

type UIConfig struct {
	FilterSecurityLevel int
}

func TerminalUI(config *UIConfig, reportChan <-chan *stats.Report, closeChan chan<- bool) {
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
		report := <-reportChan
		rows := make([][]string, 0)
		for processName, mapConnRep := range report.ProcessInfo {
			for pid, socketIdMapConnReport := range mapConnRep {
				for _, connReport := range socketIdMapConnReport {
					if connReport.SecurityLevel >= config.FilterSecurityLevel {
						rows = append(rows, []string{processName, strconv.Itoa(pid), connReport.DomainName, connReport.AdditionalInfo, strconv.Itoa(connReport.SecurityLevel)})
					}
				}
			}
		}
		table1.Rows = make([][]string, len(rows) + 1)
		table1.Rows[0] = []string{"Process", "PID", "Domain", "Additional Info"}

		sort.Slice(rows, func(i, j int) bool {
			switch(strings.Compare(rows[i][0], rows[j][0])) {
				case -1:
				return true
				case 0:
				return strings.Compare(rows[i][2], rows[j][2]) == -1
				default:
				return false
			}
		})

		for i, _ := range rows {
			idx := 1 + i
			table1.RowStyles[idx] = ui.NewStyle(ui.ColorWhite)
		}

		for i, row := range rows {
			idx := 1 + i
			table1.Rows[idx] = row[:len(row) - 1]
			secLevel, _ := strconv.Atoi(row[len(row) - 1])
			switch secLevel {
			case 0:
			case 1:
			case 2:
				table1.RowStyles[idx] = ui.NewStyle(ui.ColorWhite, ui.ColorYellow)
				continue
			case 3:
				table1.RowStyles[idx] = ui.NewStyle(ui.ColorWhite, ui.ColorRed, ui.ModifierBold)
				continue
			default:
			}
		}
		table1.SetRect(0, 3, termWidth, termHeight)
		ui.Render(title, table1)
		go func() {
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
