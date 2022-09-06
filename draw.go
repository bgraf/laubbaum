package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/bgraf/laubbaum/model"
)

type nodeState struct {
	totalHeight int
}

type renderState struct {
	nodes           map[*model.Node]nodeState
	maxWidthByLevel []int
	selected        *model.Node
}

func newRenderState() *renderState {
	return &renderState{
		nodes: make(map[*model.Node]nodeState),
	}
}

func (r *renderState) Clear() {
	r.nodes = make(map[*model.Node]nodeState)
	for i := range r.maxWidthByLevel {
		r.maxWidthByLevel[i] = 0
	}
}

func (r *renderState) updatedMaxWidth(level, width int) {
	for level >= len(r.maxWidthByLevel) {
		r.maxWidthByLevel = append(r.maxWidthByLevel, 0)
	}

	currentMax := r.maxWidthByLevel[level]
	if currentMax < width {
		r.maxWidthByLevel[level] = width
	}
}

func (r *renderState) Setup(root *model.Node, level int) {
	totalHeight := 0

	for _, child := range root.Children {
		r.Setup(child, level+1)

		childState, ok := r.nodes[child]
		if !ok {
			panic("muh")
		}

		totalHeight += childState.totalHeight
	}

	nChildren := len(root.Children)
	if nChildren > 1 {
		totalHeight += nChildren - 1 // additional spacing between children
	}

	w, rh := NodeSize(root)
	if totalHeight < rh {
		totalHeight = rh
	}

	r.updatedMaxWidth(level, w)

	r.nodes[root] = nodeState{
		totalHeight: totalHeight,
	}
}

func (r *renderState) nodeState(n *model.Node) nodeState {
	if state, ok := r.nodes[n]; ok {
		return state
	}

	panic("no such node")
}

func (r *renderState) Height(n *model.Node) int {
	return r.nodeState(n).totalHeight
}

func (r *renderState) SelectNext(k tcell.Key) {
	switch k {
	case tcell.KeyLeft:
		if r.selected.Parent != nil {
			r.selected = r.selected.Parent
		}

	case tcell.KeyRight:
		nc := len(r.selected.Children)
		if nc > 0 {
			r.selected = r.selected.Children[(nc-1)/2]
		}

	case tcell.KeyUp:
		if candidate := r.selected.PreviousSibling(); candidate != nil {
			r.selected = candidate
		}

	case tcell.KeyDown:
		if candidate := r.selected.NextSibling(); candidate != nil {
			r.selected = candidate
		}
	}
}

func (r *renderState) Draw(s tcell.Screen, root *model.Node, col, row int) {
	r.Clear()
	r.Setup(root, 0)

	if r.selected == nil {
		r.selected = root
	}

	drawRec(r, s, root, col, row, 0)
}

func drawRec(state *renderState, s tcell.Screen, n *model.Node, col, row, level int) {
	wNode, _ := NodeSize(n)
	w := state.maxWidthByLevel[level]

	isSelected := state.selected == n
	DrawNode(s, n, col, row, w, isSelected)

	height := state.Height(n)
	heightUsed := row - height/2

	childConnectors := make([]Point, len(n.Children))

	for i, child := range n.Children {
		childHeight := state.Height(child)
		offset := heightUsed + childHeight/2

		childColumnOffset := col + w + 9

		drawRec(state, s, child, childColumnOffset, offset, level+1)

		childConnectors[i] = Point{C: childColumnOffset, R: offset}

		heightUsed += childHeight + 1 // render spacing
	}

	if len(childConnectors) > 0 {
		DrawConnectors(s, Point{C: col + w, R: row}, childConnectors)
		_ = wNode
	}
}

