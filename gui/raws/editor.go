package raws // import "github.com/BenLubar/dfide/gui/raws"

import (
	"bytes"

	"github.com/BenLubar/dfide/gui"
	"github.com/BenLubar/dfide/raws"
)

func showEditor(visual visualEditor, name string, content []byte) {
	w := gui.NewWindow(name+".txt", 600, 400)
	tab := gui.NewTab()

	text := gui.NewMultiLineEntry()
	text.SetText(string(content))
	tab.Append("Text", text)
	tab.SetMargined(0, true)

	visualChange := func(newContent []byte) {
		content = newContent
		text.SetText(string(newContent))
	}

	var vcontrol gui.Control
	resetVisual := func() {
		if vcontrol != nil {
			tab.RemoveAt(0)
			gui.Destroy(vcontrol)
		}

		if visual != nil {
			visual.setName(name)
			vcontrol = visual.control()
			tab.InsertAt("Visual", 0, vcontrol)
			tab.SetMargined(0, true)
			visual.OnChange(visualChange)
		} else {
			vcontrol = nil
		}
	}
	resetVisual()

	text.OnChange(func() {
		content = []byte(text.Text())
		r := raws.NewReader(bytes.NewReader(content))
		newName, err := r.Name()
		if err == nil {
			name = newName
			w.SetTitle(name + ".txt")
		}
		visual = getVisualEditor(r)
		resetVisual()
	})

	w.SetChild(tab)
	w.Show()
}
