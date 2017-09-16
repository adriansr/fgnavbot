package main

import (
	"fmt"

	"github.com/adriansr/fgnavbot/nav"
)

func main() {
	channel := make(chan interface{})
	go nav.ReadNavaids("nav.dat.gz", channel)
	go nav.ReadAirports("apt.dat.gz", channel)

	lastApt := ""
	for terminated := 0; terminated < 2; {
		switch item := (<-channel).(type) {
		case *nav.Navaid:
			fmt.Printf("nav %s at %v\n", item.Identifier, item.Pos)
		case *nav.Airport:
			fmt.Printf("airport %s\n", item.Code)
			lastApt = item.Code
		case *nav.Runway:
			fmt.Printf("runway %s %s-%s\n", lastApt,
				item.End[0].Code, item.End[1].Code)
		case error:
			fmt.Printf("Error %v\n", item)
			terminated++
		case nav.Terminator:
			terminated++
		}
	}
}
