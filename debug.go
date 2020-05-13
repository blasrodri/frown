package main

import (
	"fmt"
	"github.com/blasrodri/frown/stats"
	"github.com/blasrodri/frown/ui"
	"strconv"
)

func debug(config *ui.UIConfig, reportChan <-chan *stats.Report, closeChan chan<- bool) {
	for {

		report := <-reportChan
		rows := make([][]string, 0)
		for processName, mapConnRep := range report.ProcessInfo {
			for pid, socketIdMapconnReport := range mapConnRep {
				for socketId, connReport := range socketIdMapconnReport {
					rows = append(rows, []string{processName, strconv.Itoa(pid), socketId, connReport.DomainName})
				}
			}
		}
		fmt.Printf("%+v\n", rows)
	}
}
