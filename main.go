package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"

	"github.com/bgraf/laubbaum/model"
)

func main() {
	fmt.Println("ok")

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("init: %+v", err)
	}

	if err := s.Init(); err != nil {
		log.Fatalf("init: %+v", err)
	}

	mod := defaultModel()
	mod.PropagateParent(nil)

	app := &application{
		screen:      s,
		diagramView: newRenderState(),
		root:        mod,
	}

	// Set default text style
	s.SetStyle(defStyle)
	s.Clear()

	for {
		drawApp(app)
		ev := s.PollEvent()
		handleEvent(app, ev)
	}
}

var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

type application struct {
	screen      tcell.Screen
	diagramView *renderState
	inputField  *InputField
	root        *model.Node
}

func drawApp(app *application) {
	s := app.screen

	_, h := s.Size()
	bc, br := 1, h/2

	s.Fill(' ', defStyle)
	app.diagramView.Draw(s, app.root, bc, br)
	if app.inputField != nil {
		app.inputField.Draw(s)
	}

	// Update screen
	s.Show()
}

func openInput(app *application, init string, clear bool) {
	w, _ := app.screen.Size()
	app.inputField = NewInputField(w/2-20, 0)
	app.inputField.SetInitialText(init)
	if clear {
		app.inputField.Clear()
	}
}

func quit(app *application) {
	app.screen.Fini()
	os.Exit(0)
}

func handleEvent(app *application, ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventResize:
		app.screen.Sync()
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyCtrlC {
			quit(app)
		}

		if app.inputField != nil {
			app.inputField.HandleEvent(ev)
			app.diagramView.selected.SetText(app.inputField.String())
			if app.inputField.IsDone() {
				app.screen.ShowCursor(-1, -1)
				app.inputField = nil
			}

			return
		}

		switch ev.Key() {
		case tcell.KeyUp:
			fallthrough
		case tcell.KeyDown:
			fallthrough
		case tcell.KeyLeft:
			fallthrough
		case tcell.KeyRight:
			app.diagramView.SelectNext(ev.Key())
		case tcell.KeyEnter:
			state := app.diagramView
			curr := state.selected

			if curr.Parent != nil {
				newChild := &model.Node{
					Text:   "",
					Parent: curr.Parent,
				}
				curr.Parent.Children = append(curr.Parent.Children, newChild)
				state.selected = newChild

				openInput(app, "", true)
			}

		case tcell.KeyTab:
			state := app.diagramView
			curr := state.selected

			newChild := &model.Node{
				Text:   "",
				Parent: curr,
			}
			curr.Children = append(curr.Children, newChild)
			state.selected = newChild

			openInput(app, newChild.Text, true)

		case tcell.KeyDelete:
			state := app.diagramView
			curr := state.selected
			if curr.Parent != nil {

				next := curr.NextSibling()
				if next == nil {
					next = curr.PreviousSibling()
				}
				if next == nil {
					next = curr.Parent
				}

				writeI := 0
				for _, c := range curr.Parent.Children {
					if c != curr {
						curr.Parent.Children[writeI] = c
						writeI++
					}
				}

				curr.Parent.Children = curr.Parent.Children[:writeI]

				state.selected = next
			}

		case tcell.KeyRune:
			switch ev.Rune() {
			case 'O':
				state := app.diagramView
				curr := state.selected

				if curr.Parent != nil {
					newChild := &model.Node{
						Text:   "new node!",
						Parent: curr.Parent,
					}

					index, _ := curr.Parent.ChildIndex(curr)

					curr.Parent.InsertChild(newChild, index)
					state.selected = newChild

					openInput(app, "", true)
				}

			case 'o':
				state := app.diagramView
				curr := state.selected

				if curr.Parent != nil {
					newChild := &model.Node{
						Text:   "new node!",
						Parent: curr.Parent,
					}

					index, _ := curr.Parent.ChildIndex(curr)

					curr.Parent.InsertChild(newChild, index+1)
					state.selected = newChild

					openInput(app, "", true)
				}

			case 'a':
				openInput(app, app.diagramView.selected.Text, false)
			case 'c':
				openInput(app, app.diagramView.selected.Text, true)
			}
		}
	}
}

func defaultModel() *model.Node {
	root := &model.Node{
		Text: "Topic",
	}

	root.Children = append(root.Children, &model.Node{
		Text:   "hallo\nwelt",
		Parent: root,
	})

	return root
}
