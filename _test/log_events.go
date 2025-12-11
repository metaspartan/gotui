//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"

	ui "github.com/metaspartan/gotui/v4"
)

// logs all events to the gotui window
// stdout can also be redirected to a file and read with `tail -f`
func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	events := ui.PollEvents()
	for {
		e := <-events
		fmt.Printf("%v", e)
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<MouseLeft>":
			return
		}
	}
}