func DrawNode(s tcell.Screen, n *model.Node, col, row int, minWidth int, isSelected bool) {
	w, h := NodeSize(n)

	if w < minWidth {
		w = minWidth
	}

	row -= h / 2

	defStyle := tcell.StyleDefault
	if isSelected {
		defStyle = tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	}

	DrawBox(s, col, row, w, h, defStyle)

	parts := strings.Split(n.Text, "\n")

	for i, part := range parts {
		drawText(s, col+2, row+1+i, col+w-2, row+1+i, tcell.StyleDefault, part)
	}
}

func NodeSize(n *model.Node) (w, h int) {
	w, h = n.InnerSize()

	w += 4 // padding and border
	h += 2 // border

	return w, h
}

func DrawBox(s tcell.Screen, col, row, w, h int, defStyle tcell.Style) {
	for c := col + 1; c < col+w-1; c++ {
		s.SetContent(c, row, '─', nil, defStyle)
		s.SetContent(c, row+h-1, '─', nil, defStyle)
	}

	for r := row + 1; r < row+h-1; r++ {
		s.SetContent(col, r, '│', nil, defStyle)
		s.SetContent(col+w-1, r, '│', nil, defStyle)

	}

	s.SetContent(col, row, '┌', nil, defStyle)
	s.SetContent(col+w-1, row, '┐', nil, defStyle)
	s.SetContent(col, row+h-1, '└', nil, defStyle)
	s.SetContent(col+w-1, row+h-1, '┘', nil, defStyle)

	// sty := tcell.StyleDefault.Bold(true)

	// drawText(s, col+2, row+1, col+w-2, row+h-2, sty, "Title ends with dot.")
	// drawText(s, col+2, row+2, col+w-2, row+h-2, defStyle, "hello world, how are you?")
}

func Drawf(s tcell.Screen, x1, y1 int, format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)

	row := y1
	col := x1
	for _, r := range text {
		s.SetContent(col, row, r, nil, tcell.StyleDefault)
		col++
	}
}

type Point struct {
	C, R int
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func DrawConnectors(s tcell.Screen, out Point, in []Point) {
	minC := 1 << 20
	minR := out.R
	maxR := out.R

	for _, p := range in {
		minR = MinInt(minR, p.R)
		minC = MinInt(minC, p.C)
		maxR = MaxInt(maxR, p.R)
	}

	if minC < out.C {
		panic("expected out.C <= min in.C")
	}

	midC := (minC + out.C) / 2

	for c := out.C; c < midC; c++ {
		s.SetContent(c, out.R, tcell.RuneHLine, nil, tcell.StyleDefault)
	}

	for r := minR; r < maxR+1; r++ {
		s.SetContent(midC, r, tcell.RuneVLine, nil, tcell.StyleDefault)
	}

	for _, p := range in {
		for c := midC + 1; c < p.C; c++ {
			s.SetContent(c, p.R, tcell.RuneHLine, nil, tcell.StyleDefault)
		}

		s.SetContent(p.C, p.R, tcell.RuneRTee, nil, tcell.StyleDefault)

		switch p.R {
		case minR:
			s.SetContent(midC, p.R, tcell.RuneULCorner, nil, tcell.StyleDefault)
		case maxR:
			s.SetContent(midC, p.R, tcell.RuneLLCorner, nil, tcell.StyleDefault)
		default:
			s.SetContent(midC, p.R, tcell.RuneLTee, nil, tcell.StyleDefault)
		}
	}

	right, _, _, _ := s.GetContent(midC+1, out.R)

	// Mid-Connector
	switch {
	case minR < out.R && maxR > out.R && right == tcell.RuneHLine:
		s.SetContent(midC, out.R, '┼', nil, tcell.StyleDefault)
	case minR < out.R && maxR > out.R:
		s.SetContent(midC, out.R, '┤', nil, tcell.StyleDefault)
	case minR < out.R:
		s.SetContent(midC, out.R, '┘', nil, tcell.StyleDefault)
	case maxR > out.R:
		s.SetContent(midC, out.R, '┐', nil, tcell.StyleDefault)
	default:
		s.SetContent(midC, out.R, tcell.RuneHLine, nil, tcell.StyleDefault)
	}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
