# Changelog

## Fork to gotui [1.0.0] - Upgrade (2025)

### Added
- **Backend Migration**: Switched from `termbox-go` to `tcell` v2 for modern terminal support.
- **New Widgets**:
  - `Heatmap`: Visualizes 2D data intensities (`widgets/heatmap.go`).
  - `Calendar`: Month/Year view with day highlighting (`widgets/calendar.go`).
  - `TextArea`: Multi-line text input with navigation and editing (`widgets/textarea.go`).
  - `Input`: Single-line text input with password support (`widgets/input.go`).
- **Styling & Alignment**:
  - **TrueColor**: Full 24-bit RGB support via `tcell.NewRGBColor` / `ui.NewRGBColor`.
  - **Extended Colors**: Added 48 named CSS colors to `style.go`.
  - **Bottom Titles**: `Block` now supports `TitleBottom`.
  - **Alignment**:
    - `Block` title alignment (Left, Center, Right) for both top and bottom titles.
    - `Paragraph` vertical and horizontal alignment.
    - `List` text alignment (Left, Center, Right).
- **Features**:
  - **Mouse Wheel**: Added support for `<MouseWheelUp>` and `<MouseWheelDown>` events.
- **Examples**: Added `heatmap.go`, `calendar.go`, `alignment.go`, `colors.go`, `textarea.go`, `input.go`. Updated examples with mouse support.
- **TrueColor**: Full 24-bit RGB color support (`NewRGBColor`).
- **Feature**: Bottom Border Titles (`Block.TitleBottom`).
- **Input**: Improved mouse and key event handling (modern protocols).
- **Transparency**: Default background color is now transparent (`ColorDefault`).

### Changed
- **Module**: Renamed to `github.com/metaspartan/gotui`.
- **Package**: Renamed `termui` to `gotui`.
- **Cleanup**: Removed legacy copyright headers and deprecated comments.