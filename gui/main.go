package gui // import "github.com/BenLubar/dfide/gui"

import "github.com/andlabs/ui"

// Main is the entry point of the dfide GUI. If the error returned is non-nil,
// there was a problem starting the program. This function will only return
// once the IDE is ready to exit.
func Main() error {
	return ui.Main(main)
}

func main() {
	// currently, this is just the example from the andlabs/ui wiki.

	input := ui.NewEntry()
	button := ui.NewButton("Greet")
	greeting := ui.NewLabel("")
	box := ui.NewVerticalBox()
	box.Append(ui.NewLabel("Enter your name:"), false)
	box.Append(input, false)
	box.Append(button, false)
	box.Append(greeting, false)
	window := ui.NewWindow("Hello", 200, 100, false)
	window.SetMargined(true)
	window.SetChild(box)
	button.OnClicked(func(*ui.Button) {
		greeting.SetText("Hello, " + input.Text() + "!")
	})
	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	window.Show()
}
