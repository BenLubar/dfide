package raws // import "github.com/BenLubar/dfide/gui/raws"

import (
	"github.com/BenLubar/dfide/gui"
	"github.com/BenLubar/dfide/gui/file"
	"github.com/BenLubar/dfide/raws/language"
)

type visualEditor interface {
	control() gui.Control
}

func showEditor(visual visualEditor, f *file.File, content []byte) {
	w := gui.NewWindow(f.Name, 600, 400)
	tab := gui.NewTab()
	if visual != nil {
		tab.Append("Visual", visual.control())
	}
	// TODO: create text editor tab.Append("Text", ???)
	w.SetChild(tab)
	w.Show()
}

type errorEditor struct {
	err error
}

func (e *errorEditor) control() gui.Control {
	return gui.NewLabel("Error: " + e.err.Error())
}

type languageEditor struct {
	tags []language.Tag
	box  gui.Box
}

func (e *languageEditor) control() gui.Control {
	box3 := gui.NewVerticalBox()
	e.box = box3

	box1 := gui.NewVerticalBox()
	box1.SetPadded(true)
	box2 := gui.NewHorizontalBox()
	box2.SetPadded(true)
	box2.Append(gui.NewLabel("Add"), false)
	addWordButton := gui.NewButton("Word")
	addWordButton.OnClick(func() {
		w := new(language.Word)
		e.tags = append(e.tags, language.Tag{Word: w})
		box3.Append(e.addWord(w), false)
	})
	box2.Append(addWordButton, false)
	addSymbolButton := gui.NewButton("Symbol")
	addSymbolButton.OnClick(func() {
		s := new(language.Symbol)
		e.tags = append(e.tags, language.Tag{Symbol: s})
		box3.Append(e.addSymbol(s), false)
	})
	box2.Append(addSymbolButton, false)
	addTranslationButton := gui.NewButton("Translation")
	addTranslationButton.OnClick(func() {
		t := new(language.Translation)
		e.tags = append(e.tags, language.Tag{Translation: t})
		box3.Append(e.addTranslation(t), false)
	})
	box2.Append(addTranslationButton, false)
	box1.Append(box2, false)
	box3.SetScrollable(gui.ScrollableAlways)
	box1.Append(box3, true)
	for _, t := range e.tags {
		var c gui.Control
		if t.Symbol != nil {
			c = e.addSymbol(t.Symbol)
		} else if t.Translation != nil {
			c = e.addTranslation(t.Translation)
		} else {
			c = e.addWord(t.Word)
		}
		box3.Append(c, false)
	}
	return box1
}

func (e *languageEditor) deleteTag(c gui.Control, tag language.Tag) {
	for i, t := range e.tags {
		if t == tag {
			e.tags = append(e.tags[:i], e.tags[i+1:]...)
			e.box.RemoveAt(i)
			gui.Destroy(c)
			return
		}
	}
}

func (e *languageEditor) addSymbol(s *language.Symbol) gui.Control {
	_ = s.ID
	for _, w := range s.Words {
		_ = w
	}

	// TODO
	return gui.NewLabel("TODO: symbol editor")
}

func (e *languageEditor) addTranslation(t *language.Translation) gui.Control {
	_ = t.ID
	for _, w := range t.Words {
		_ = w.English
		_ = w.Native
	}

	// TODO
	return gui.NewLabel("TODO: translation editor")
}

func (e *languageEditor) addWord(w *language.Word) gui.Control {
	vbox := gui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := gui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)

	idEntry := gui.NewEntry()
	idEntry.SetText(w.ID)
	idEntry.OnChange(func() {
		w.ID = idEntry.Text()
	})
	hbox.Append(idEntry, true)

	deleteButton := gui.NewButton("Delete")
	deleteButton.OnClick(func() {
		e.deleteTag(vbox, language.Tag{Word: w})
	})
	hbox.Append(deleteButton, false)

	tab := gui.NewTab()
	tab.Append("Noun", e.addNoun(&w.Noun))
	tab.Append("Adjective", e.addAdjective(&w.Adjective))
	tab.Append("Prefix", e.addPrefix(&w.Prefix))
	tab.Append("Verb", e.addVerb(&w.Verb))
	vbox.Append(tab, true)

	return vbox
}

func (e *languageEditor) addNoun(n **language.Noun) gui.Control {
	// TODO
	return gui.NewLabel("TODO: noun editor")
}

func (e *languageEditor) addAdjective(a **language.Adjective) gui.Control {
	// TODO
	return gui.NewLabel("TODO: adjective editor")
}

func (e *languageEditor) addPrefix(p **language.Prefix) gui.Control {
	// TODO
	return gui.NewLabel("TODO: prefix editor")
}

func (e *languageEditor) addVerb(v **language.Verb) gui.Control {
	// TODO
	return gui.NewLabel("TODO: verb editor")
}
