package main // import "github.com/BenLubar/dfide"

import (
	"time"

	"github.com/BenLubar/dfide/gui"
	"github.com/BenLubar/dfide/gui/file"
	"github.com/BenLubar/dfide/gui/raws"
)

func main() {
	if err := gui.Main(uiMain); err != nil {
		panic(err)
	}
}

func uiMain() {
	vbox := gui.NewVerticalBox()
	vbox.SetPadded(true)
	vbox.SetScrollable(gui.ScrollableOverflow)

	openButton := gui.NewButton("Open Raw File")
	openButton.OnClick(func() {
		file.Open(gui.MainWindow, raws.OpenFile, ".txt")
	})
	vbox.Append(openButton, false)

	newLanguageButton := gui.NewButton("New Language Raws")
	newLanguageButton.OnClick(func() {
		raws.OpenFile(&file.File{
			Name:         "language_untitled.txt",
			ContentType:  "text/plain",
			LastModified: time.Now(),
			Contents:     []byte("language_untitled\n\n[OBJECT:LANGUAGE]\n"),
		})
	})
	vbox.Append(newLanguageButton, false)

	gui.MainWindow.SetChild(vbox)
	gui.MainWindow.SetMargined(true)
	gui.MainWindow.Show()
}
