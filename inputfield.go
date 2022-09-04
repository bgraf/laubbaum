package main

import "github.com/gdamore/tcell/v2"

type InputField struct {
	buffer []rune
	width  int
	x, y   int
}

func NewInputField(x, y int) *InputField {
	return &InputField{
		width: 20,
		x:     x,
		y:     y,
	}
}

func (i *InputField) updateWidth() {
	if i.width < len(i.buffer) {
		i.width = len(i.buffer)
	}
}

func (inp *InputField) Draw(s tcell.Screen) {
	inp.updateWidth()

	for i := 0; i < inp.width; i++ {
		s.SetContent(inp.x+i, inp.y+0, tcell.RuneHLine, nil, defStyle)
		s.SetContent(inp.x+i, inp.y+2, tcell.RuneHLine, nil, defStyle)
	}

	for i, r := range inp.buffer {
		s.SetContent(inp.x+i, inp.y+1, r, nil, defStyle)
	}
	for i := len(inp.buffer); i < inp.width; i++ {
		s.SetContent(inp.x+i, inp.y+1, ' ', nil, defStyle)
	}

	s.ShowCursor(inp.x+len(inp.buffer), inp.y+1)
}

func (i *InputField) HandleEvent(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyBackspace2:
		if len(i.buffer) > 0 {
			i.buffer = i.buffer[:len(i.buffer)-1]
		}
	case tcell.KeyRune:
		i.buffer = append(i.buffer, ev.Rune())
	case tcell.KeyEnter:
		return true
	}

	return false
}
