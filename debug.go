package main

import (
	"fmt"
	"strconv"

	"github.com/blasrodri/frown/stats"
	"github.com/blasrodri/frown/ui"
)

func debug(config *ui.UIConfig, reportChan <-chan *stats.Report, closeChan chan<- bool) {
	for {

		report := <-reportChan
		rows := make([][]string, 0)
		for processName, mapConnRep := range report.ProcessInfo {
			for pid, SocketIDMapconnReport := range mapConnRep {
				for SocketID, connReport := range SocketIDMapconnReport {
					rows = append(rows, []string{processName, strconv.Itoa(pid), SocketID, connReport.DomainName})
				}
			}
		}
		fmt.Printf("%+v\n", rows)
	}
}
