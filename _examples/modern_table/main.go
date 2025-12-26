package main

import (
	"log"

	ui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := widgets.NewTable()
	table.Title = "Modern Table Demo"
	table.Rows = [][]string{
		{"Name", "Address", "Email"},
		{"Alaina Batz", "18312 Lake Alabama bury\nBrettfort, Tennessee 04176", "elvisryan@predovic.net"},
		{"Caleigh Conn", "76367 Lake Illinois furt\nDessiebury, Pennsylvania 74574", "jettstoltenberg@damore.name"},
		{"Christelle Jewess", "6157 Lake South Carolina fort\nKeeblerborough, North Carolina 72673", "queeniegleichner@kunze.net"},
		{"Deron Schmeler", "88886 South Hawaii mouth\nAlicialand, Mississippi 71041", "justonmarquardt@tremblay.biz"},
		{"Ezra Koss", "4421 Port Arkansas chester\nHilllbury, Arizona 34088", "aliciamurazik@kris.io"},
	}
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	// Full Width Setup
	termWidth, termHeight := ui.TerminalDimensions()
	table.SetRect(0, 0, termWidth, termHeight)

	// Modern Features
	table.ShowCursor = true
	table.CursorColor = ui.ColorBlue
	table.ShowLocation = true

	// Zebra Striping & Header Style
	table.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorClear, ui.ModifierBold) // Header
	table.RowStyles[2] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack)                  // Stripe
	table.RowStyles[4] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack)                  // Stripe

	// Selection Style (Blue background like the mockup)
	table.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorBlue)
	table.SelectedRow = 1 // Start on first data row

	ui.Render(table)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			table.SetRect(0, 0, payload.Width, payload.Height)
			ui.Clear()
			ui.Render(table)
		case "<Up>":
			table.ScrollUp()
		case "<Down>":
			table.ScrollDown()
		case "<Home>":
			table.ScrollTop()
		case "<End>":
			table.ScrollBottom()
		}
		ui.Render(table)
	}
}
