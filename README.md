# gotui

gotui is a cross-platform and fully-customizable terminal dashboard and widget library built on top of [tcell](https://github.com/gdamore/tcell). It is a modern fork of `termui`, inspired by [blessed-contrib](https://github.com/yaronn/blessed-contrib), [tui-rs](https://github.com/fdehau/tui-rs), and [ratatui](https://github.com/ratatui-org/ratatui) and written purely in Go by Carsen Klock.

## Note

This is a modern fork of termui for 2025, heavily upgraded to support TrueColor, modern terminal events, and new layouts.

## Versions

gotui is compatible with Go 1.24+.

## Features

- **Backend**: Native `tcell` support for TrueColor (24-bit RGB), mouse events, and resize handling.
- **Widgets**:
  - **Charts**: BarChart, StackedBarChart, PieChart, Plot (Line/Scatter), Sparkline, Gauge, Heatmap.
  - **Text**: Paragraph (with wrapping and alignment), List, Tree.
  - **Layout**: Grid system, Flexbox-like blocks, Tabs.
  - **Interactive**: Calendar, Tables, Input, TextArea.
- **Styling**:
  - Full RGB Color support.
  - Border titles (Top and Bottom) with alignment (Left, Center, Right).
  - Rich styling parser for text.
- **Compatibility**: Works with modern terminals (iTerm2, Kitty, Alacritty, Ghostty).

## Installation

### Go modules

It is not necessary to `go get` gotui, since Go will automatically manage any imported dependencies for you.

```bash
go get github.com/metaspartan/gotui
```

## Hello World

```go
package main

import (
	"log"

	ui "github.com/metaspartan/gotui"
	"github.com/metaspartan/gotui/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Text = "Hello World!"
	p.SetRect(0, 0, 25, 5)

	ui.Render(p)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}
}
```

## Widgets

- [BarChart](./_examples/barchart.go)
- [Calendar](./_examples/calendar.go)
- [Canvas](./_examples/canvas.go) (for drawing braille dots)
- [Gauge](./_examples/gauge.go)
- [Heatmap](./_examples/heatmap.go)
- [Image](./_examples/image.go)
- [Input](./_examples/input.go)
- [List](./_examples/list.go)
- [Tree](./_examples/tree.go)
- [Paragraph](./_examples/paragraph.go)
- [PieChart](./_examples/piechart.go)
- [Plot](./_examples/plot.go) (for scatterplots and linecharts)
- [Sparkline](./_examples/sparkline.go)
- [StackedBarChart](./_examples/stacked_barchart.go)
- [Table](./_examples/table.go)
- [Tabs](./_examples/tabs.go)
- [TextArea](./_examples/textarea.go)

Run an example with `go run _examples/{example}.go` or run each example consecutively with `make run-examples`.

## Uses

- [mactop](https://github.com/context-labs/mactop)

(Submit your projects via a PR)

## Acknowledgments

- [termui](https://github.com/gizak/termui)

## Author(s)

gotui Author: Carsen Klock - [X](https://x.com/carsenklock)

termui Author: Zack Guo - [Github](https://github.com/gizak)

## Related Works

- [blessed-contrib](https://github.com/yaronn/blessed-contrib)
- [gocui](https://github.com/jroimartin/gocui)
- [termdash](https://github.com/mum4k/termdash)
- [tui-rs](https://github.com/fdehau/tui-rs)
- [tview](https://github.com/rivo/tview)
- [termui](https://github.com/gizak/termui)
- [ratatui](https://github.com/ratatui-org/ratatui)

## License

[MIT](http://opensource.org/licenses/MIT)
