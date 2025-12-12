<div align="center">
  <img src="./logo.png" width="300" alt="gotui Logo" />
  <h1>gotui</h1>
  <p>
    <strong>A modern, high-performance Terminal User Interface (TUI) library for Go.</strong>
  </p>

[![Go Report Card](https://goreportcard.com/badge/github.com/metaspartan/gotui)](https://goreportcard.com/report/github.com/metaspartan/gotui/v4)
[![GoDoc](https://godoc.org/github.com/metaspartan/gotui?status.svg)](https://pkg.go.dev/github.com/metaspartan/gotui/v4)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/metaspartan/gotui/blob/master/LICENSE)

</div>

**gotui** by Carsen Klock is a fully-customizable dashboard and widget Go library built on top of [tcell](https://github.com/gdamore/tcell). It is a modernized enhanced fork of [termui](https://github.com/gizak/termui), engineered for valid **TrueColor** support, **high-performance rendering**, flex layouts, rounded borders, input, and for feature parity with robust libraries like [ratatui](https://github.com/ratatui-org/ratatui).

---

![gotui](./gotui.gif)

## ⚡ Features

- **🚀 High Performance**: optimized rendering engine capable of **~3000 FPS** frame operations with zero-allocation drawing loops.
- **🎨 TrueColor Support**: Full 24-bit RGB color support for modern terminals (Ghostty, Alacritty, Kitty, iTerm2).
- **📐 Flexible Layouts**: 
  - **Flex**: Mixed fixed/proportional layouts.
  - **Grid**: 12-column dynamic grid system.
  - **Absolutes**: Exact coordinates when needed.
- **🌐 SSH / Remote Apps**: Turn any TUI into a zero-install SSH accessible application (multi-tenant support).
- **📊 Rich Widgets**:
  - **Charts**: BarChart, StackedBarChart, PieChart, DonutChart, RadarChart (Spider), FunnelChart, TreeMap, Sparkline, Plot (Scatter/Line).
  - **Gauges**: Gauge, LineGauge (with pixel-perfect Braille/Block styles).
  - **Interaction**: Input, TextArea, List, Table, Scrollbar.
  - **Misc**: TabPane, Image (block-based), Canvas (Braille), Heatmap.
- **🖱️ Mouse Support**: Full mouse event support (Click, Scroll Wheel, Drag).
- **🔧 Customizable**: Themes, rounded borders, border titles (alignment).

## 🆚 Comparison

| Feature | `gotui` | `termui` |
| :--- | :---: | :---: |
| **Renderer** | `tcell` (Optimized) | `termbox` |
| **Performance (FPS)** | **~3300** (Heavy Load) | ~1700 |
| **Widgets Available** | **23+** (Calendar, Tree, Flex...) | ~12 |
| **Layout System** | **Flex + Grid + Absolute** | Grid |
| **Customization** | **High** (Rounded Borders, Alignments) | Basic |
| **Pixel-Perfect** | **Yes** (Braille/Block/Space) | No |
| **Mouse Support** | **Full** (Wheel/Click/Drag) | Click |
| **TrueColor** | **Yes** | No |
| **SSH / Multi-User** | **Native** (Backend API) | No (Global State) |
| **Modern Terminal Support** | **All** (iterm, ghostty, etc.) | No |

`gotui` is backward compatible with `termui` and can mostly be used as a drop-in replacement.

## 📦 Installation

`gotui` uses Go modules.

```bash
go get github.com/metaspartan/gotui/v4
```

Requires Go Lang 1.24 or higher.

## 🚀 Quick Start

Create a `main.go`:

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
	p.Title = "Hello World"
	p.Text = "PRESS q TO QUIT.\n\nCombined with modern widgets, gotui aims to provide the best TUI experience in Go."
	p.SetRect(0, 0, 50, 5)
	p.TitleStyle.Fg = ui.ColorYellow
	p.BorderStyle.Fg = ui.ColorSkyBlue

	ui.Render(p)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
```

## 📚 Widgets Gallery*

Run the main dashboard demo: `go run _examples/dashboard/main.go`

Run individual examples: `go run _examples/<name>/main.go`

*Widget Screenshots are auto generated.

| Widget/Example | Screenshot | Code |
| :--- | :---: | :--- |
| **Alignment** | <img src="_examples/alignment/screenshot.png" height="80" /> | [View Example Code](_examples/alignment/main.go) |
| **Background** | <img src="_examples/background/screenshot.png" height="80" /> | [View Example Code](_examples/background/main.go) |
| **Barchart** | <img src="_examples/barchart/screenshot.png" height="80" /> | [View Example Code](_examples/barchart/main.go) |
| **Block** | <img src="_examples/block/screenshot.png" height="80" /> | [View Example Code](_examples/block/main.go) |
| **Block Multi Title** | <img src="_examples/block_multi_title/screenshot.png" height="80" /> | [View Example Code](_examples/block_multi_title/main.go) |
| **Calendar** | <img src="_examples/calendar/screenshot.png" height="80" /> | [View Example Code](_examples/calendar/main.go) |
| **Canvas** | <img src="_examples/canvas/screenshot.png" height="80" /> | [View Example Code](_examples/canvas/main.go) |
| **Collapsed Borders** | <img src="_examples/collapsed_borders/screenshot.png" height="80" /> | [View Example Code](_examples/collapsed_borders/main.go) |
| **Colors** | <img src="_examples/colors/screenshot.png" height="80" /> | [View Example Code](_examples/colors/main.go) |
| **Dashboard** | <img src="_examples/dashboard/screenshot.png" height="80" /> | [View Example Code](_examples/dashboard/main.go) |
| **Demo** | <img src="_examples/demo/screenshot.png" height="80" /> | [View Example Code](_examples/demo/main.go) |
| **Donutchart** | <img src="_examples/donutchart/screenshot.png" height="80" /> | [View Example Code](_examples/donutchart/main.go) |
| **Flex** | <img src="_examples/flex/screenshot.png" height="80" /> | [View Example Code](_examples/flex/main.go) |
| **Funnelchart** | <img src="_examples/funnelchart/screenshot.png" height="80" /> | [View Example Code](_examples/funnelchart/main.go) |
| **Gauge** | <img src="_examples/gauge/screenshot.png" height="80" /> | [View Example Code](_examples/gauge/main.go) |
| **Grid** | <img src="_examples/grid/screenshot.png" height="80" /> | [View Example Code](_examples/grid/main.go) |
| **Heatmap** | <img src="_examples/heatmap/screenshot.png" height="80" /> | [View Example Code](_examples/heatmap/main.go) |
| **Hello World** | <img src="_examples/hello_world/screenshot.png" height="80" /> | [View Example Code](_examples/hello_world/main.go) |
| **Image** | <img src="_examples/image/screenshot.png" height="80" /> | [View Example Code](_examples/image/main.go) |
| **Input** | <img src="_examples/input/screenshot.png" height="80" /> | [View Example Code](_examples/input/main.go) |
| **Linechart** | <img src="_examples/linechart/screenshot.png" height="80" /> | [View Example Code](_examples/linechart/main.go) |
| **Linegauge** | <img src="_examples/linegauge/screenshot.png" height="80" /> | [View Example Code](_examples/linegauge/main.go) |
| **List** | <img src="_examples/list/screenshot.png" height="80" /> | [View Example Code](_examples/list/main.go) |
| **Logo** | <img src="_examples/logo/screenshot.png" height="80" /> | [View Example Code](_examples/logo/main.go) |
| **Modern Demo** | <img src="_examples/modern_demo/screenshot.png" height="80" /> | [View Example Code](_examples/modern_demo/main.go) |
| **Paragraph** | <img src="_examples/paragraph/screenshot.png" height="80" /> | [View Example Code](_examples/paragraph/main.go) |
| **Piechart** | <img src="_examples/piechart/screenshot.png" height="80" /> | [View Example Code](_examples/piechart/main.go) |
| **Plot** | <img src="_examples/plot/screenshot.png" height="80" /> | [View Example Code](_examples/plot/main.go) |
| **Radarchart** | <img src="_examples/radarchart/screenshot.png" height="80" /> | [View Example Code](_examples/radarchart/main.go) |
| **Scrollbar** | <img src="_examples/scrollbar/screenshot.png" height="80" /> | [View Example Code](_examples/scrollbar/main.go) |
| **Sparkline** | <img src="_examples/sparkline/screenshot.png" height="80" /> | [View Example Code](_examples/sparkline/main.go) |
| **SSH Dashboard** | <img src="_examples/dashboard/screenshot.png" height="80" /> | [View Example Code](_examples/ssh-dashboard/main.go) |
| **Stacked Barchart** | <img src="_examples/stacked_barchart/screenshot.png" height="80" /> | [View Example Code](_examples/stacked_barchart/main.go) |
| **Table** | <img src="_examples/table/screenshot.png" height="80" /> | [View Example Code](_examples/table/main.go) |
| **Tabs** | <img src="_examples/tabs/screenshot.png" height="80" /> | [View Example Code](_examples/tabs/main.go) |
| **Textarea** | <img src="_examples/textarea/screenshot.png" height="80" /> | [View Example Code](_examples/textarea/main.go) |
| **Tree** | <img src="_examples/tree/screenshot.png" height="80" /> | [View Example Code](_examples/tree/main.go) |
| **Treemap** | <img src="_examples/treemap/screenshot.png" height="80" /> | [View Example Code](_examples/treemap/main.go) |

## 🛠️ Advanced Usage

### Customizing Borders
`gotui` supports **Rounded Borders** and various title alignments.

```go
p.Border = true
p.BorderRounded = true   // ╭───╮ instead of ┌───┐
p.Title = "My Title"
p.TitleAlignment = ui.AlignLeft // or AlignCenter, AlignRight
p.TitleBottom = "Page 1"
p.TitleBottomAlignment = ui.AlignRight
```

### Handling Mouse Events
Events include `MouseLeft`, `MouseRight`, `MouseRelease`, `MouseWheelUp`, `MouseWheelDown`.

```go
uiEvents := ui.PollEvents()
for e := range uiEvents {
    if e.Type == ui.MouseEvent {
        // e.ID is "MouseLeft", "MouseWheelUp", etc.
        // e.Payload.X, e.Payload.Y are coordinates
    }
}
```

### 🌐 Serving over SSH

You can easily serve your TUI over SSH (like standard CLI apps) using `ui.InitWithConfig` and a library like `gliderlabs/ssh`.

```go
func sshHandler(sess ssh.Session) {
    // 1. Create a custom backend for this session
    app, err := ui.NewBackend(&ui.InitConfig{
        CustomTTY: sess, // ssh.Session implements io.ReadWriter
    })
    if err != nil {
        return // Handle error appropriately
    }
    defer app.Close()

    // 2. Use the app instance instead of global ui.* functions
    p := widgets.NewParagraph()
    p.Text = "Hello SSH User!"
    p.SetRect(0, 0, 20, 5)
    
    app.Render(p) // Renders to the SSH client only!
}
```

Check `_examples/ssh-dashboard` for a full multi-user demo.

## 🤝 Contributing

Contributions are welcome! Please submit a Pull Request.

1. Fork the repo.
2. Create your feature branch (`git checkout -b feature/my-new-feature`).
3. Commit your changes (`git commit -am 'Add some feature'`).
4. Push to the branch (`git push origin feature/my-new-feature`).
5. Create a new Pull Request.

## Projects using gotui

- [mactop](https://github.com/metaspartan/mactop)

Submit a PR to add yours here!

## Author(s)

Carsen Klock (https://x.com/carsenklock)

Zack Guo (https://github.com/gizak)

## License

[MIT](http://opensource.org/licenses/MIT)

## Acknowledgments

Original [termui](https://github.com/gizak/termui) by Zack Guo.

Inspired by [Ratatui](https://github.com/ratatui-org/ratatui).
