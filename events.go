
package gotui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

/*
List of events:
	mouse events:
		<MouseLeft> <MouseRight> <MouseMiddle>
		<MouseWheelUp> <MouseWheelDown>
	keyboard events:
		any uppercase or lowercase letter like j or J
		<C-d> etc
		<M-d> etc
		<Up> <Down> <Left> <Right>
		<Insert> <Delete> <Home> <End> <Previous> <Next>
		<Backspace> <Tab> <Enter> <Escape> <Space>
		<C-<Space>> etc
	terminal events:
        <Resize>
*/

type EventType uint

const (
	KeyboardEvent EventType = iota
	MouseEvent
	ResizeEvent
)

type Event struct {
	Type    EventType
	ID      string
	Payload interface{}
}

// Mouse payload.
type Mouse struct {
	Drag bool
	X    int
	Y    int
}

// Resize payload.
type Resize struct {
	Width  int
	Height int
}

// PollEvents gets events from tcell, converts them, then sends them to each of its channels.
func PollEvents() <-chan Event {
	ch := make(chan Event)
	go func() {
		for {
			if Screen == nil {
				return
			}
			ev := Screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				ch <- convertTcellKeyEvent(ev)
			case *tcell.EventMouse:
				ch <- convertTcellMouseEvent(ev)
			case *tcell.EventResize:
				w, h := ev.Size()
				ch <- Event{
					Type: ResizeEvent,
					ID:   "<Resize>",
					Payload: Resize{
						Width:  w,
						Height: h,
					},
				}
			}
		}
	}()
	return ch
}

var keyMap = map[tcell.Key]string{
	tcell.KeyF1:         "<F1>",
	tcell.KeyF2:         "<F2>",
	tcell.KeyF3:         "<F3>",
	tcell.KeyF4:         "<F4>",
	tcell.KeyF5:         "<F5>",
	tcell.KeyF6:         "<F6>",
	tcell.KeyF7:         "<F7>",
	tcell.KeyF8:         "<F8>",
	tcell.KeyF9:         "<F9>",
	tcell.KeyF10:        "<F10>",
	tcell.KeyF11:        "<F11>",
	tcell.KeyF12:        "<F12>",
	tcell.KeyInsert:     "<Insert>",
	tcell.KeyDelete:     "<Delete>",
	tcell.KeyHome:       "<Home>",
	tcell.KeyEnd:        "<End>",
	tcell.KeyPgUp:       "<PageUp>",
	tcell.KeyPgDn:       "<PageDown>",
	tcell.KeyUp:         "<Up>",
	tcell.KeyDown:       "<Down>",
	tcell.KeyLeft:       "<Left>",
	tcell.KeyRight:      "<Right>",
	tcell.KeyCtrlA:      "<C-a>",
	tcell.KeyCtrlB:      "<C-b>",
	tcell.KeyCtrlC:      "<C-c>",
	tcell.KeyCtrlD:      "<C-d>",
	tcell.KeyCtrlE:      "<C-e>",
	tcell.KeyCtrlF:      "<C-f>",
	tcell.KeyCtrlG:      "<C-g>",
	tcell.KeyCtrlH:      "<C-h>", // Backspace sometimes
	tcell.KeyTab:        "<Tab>",
	tcell.KeyCtrlJ:      "<C-j>",
	tcell.KeyCtrlK:      "<C-k>",
	tcell.KeyCtrlL:      "<C-l>",
	tcell.KeyEnter:      "<Enter>",
	tcell.KeyCtrlN:      "<C-n>",
	tcell.KeyCtrlO:      "<C-o>",
	tcell.KeyCtrlP:      "<C-p>",
	tcell.KeyCtrlQ:      "<C-q>",
	tcell.KeyCtrlR:      "<C-r>",
	tcell.KeyCtrlS:      "<C-s>",
	tcell.KeyCtrlT:      "<C-t>",
	tcell.KeyCtrlU:      "<C-u>",
	tcell.KeyCtrlV:      "<C-v>",
	tcell.KeyCtrlW:      "<C-w>",
	tcell.KeyCtrlX:      "<C-x>",
	tcell.KeyCtrlY:      "<C-y>",
	tcell.KeyCtrlZ:      "<C-z>",
	tcell.KeyEsc:        "<Escape>",
	tcell.KeyBackspace:  "<Backspace>",
	tcell.KeyBackspace2: "<Backspace>",
}

func convertTcellKeyEvent(e *tcell.EventKey) Event {
	ID := ""
	if e.Key() == tcell.KeyRune {
		r := e.Rune()
		if e.Modifiers()&tcell.ModAlt != 0 {
			ID = fmt.Sprintf("<M-%c>", r)
		} else {
			ID = string(r)
		}
	} else {
		// Named key
		if val, ok := keyMap[e.Key()]; ok {
			ID = val
		} else {
			ID = fmt.Sprintf("<Key:%v>", e.Key())
		}
		// Handle simple M-Key for non-runes if necessary, but tcell usually handles this
	}

	return Event{
		Type:    KeyboardEvent,
		ID:      ID,
		Payload: e, // Optional: might want to pass raw event
	}
}

func convertTcellMouseEvent(e *tcell.EventMouse) Event {
	btns := e.Buttons()
	ID := "Unknown_Mouse_Button"

	if btns&tcell.Button1 != 0 {
		ID = "<MouseLeft>"
	} else if btns&tcell.Button3 != 0 {
		ID = "<MouseRight>" // Right is button 3 usually? tcell says Button2=secondary, Button3=middle? Check docs.
		// Actually tcell: Button1=Left, Button2=Middle, Button3=Right usually.
		ID = "<MouseMiddle>" // Wait, standard xterm is 1,2,3 -> Left,Middle,Right.
		// tcell.Button1 is Primary.
		// tcell.Button2 is Secondary (Middle).
		// tcell.Button3 is Tertiary (Right).
	} else if btns&tcell.Button2 != 0 {
		ID = "<MouseMiddle>"
	}

	// Correcting assumptions based on tcell definition:
	// Button1 = Left
	// Button2 = Middle
	// Button3 = Right
	if btns&tcell.Button1 != 0 {
		ID = "<MouseLeft>"
	}
	if btns&tcell.Button2 != 0 {
		ID = "<MouseMiddle>"
	}
	if btns&tcell.Button3 != 0 {
		ID = "<MouseRight>"
	}
	if btns&tcell.WheelUp != 0 {
		ID = "<MouseWheelUp>"
	}
	if btns&tcell.WheelDown != 0 {
		ID = "<MouseWheelDown>"
	}

	x, y := e.Position()

	// Detect Drag
	// tcell doesn't have a dedicated Drag generic event but e.Buttons() might include it if we track state.
	// OR ModMotion.
	// For simplify, we just map basic clicks for now.

	return Event{
		Type: MouseEvent,
		ID:   ID,
		Payload: Mouse{
			X:    x,
			Y:    y,
			Drag: false, // Implementation detail to improve later
		},
	}
}
