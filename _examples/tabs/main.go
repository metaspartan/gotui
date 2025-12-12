package main

import (
	"image"
	"log"
	"time"

	ui "github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	header := widgets.NewParagraph()
	header.Text = "gotui Tabs Example"
	header.SetRect(50, 3, 75, 4)
	header.Border = false
	header.TextStyle = ui.NewStyle(ui.ColorWhite, ui.ColorClear, ui.ModifierBold)
	header.TextAlignment = ui.AlignRight

	tabpane := widgets.NewTabPane("Tab 1", "Tab 2", "Tab 3", "Tab 4", "Tab 5")
	tabpane.SetRect(0, 3, 50, 4)
	tabpane.Border = false

	tabpane.ActiveTabStyle = ui.NewStyle(ui.ColorWhite, ui.ColorGold)
	tabpane.InactiveTabStyle = ui.NewStyle(ui.ColorWhite, ui.ColorGreen)

	tabpane.PadLeft = 1
	tabpane.PadRight = 1
	tabpane.TabGap = 1
	tabpane.Separator = ""

	renderTab := func() *widgets.Paragraph {
		p := widgets.NewParagraph()
		p.SetRect(0, 4, 75, 12)
		p.Border = true
		p.BorderType = ui.BorderLine
		p.BorderStyle.Fg = ui.ColorGold

		switch tabpane.ActiveTabIndex {
		case 0:
			p.Text = "Hello, World!"
			p.BorderRounded = true
			tabpane.ActiveTabStyle = ui.NewStyle(ui.ColorWhite, ui.ColorBlue)
			p.BorderStyle.Fg = ui.ColorGold
		case 1:
			p.Text = "Welcome to the gotui tabs example!"
			p.BorderType = ui.BorderThick
			tabpane.ActiveTabStyle = ui.NewStyle(ui.ColorWhite, ui.ColorPink)
			p.BorderStyle.Fg = ui.ColorGreen
		case 2:
			p.Text = "Look! I'm different than others!"
			p.BorderType = ui.BorderDouble
			tabpane.ActiveTabStyle = ui.NewStyle(ui.ColorWhite, ui.ColorMagenta)
			p.BorderStyle.Fg = ui.ColorMagenta
		case 3:
			p.Text = "I know, these are some basic changes. But I think you got the main idea."
			p.BorderType = ui.BorderLine
			tabpane.ActiveTabStyle = ui.NewStyle(ui.ColorWhite, ui.ColorRed)
			p.BorderStyle.Fg = ui.ColorRed
		case 4:
			p.Text = "Custom Block Borders! (█ ▀ ▄ ▌ ▐)"
			p.BorderType = ui.BorderBlock
			tabpane.ActiveTabStyle = ui.NewStyle(ui.ColorBlack, ui.ColorGold)
			p.BorderStyle.Fg = ui.ColorGold
		}
		return p
	}

	footer := widgets.NewParagraph()
	footer.Text = "◄ ► to change tab | Press q to quit"
	footer.SetRect(0, 12, 75, 13)
	footer.Border = false
	footer.TextStyle = ui.NewStyle(ui.ColorWhite)
	footer.TitleAlignment = ui.AlignCenter

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Millisecond * 100)

	ui.Clear()
	content := renderTab()
	ui.Render(header, tabpane, content, footer)

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "h", "<Left>":
				tabpane.FocusLeft()
			case "l", "<Right>":
				tabpane.FocusRight()
			case "<MouseLeft>":
				payload := e.Payload.(ui.Mouse)
				x, y := payload.X, payload.Y
				idx := tabpane.ResolveClick(image.Pt(x, y))
				if idx != -1 {
					tabpane.ActiveTabIndex = idx
				}
			}
			ui.Clear()
			content := renderTab()
			ui.Render(header, tabpane, content, footer)

		case <-ticker.C:
			ui.Clear()
			content := renderTab()
			ui.Render(header, tabpane, content, footer)
		}
	}
}
