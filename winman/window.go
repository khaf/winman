package winman

import (
	"fmt"

	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
)

type Window struct {
	*cview.Box
	root          cview.Primitive // The item to be positioned. May be nil for an empty item.
	manager       *Manager
	buttons       []*Button
	border        bool
	restoreX      int
	restoreY      int
	restoreWidth  int
	restoreHeight int
	maximized     bool
	Draggable     bool
	Resizable     bool
	modal         bool
}

// NewWindow creates a new window in this window manager
func NewWindow() *Window {
	window := &Window{
		Box: cview.NewBox().SetBackgroundColor(tcell.ColorDefault),
	}
	window.restoreX, window.restoreY, window.restoreHeight, window.restoreWidth = window.GetRect()
	window.SetBorder(true)
	return window
}

func (w *Window) SetRoot(root cview.Primitive) *Window {
	w.root = root
	return w
}

func (w *Window) GetRoot() cview.Primitive {
	return w.root
}

func (w *Window) Draw(screen tcell.Screen) {
	if w.Box.HasFocus() && !w.HasFocus() {
		w.Box.Blur()
	}
	w.Box.Draw(screen)

	if w.root != nil {
		x, y, width, height := w.GetInnerRect()
		w.root.SetRect(x, y, width, height)
		w.root.Draw(cview.NewClipRegion(screen, x, y, width, height))
	}

	if w.border {
		x, y, width, height := w.GetRect()
		screen = cview.NewClipRegion(screen, x, y, width, height)
		for _, button := range w.buttons {
			buttonX, buttonY := button.offsetX+x, button.offsetY+y
			if button.offsetX < 0 {
				buttonX += width
			}
			if button.offsetY < 0 {
				buttonY += height
			}

			//screen.SetContent(buttonX, buttonY, button.Symbol, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
			cview.Print(screen, cview.Escape(fmt.Sprintf("[%c]", button.Symbol)), buttonX-1, buttonY, 9, 0, tcell.ColorYellow)
		}
	}
}

func (w *Window) checkManager() {
	if w.manager == nil {
		panic("Window must be added to a Window Manager to call this method")
	}
}

func (w *Window) Show() *Window {
	w.checkManager()
	w.manager.Show(w)
	return w
}

func (w *Window) Hide() *Window {
	w.checkManager()
	w.manager.Hide(w)
	return w
}

func (w *Window) Maximize() *Window {
	w.checkManager()
	w.restoreX, w.restoreY, w.restoreHeight, w.restoreWidth = w.GetRect()
	w.SetRect(w.manager.GetInnerRect())
	w.maximized = true
	return w
}

func (w *Window) Restore() *Window {
	w.SetRect(w.restoreX, w.restoreY, w.restoreHeight, w.restoreWidth)
	w.maximized = false
	return w
}

func (w *Window) ShowModal() *Window {
	w.checkManager()
	w.manager.ShowModal(w)
	return w
}

func (w *Window) Center() *Window {
	w.checkManager()
	mx, my, mw, mh := w.manager.GetInnerRect()
	x, y, width, height := w.GetRect()
	x = mx + (mw-width)/2
	y = my + (mh-height)/2
	w.SetRect(x, y, width, height)
	return w
}

// Focus is called when this primitive receives focus.
func (w *Window) Focus(delegate func(p cview.Primitive)) {
	if w.root != nil {
		delegate(w.root)
		w.Box.Focus(nil)
	} else {
		delegate(w.Box)
	}
}

func (w *Window) Blur() {
	if w.root != nil {
		w.root.Blur()
	}
	w.Box.Blur()
}

func (w *Window) IsMaximized() bool {
	return w.maximized
}

// SetBorder sets the flag indicating whether or not the box should have a
// border.
func (w *Window) SetBorder(show bool) *Window {
	w.border = show
	w.Box.SetBorder(show)
	return w
}

// HasFocus returns whether or not this primitive has focus.
func (w *Window) HasFocus() bool {
	if w.root != nil {
		return w.root.GetFocusable().HasFocus()
	} else {
		return w.Box.HasFocus()
	}
}

func (w *Window) MouseHandler() func(action cview.MouseAction, event *tcell.EventMouse, setFocus func(p cview.Primitive)) (consumed bool, capture cview.Primitive) {
	return w.WrapMouseHandler(func(action cview.MouseAction, event *tcell.EventMouse, setFocus func(p cview.Primitive)) (consumed bool, capture cview.Primitive) {
		if action == cview.MouseLeftClick {
			x, y := event.Position()
			wx, wy, width, _ := w.GetRect()
			if y == wy {
				for _, button := range w.buttons {
					if button.offsetX >= 0 && x == wx+button.offsetX || button.offsetX < 0 && x == wx+width+button.offsetX {
						if button.ClickHandler != nil {
							button.ClickHandler()
						}
						return true, nil
					}
				}
			}
		}
		if w.root != nil {
			return w.root.MouseHandler()(action, event, setFocus)
		}
		return false, nil
	})
}

func (w *Window) AddButton(button *Button) *Window {
	w.buttons = append(w.buttons, button)

	offsetLeft, offsetRight := 2, -3
	for _, button := range w.buttons {
		if button.Alignment == ButtonRight {
			button.offsetX = offsetRight
			offsetRight -= 3
		} else {
			button.offsetX = offsetLeft
			offsetLeft += 3
		}
	}

	return w
}

func (w *Window) GetButton(i int) *Button {
	if i < 0 || i >= len(w.buttons) {
		return nil
	}
	return w.buttons[i]
}

func (w *Window) ButtonCount() int {
	return len(w.buttons)
}
