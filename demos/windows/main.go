// Demo code for the Form primitive.
package main

import (
	"fmt"

	"gitlab.com/tslocum/cview"
)

func main() {

	app := cview.NewApplication()
	wm := cview.NewWindowManager()
	createForm := func(i int) {
		form := cview.NewForm().
			AddDropDown("Title", []string{"Mr.", "Ms.", "Mrs.", "Dr.", "Prof."}, 0, nil).
			AddInputField("First name", "", 20, nil, nil).
			AddInputField("Last name", "", 20, nil, nil).
			AddPasswordField("Password", "", 10, '*', nil).
			AddCheckBox("", "Age 18+", false, nil).
			AddButton("Save", nil).
			AddButton("Quit", func() {
				app.Stop()
			})

		window := wm.NewWindow(form, true)
		window.Draggable = i%2 == 0
		window.Resizable = i%3 == 0
		title := fmt.Sprintf("Window%d - Draggable: %t, Resizable: %t", i, window.Draggable, window.Resizable)
		window.SetBorder(true).SetTitle(title).SetTitleAlign(cview.AlignCenter)
		window.SetRect(2+i*2, 2+i, 50, 30)
		window.AddButton(&cview.WindowButton{
			Symbol:       'X',
			Alignment:    cview.AlignLeft,
			ClickHandler: func() { window.Close() },
		})
		var maxMinButton *cview.WindowButton
		maxMinButton = &cview.WindowButton{
			Symbol:    '▴',
			Alignment: cview.AlignRight,
			ClickHandler: func() {
				if window.IsMaximized() {
					window.Restore()
					maxMinButton.Symbol = '▴'
				} else {
					window.Maximize()
					maxMinButton.Symbol = '▾'
				}
			},
		}
		window.AddButton(maxMinButton)

		window.Open()
	}

	for i := 0; i < 10; i++ {
		createForm(i)
	}

	if err := app.SetRoot(wm, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
