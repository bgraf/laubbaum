package main

import (
	"github.com/gdamore/tcell/v2"
)

type InputField struct {
	buffer  []rune
	width   int
	height  int
	x, y    int
	isDone  bool
	initial string
	cursor  int
}

func NewInputField(x, y int) *InputField {
	return &InputField{
		width:  20,
		height: 1,
		x:      x,
		y:      y,
	}
}

func (i *InputField) SetInitialText(initial string) {
	i.initial = initial
	i.reset()
}

func (i *InputField) updateWidth() {
	if i.width < len(i.buffer) {
		i.width = len(i.buffer)
	}
}

func (i *InputField) String() string {
	return string(i.buffer)
}

func (i *InputField) reset() {
	i.Clear()
	i.buffer = []rune(i.initial)
	i.cursor = len(i.buffer)
}

func (i *InputField) Clear() {
	i.buffer = []rune{}
}

func (inp *InputField) Draw(s tcell.Screen) {
	inp.updateWidth()

	line := 0
	col := 0

	for i, c := range inp.buffer {
		if i == inp.cursor {
			s.ShowCursor(inp.x+col, inp.y+line+1)
		}

		if c == '\n' {
			line++
			col = 0
			continue
		}
		s.SetContent(inp.x+col, inp.y+1+line, c, nil, defStyle)

		col++

	}

	if inp.cursor >= len(inp.buffer) {
		s.ShowCursor(inp.x+col, inp.y+line+1)
	}
}

func (i *InputField) IsDone() bool {
	return i.isDone
}

func (i *InputField) Insert(char rune) {
	if len(i.buffer) <= i.cursor {
		i.buffer = append(i.buffer, char)
		i.incrCursor()

		return
	}

	i.buffer = append(i.buffer, ' ')
	copy(i.buffer[i.cursor+1:], i.buffer[i.cursor:])
	i.buffer[i.cursor] = char
	i.incrCursor()
}

func (i *InputField) incrCursor() {
	if i.cursor < len(i.buffer) {
		i.cursor++
	}
}

func (i *InputField) decrCursor() {
	if 0 < i.cursor {
		i.cursor--
	}
}

func (i *InputField) HandleEvent(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyBackspace2:
		if len(i.buffer) > 0 {
			i.buffer = i.buffer[:len(i.buffer)-1]
			i.decrCursor()
		}
	case tcell.KeyRune:
		i.Insert(ev.Rune())
	case tcell.KeyEnter:
		if ev.Modifiers()&tcell.ModAlt > 0 {
			i.Insert('\n')
		} else {
			i.isDone = true
		}
	case tcell.KeyEscape:
		i.reset()
		i.isDone = true

	case tcell.KeyLeft:
		i.decrCursor()
	case tcell.KeyRight:
		i.incrCursor()
	}
}
