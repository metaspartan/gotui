package widgets

import (
	"fmt"
	"image"
	"strings"

	rw "github.com/mattn/go-runewidth"
	ui "github.com/metaspartan/gotui"
)

const treeIndent = "  "

// TreeNode is a tree node.
type TreeNode struct {
	Value    fmt.Stringer
	Expanded bool
	Nodes    []*TreeNode

	// level stores the node level in the tree.
	level int
}

// TreeWalkFn is a function used for walking a Tree.
// To interrupt the walking process function should return false.
type TreeWalkFn func(*TreeNode) bool

func (tn *TreeNode) parseStyles(style ui.Style) []ui.Cell {
	var sb strings.Builder
	if len(tn.Nodes) == 0 {
		sb.WriteString(strings.Repeat(treeIndent, tn.level+1))
	} else {
		sb.WriteString(strings.Repeat(treeIndent, tn.level))
		if tn.Expanded {
			sb.WriteRune(ui.Theme.Tree.Expanded)
		} else {
			sb.WriteRune(ui.Theme.Tree.Collapsed)
		}
		sb.WriteByte(' ')
	}
	sb.WriteString(tn.Value.String())
	return ui.ParseStyles(sb.String(), style)
}

// Tree is a tree widget.
type Tree struct {
	ui.Block
	TextStyle        ui.Style
	SelectedRowStyle ui.Style
	WrapText         bool
	SelectedRow      int

	nodes []*TreeNode
	// rows is flatten nodes for rendering.
	rows   []*TreeNode
	topRow int
}

// NewTree creates a new Tree widget.
func NewTree() *Tree {
	return &Tree{
		Block:            *ui.NewBlock(),
		TextStyle:        ui.Theme.Tree.Text,
		SelectedRowStyle: ui.Theme.Tree.Text,
		WrapText:         true,
	}
}

func (t *Tree) SetNodes(nodes []*TreeNode) {
	t.nodes = nodes
	t.prepareNodes()
}

func (t *Tree) prepareNodes() {
	t.rows = make([]*TreeNode, 0)
	for _, node := range t.nodes {
		t.prepareNode(node, 0)
	}
}

func (t *Tree) prepareNode(node *TreeNode, level int) {
	t.rows = append(t.rows, node)
	node.level = level

	if node.Expanded {
		for _, n := range node.Nodes {
			t.prepareNode(n, level+1)
		}
	}
}

func (t *Tree) Walk(fn TreeWalkFn) {
	for _, n := range t.nodes {
		if !t.walk(n, fn) {
			break
		}
	}
}

func (t *Tree) walk(n *TreeNode, fn TreeWalkFn) bool {
	if !fn(n) {
		return false
	}

	for _, node := range n.Nodes {
		if !t.walk(node, fn) {
			return false
		}
	}

	return true
}

func (t *Tree) Draw(buf *ui.Buffer) {
	t.Block.Draw(buf)
	point := t.Inner.Min

	// adjusts view into widget
	if t.SelectedRow >= t.Inner.Dy()+t.topRow {
		t.topRow = t.SelectedRow - t.Inner.Dy() + 1
	} else if t.SelectedRow < t.topRow {
		t.topRow = t.SelectedRow
	}

	// draw rows
	for row := t.topRow; row < len(t.rows) && point.Y < t.Inner.Max.Y; row++ {
		cells := t.rows[row].parseStyles(t.TextStyle)
		if t.WrapText {
			cells = ui.WrapCells(cells, uint(t.Inner.Dx()))
		}
		for j := 0; j < len(cells) && point.Y < t.Inner.Max.Y; j++ {
			style := cells[j].Style
			if row == t.SelectedRow {
				style = t.SelectedRowStyle
			}
			if point.X+1 == t.Inner.Max.X+1 && len(cells) > t.Inner.Dx() {
				buf.SetCell(ui.NewCell(ui.ELLIPSES, style), point.Add(image.Pt(-1, 0)))
			} else {
				buf.SetCell(ui.NewCell(cells[j].Rune, style), point)
				point = point.Add(image.Pt(rw.RuneWidth(cells[j].Rune), 0))
			}
		}
		point = image.Pt(t.Inner.Min.X, point.Y+1)
	}

	// draw UP_ARROW if needed
	if t.topRow > 0 {
		buf.SetCell(
			ui.NewCell(ui.UP_ARROW, ui.NewStyle(ui.ColorWhite)),
			image.Pt(t.Inner.Max.X-1, t.Inner.Min.Y),
		)
	}

	// draw DOWN_ARROW if needed
	if len(t.rows) > int(t.topRow)+t.Inner.Dy() {
		buf.SetCell(
			ui.NewCell(ui.DOWN_ARROW, ui.NewStyle(ui.ColorWhite)),
			image.Pt(t.Inner.Max.X-1, t.Inner.Max.Y-1),
		)
	}
}

// ScrollAmount scrolls by amount given. If amount is < 0, then scroll up.
// There is no need to set t.topRow, as this will be set automatically when drawn,
// since if the selected item is off screen then the topRow variable will change accordingly.
func (t *Tree) ScrollAmount(amount int) {
	if len(t.rows)-int(t.SelectedRow) <= amount {
		t.SelectedRow = len(t.rows) - 1
	} else if int(t.SelectedRow)+amount < 0 {
		t.SelectedRow = 0
	} else {
		t.SelectedRow += amount
	}
}

func (t *Tree) SelectedNode() *TreeNode {
	if len(t.rows) == 0 {
		return nil
	}
	return t.rows[t.SelectedRow]
}

func (t *Tree) ScrollUp() {
	t.ScrollAmount(-1)
}

func (t *Tree) ScrollDown() {
	t.ScrollAmount(1)
}

func (t *Tree) ScrollPageUp() {
	// If an item is selected below top row, then go to the top row.
	if t.SelectedRow > t.topRow {
		t.SelectedRow = t.topRow
	} else {
		t.ScrollAmount(-t.Inner.Dy())
	}
}

func (t *Tree) ScrollPageDown() {
	t.ScrollAmount(t.Inner.Dy())
}

func (t *Tree) ScrollHalfPageUp() {
	t.ScrollAmount(-int(ui.FloorFloat64(float64(t.Inner.Dy()) / 2)))
}

func (t *Tree) ScrollHalfPageDown() {
	t.ScrollAmount(int(ui.FloorFloat64(float64(t.Inner.Dy()) / 2)))
}

func (t *Tree) ScrollTop() {
	t.SelectedRow = 0
}

func (t *Tree) ScrollBottom() {
	t.SelectedRow = len(t.rows) - 1
}

func (t *Tree) Collapse() {
	t.rows[t.SelectedRow].Expanded = false
	t.prepareNodes()
}

func (t *Tree) Expand() {
	node := t.rows[t.SelectedRow]
	if len(node.Nodes) > 0 {
		t.rows[t.SelectedRow].Expanded = true
	}
	t.prepareNodes()
}

func (t *Tree) ToggleExpand() {
	node := t.rows[t.SelectedRow]
	if len(node.Nodes) > 0 {
		node.Expanded = !node.Expanded
	}
	t.prepareNodes()
}

func (t *Tree) ExpandAll() {
	t.Walk(func(n *TreeNode) bool {
		if len(n.Nodes) > 0 {
			n.Expanded = true
		}
		return true
	})
	t.prepareNodes()
}

func (t *Tree) CollapseAll() {
	t.Walk(func(n *TreeNode) bool {
		n.Expanded = false
		return true
	})
	t.prepareNodes()
}
