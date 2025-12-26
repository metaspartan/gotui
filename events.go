package gotui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v3"
)

// PollEvents polls for events.
func PollEvents() <-chan Event {
	return DefaultBackend.PollEvents()
}

// PollEventsWithContext polls for events with a context.
func PollEventsWithContext(ctx context.Context) <-chan Event {
	return DefaultBackend.PollEventsWithContext(ctx)
}

// PollEvents polls for events.
func (b *Backend) PollEvents() <-chan Event {
	ch := make(chan Event)
	go func() {
		defer close(ch)
		for {
			if b.Screen == nil {
				return
			}
			ev, ok := <-b.Screen.EventQ()
			if !ok {
				return
			}
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

// PollEventsWithContext polls for events with a context.
func (b *Backend) PollEventsWithContext(ctx context.Context) <-chan Event {
	ch := make(chan Event)
	go func() {
		defer close(ch)
		events := make(chan Event)
		stopPoller := make(chan struct{})
		go b.runPoller(events, stopPoller)
		for {
			select {
			case <-ctx.Done():
				close(stopPoller)
				if b.Screen != nil {
					// Non-blocking send to avoid deadlock if runPoller already exited
					select {
					case b.Screen.EventQ() <- tcell.NewEventInterrupt(nil):
					default:
					}
				}
				return
			case ev, ok := <-events:
				if !ok {
					return
				}
				select {
				case ch <- ev:
				case <-ctx.Done():
					close(stopPoller)
					if b.Screen != nil {
						// Non-blocking send to avoid deadlock if runPoller already exited
						select {
						case b.Screen.EventQ() <- tcell.NewEventInterrupt(nil):
						default:
						}
					}
					return
				}
			}
		}
	}()
	return ch
}
func (b *Backend) runPoller(events chan<- Event, stop <-chan struct{}) {
	defer close(events)
	for {
		select {
		case <-stop:
			return
		default:
			if b.Screen == nil {
				return
			}
			ev, ok := <-b.Screen.EventQ()
			if !ok {
				return
			}
			if _, isInterrupt := ev.(*tcell.EventInterrupt); isInterrupt {
				return
			}
			b.processEvent(ev, events, stop)
		}
	}
}
func (b *Backend) processEvent(ev tcell.Event, events chan<- Event, stop <-chan struct{}) {
	var converted Event
	switch ev := ev.(type) {
	case *tcell.EventKey:
		converted = convertTcellKeyEvent(ev)
	case *tcell.EventMouse:
		converted = convertTcellMouseEvent(ev)
	case *tcell.EventResize:
		w, h := ev.Size()
		converted = Event{
			Type: ResizeEvent,
			ID:   "<Resize>",
			Payload: Resize{
				Width:  w,
				Height: h,
			},
		}
	default:
		return
	}
	select {
	case events <- converted:
	case <-stop:
	}
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
	tcell.KeyCtrlH:      "<C-h>",
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
		s := e.Str()
		if e.Modifiers()&tcell.ModAlt != 0 {
			ID = fmt.Sprintf("<M-%s>", s)
		} else {
			ID = s
		}
	} else {
		if val, ok := keyMap[e.Key()]; ok {
			ID = val
		} else {
			ID = fmt.Sprintf("<Key:%v>", e.Key())
		}
	}
	return Event{
		Type:    KeyboardEvent,
		ID:      ID,
		Payload: e,
	}
}
func convertTcellMouseEvent(e *tcell.EventMouse) Event {
	btns := e.Buttons()
	ID := "Unknown_Mouse_Button"
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
	if btns == tcell.ButtonNone {
		ID = "<MouseRelease>"
	}
	x, y := e.Position()
	return Event{
		Type: MouseEvent,
		ID:   ID,
		Payload: Mouse{
			X:    x,
			Y:    y,
			Drag: false,
		},
	}
}
