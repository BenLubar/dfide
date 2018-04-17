package main // import "github.com/BenLubar/dfide"

import (
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
	openButton := gui.NewButton("Open Raw File")
	openButton.OnClick(func() {
		file.Open(gui.MainWindow, raws.OpenFile, ".txt")
	})
	gui.MainWindow.SetChild(openButton)
	gui.MainWindow.SetMargined(true)
	gui.MainWindow.Show()
}
