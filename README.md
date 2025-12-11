# gotui

[![Go Report Card](https://goreportcard.com/badge/github.com/metaspartan/gotui)](https://goreportcard.com/report/github.com/metaspartan/gotui/v4)
[![GoDoc](https://godoc.org/github.com/metaspartan/gotui?status.svg)](https://godoc.org/github.com/metaspartan/gotui/v4)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/metaspartan/gotui/blob/master/LICENSE)

gotui is a cross-platform and fully-customizable terminal dashboard and widget library built on top of [tcell](https://github.com/gdamore/tcell). It is a modern fork of [termui](https://github.com/gizak/termui), inspired by [ratatui](https://github.com/ratatui-org/ratatui) and written purely in Go by Carsen Klock.

![Logo](./logo.png)

## Note

This is a modern fork of termui for 2025, heavily upgraded to support TrueColor, modern terminal events, better performance, and new layouts.

## Versions

gotui is compatible with Go 1.24+.

## Features

- **Backend**: Native `tcell` support for TrueColor (24-bit RGB), mouse events, and resize handling.
- **Gauges**: Progress bars and gauges.
- **Charts**:
  - **BarChart**: Stacked and standard bar charts.
  - **PieChart**: Pie and Donut charts.
  - **RadarChart**: Spider/Radar charts.
  - **TreeMap**: Hierarchical data visualization.
  - **FunnelChart**: Process flow/conversion charts.
  - **Sparkline**: Mini sparklines.
  - **Plot**: Line, Scatter, and Braille-mode charts.
- **Maps**:
  - **World Map**: High-resolution world map example using the generic `Canvas` widget (see `_examples/canvas.go`).
- **New Widgets**:
  - **LineGauge**: Thin, character-based progress bar with alignment options (Block, Dots, custom runic styles).
  - **Scrollbar**: Ratatui-compatible scrollbars (Vertical/Horizontal) with mouse and keyboard support.
  - **Logo**: Pixel-perfect block-style logo renderer.
- **Performance**:
  - **Optimized Rendering**: `Buffer` uses flat slices for O(1) access, providing 2-3x speedup.
  - **Zero Allocations**: Drawing loops minimized for high-fps scenes (~3000 FPS potential).
- **Layout**:
  - **Grid**: Responsive grid layout.
  - **Tabs**: Tabbed navigation.
  - **Interactive**: Calendar, Tables, Input, TextArea.
- **Styling**:
  - **Rounded Borders**: Optional rounded corners for blocks.
  - Full RGB Color support.
  - Border titles (Top and Bottom) with alignment (Left, Center, Right).
  - Rich styling parser for text.
  - **Collapsed Borders**: Support for merging adjacent block borders using `BorderCollapse`.
- **Compatibility**: Works with modern terminals (iTerm2, Kitty, Alacritty, Ghostty).

## Installation

### Go modules

It is not necessary to `go get` gotui, since Go will automatically manage any imported dependencies for you.

```bash
go get github.com/metaspartan/gotui/v4
```

## Hello World

```go
package main

import (
	"log"

	ui "github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
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
- [Block](./_examples/block.go)
- [Collapsed Borders](./_examples/collapsed_borders.go)
- [Calendar](./_examples/calendar.go)
- [Canvas](./_examples/canvas.go)
- [Gauge](./_examples/gauge.go)
- [Heatmap](./_examples/heatmap.go)
- [Image](./_examples/image.go)
- [Input](./_examples/input.go)
- [List](./_examples/list.go)
- [Logo](./_examples/logo.go)
- [LineGauge](./_examples/linegauge.go)
- [Tree](./_examples/tree.go)
- [Paragraph](./_examples/paragraph.go)
- [PieChart](./_examples/piechart.go)
- [Plot](./_examples/plot.go) (for scatterplots and linecharts)
- [Sparkline](./_examples/sparkline.go)
- [StackedBarChart](./_examples/stacked_barchart.go)
- [Scrollbar](./_examples/scrollbar.go)
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
