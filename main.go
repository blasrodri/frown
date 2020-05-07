package main

import (
	"fmt"
	"github.com/blasrodri/frown/lsof"
	"log"
	"time"
)

func main() {
	for {
		time.Sleep(500 * time.Duration(time.Millisecond))

		userPids, err := lsof.GetUserProcessList()
		if err != nil {
			log.Fatal(err)
		}
		for _, uProc := range userPids {
			fmt.Printf("%+v\n", *uProc)
		}
		connDetails, err := lsof.MonitorUserConnections()
		if err != nil {
			log.Fatal(err)
		}
		for _, connDeets := range connDetails {
			fmt.Printf("%+v\n", *connDeets)
		}
	}
}
