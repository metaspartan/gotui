package gotui

func NewGrid() *Grid {
	g := &Grid{
		Block: *NewBlock(),
	}
	g.Border = false
	return g
}

func NewCol(ratio float64, i ...interface{}) GridItem {
	_, ok := i[0].(Drawable)
	entry := i[0]
	if !ok {
		entry = i
	}
	return GridItem{
		Type:   col,
		Entry:  entry,
		IsLeaf: ok,
		ratio:  ratio,
	}
}

func NewRow(ratio float64, i ...interface{}) GridItem {
	_, ok := i[0].(Drawable)
	entry := i[0]
	if !ok {
		entry = i
	}
	return GridItem{
		Type:   row,
		Entry:  entry,
		IsLeaf: ok,
		ratio:  ratio,
	}
}

func (g *Grid) Set(entries ...interface{}) {
	entry := GridItem{
		Type:   row,
		Entry:  entries,
		IsLeaf: false,
		ratio:  1.0,
	}
	g.setHelper(entry, 1.0, 1.0)
}

func (g *Grid) setHelper(item GridItem, parentWidthRatio, parentHeightRatio float64) {
	var HeightRatio float64
	var WidthRatio float64
	switch item.Type {
	case col:
		HeightRatio = 1.0
		WidthRatio = item.ratio
	case row:
		HeightRatio = item.ratio
		WidthRatio = 1.0
	}
	item.WidthRatio = parentWidthRatio * WidthRatio
	item.HeightRatio = parentHeightRatio * HeightRatio

	if item.IsLeaf {
		g.Items = append(g.Items, &item)
	} else {
		XRatio := 0.0
		YRatio := 0.0
		cols := false
		rows := false

		children := InterfaceSlice(item.Entry)

		for i := 0; i < len(children); i++ {
			if children[i] == nil {
				continue
			}
			child, _ := children[i].(GridItem)

			child.XRatio = item.XRatio + (item.WidthRatio * XRatio)
			child.YRatio = item.YRatio + (item.HeightRatio * YRatio)

			switch child.Type {
			case col:
				cols = true
				XRatio += child.ratio
				if rows {
					item.HeightRatio /= 2
				}
			case row:
				rows = true
				YRatio += child.ratio
				if cols {
					item.WidthRatio /= 2
				}
			}

			g.setHelper(child, item.WidthRatio, item.HeightRatio)
		}
	}
}

func (g *Grid) Draw(buf *Buffer) {
	width := float64(g.Dx()) + 1
	height := float64(g.Dy()) + 1

	for _, item := range g.Items {
		entry, _ := item.Entry.(Drawable)

		x := int(width*item.XRatio) + g.Min.X
		y := int(height*item.YRatio) + g.Min.Y
		w := int(width * item.WidthRatio)
		h := int(height * item.HeightRatio)

		if x+w > g.Dx() {
			w--
		}
		if y+h > g.Dy() {
			h--
		}

		entry.SetRect(x, y, x+w, y+h)

		entry.Lock()
		entry.Draw(buf)
		entry.Unlock()
	}
}
