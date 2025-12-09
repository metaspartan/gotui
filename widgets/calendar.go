// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"fmt"
	"image"
	"time"

	ui "github.com/metaspartan/gotui"
)

type Calendar struct {
	ui.Block
	Month       time.Month
	Year        int
	CurrentDay  int
	SelectedDay int
	HeaderStyle ui.Style
	DayStyle    ui.Style
}

func NewCalendar() *Calendar {
	now := time.Now()
	return &Calendar{
		Block:       *ui.NewBlock(),
		Month:       now.Month(),
		Year:        now.Year(),
		CurrentDay:  now.Day(),
		SelectedDay: now.Day(),
		HeaderStyle: ui.Theme.Block.Title,
		DayStyle:    ui.Theme.Paragraph.Text,
	}
}

func (c *Calendar) Draw(buf *ui.Buffer) {
	c.Block.Draw(buf)

	// Center Header: Month Year
	header := fmt.Sprintf("%s %d", c.Month, c.Year)
	headerX := c.Inner.Min.X + (c.Inner.Dx()-len(header))/2
	buf.SetString(header, c.HeaderStyle, image.Pt(headerX, c.Inner.Min.Y))

	// Days of week
	y := c.Inner.Min.Y + 2
	x := c.Inner.Min.X // Padding?

	// Simple layout: 3 chars per day + 9 gap
	days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	for i, d := range days {
		buf.SetString(d, c.HeaderStyle, image.Pt(x+(i*4), y))
	}
	y++

	// Determine start day
	t := time.Date(c.Year, c.Month, 1, 0, 0, 0, 0, time.Local)
	startWeekday := int(t.Weekday()) // 0=Sun

	// Determine days in month
	nextMonth := t.AddDate(0, 1, 0)
	daysInMonth := nextMonth.Sub(t).Hours() / 24

	day := 1
	row := 0

	// Initial padding
	col := startWeekday

	for day <= int(daysInMonth) {
		dx := x + (col * 4)
		dy := y + row

		style := c.DayStyle
		if day == c.CurrentDay {
			style.Bg = ui.ColorGreen
			style.Fg = ui.ColorBlack
		}

		buf.SetString(
			fmt.Sprintf("%3d", day),
			style,
			image.Pt(dx, dy),
		)

		day++
		col++
		if col > 6 {
			col = 0
			row++
		}
	}
}
