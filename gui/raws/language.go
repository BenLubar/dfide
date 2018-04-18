package raws // import "github.com/BenLubar/dfide/gui/raws"

import (
	"github.com/BenLubar/dfide/gui"
	"github.com/BenLubar/dfide/raws"
	"github.com/BenLubar/dfide/raws/language"
)

type languageEditor struct {
	baseVisualEditor
	tags []language.Tag
	box  gui.Box
}

func (e *languageEditor) control() gui.Control {
	e.box = gui.NewVerticalBox()
	e.box.SetPadded(true)
	e.box.SetScrollable(gui.ScrollableAlways)

	vbox := gui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox := gui.NewHorizontalBox()
	hbox.SetPadded(true)
	hbox.Append(gui.NewLabel(""), true) // center-align
	hbox.Append(gui.NewLabel("Add"), false)
	addWordButton := gui.NewButton("Word")
	addWordButton.OnClick(func() {
		w := new(language.Word)
		e.tags = append(e.tags, language.Tag{Word: w})
		e.box.Append(e.addWord(w), false)
		e.onChange()
	})
	hbox.Append(addWordButton, false)
	addSymbolButton := gui.NewButton("Symbol")
	addSymbolButton.OnClick(func() {
		s := new(language.Symbol)
		e.tags = append(e.tags, language.Tag{Symbol: s})
		e.box.Append(e.addSymbol(s), false)
		e.onChange()
	})
	hbox.Append(addSymbolButton, false)
	addTranslationButton := gui.NewButton("Translation")
	addTranslationButton.OnClick(func() {
		t := new(language.Translation)
		e.tags = append(e.tags, language.Tag{Translation: t})
		e.box.Append(e.addTranslation(t), false)
		e.onChange()
	})
	hbox.Append(addTranslationButton, false)
	hbox.Append(gui.NewLabel(""), true) // center-align
	vbox.Append(hbox, false)
	vbox.Append(e.box, true)
	for _, t := range e.tags {
		var c gui.Control
		if t.Symbol != nil {
			c = e.addSymbol(t.Symbol)
		} else if t.Translation != nil {
			c = e.addTranslation(t.Translation)
		} else {
			c = e.addWord(t.Word)
		}
		e.box.Append(c, false)
	}
	return vbox
}

func (e *languageEditor) deleteTag(c gui.Control, tag language.Tag) {
	for i, t := range e.tags {
		if t == tag {
			e.tags = append(e.tags[:i], e.tags[i+1:]...)
			e.box.RemoveAt(i)
			gui.Destroy(c)
			e.onChange()
			return
		}
	}
}

func (e *languageEditor) addTagBase(label string, id *string, tag language.Tag) gui.Box {
	vbox := gui.NewVerticalBox()

	hbox := gui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)

	hbox.Append(gui.NewLabel(label), false)

	idEntry := gui.NewEntry()
	idEntry.SetText(*id)
	idEntry.OnChange(func() {
		*id = idEntry.Text()
		e.onChange()
	})
	hbox.Append(idEntry, true)

	deleteButton := gui.NewButton("Delete")
	deleteButton.OnClick(func() {
		e.deleteTag(vbox, tag)
	})
	hbox.Append(deleteButton, false)

	return vbox
}

func (e *languageEditor) addSymbol(s *language.Symbol) gui.Control {
	vbox := e.addTagBase("Symbol:", &s.ID, language.Tag{Symbol: s})

	wordEntries := make([]gui.Entry, len(s.Words)+1)

	var entryOnChange func(entry gui.Entry) func()
	entryOnChange = func(entry gui.Entry) func() {
		return func() {
			if entry == wordEntries[len(wordEntries)-1] {
				if entry.Text() != "" {
					s.Words = append(s.Words, entry.Text())
					e.onChange()

					newDummy := gui.NewEntry()
					newDummy.OnChange(entryOnChange(newDummy))
					wordEntries = append(wordEntries, newDummy)
					vbox.Append(newDummy, false)
				}
			} else {
				for i, we := range wordEntries {
					if entry == we {
						if entry.Text() != "" {
							s.Words[i] = entry.Text()
						} else {
							s.Words = append(s.Words[:i], s.Words[i+1:]...)
							wordEntries = append(wordEntries[:i], wordEntries[i+1:]...)
							vbox.RemoveAt(i + 1)
							gui.Destroy(entry)
						}
						e.onChange()
						return
					}
				}
			}
		}
	}

	for i, w := range s.Words {
		entry := gui.NewEntry()
		entry.SetText(w)
		entry.OnChange(entryOnChange(entry))
		wordEntries[i] = entry
		vbox.Append(entry, false)
	}

	entry := gui.NewEntry()
	entry.OnChange(entryOnChange(entry))
	wordEntries[len(s.Words)] = entry
	vbox.Append(entry, false)

	return vbox
}

func (e *languageEditor) addTranslation(t *language.Translation) gui.Control {
	vbox := e.addTagBase("Translation:", &t.ID, language.Tag{Translation: t})

	var addEntry func(english, native string)

	wordEntries := make([]gui.Box, 0, len(t.Words)+1)

	var entryOnChange func(english, native gui.Entry, hbox gui.Box) func()
	entryOnChange = func(english, native gui.Entry, hbox gui.Box) func() {
		return func() {
			if hbox == wordEntries[len(wordEntries)-1] {
				if english.Text() != "" || native.Text() != "" {
					t.Words = append(t.Words, language.TWord{
						English: english.Text(),
						Native:  native.Text(),
					})
					e.onChange()
				}
			} else {
				for i, we := range wordEntries {
					if hbox == we {
						if english.Text() != "" || native.Text() != "" {
							t.Words[i].English = english.Text()
							t.Words[i].Native = native.Text()
						} else {
							t.Words = append(t.Words[:i], t.Words[i+1:]...)
							wordEntries = append(wordEntries[:i], wordEntries[i+1:]...)
							vbox.RemoveAt(i + 1)
							gui.Destroy(hbox)
						}
						e.onChange()
						return
					}
				}
			}
		}
	}

	addEntry = func(english, native string) {
		newEnglish := gui.NewEntry()
		newNative := gui.NewEntry()

		newEnglish.SetText(english)
		newNative.SetText(native)

		hbox := gui.NewHorizontalBox()
		hbox.SetPadded(true)
		hbox.Append(newEnglish, true)
		hbox.Append(newNative, true)

		wordEntries = append(wordEntries, hbox)

		change := entryOnChange(newEnglish, newNative, hbox)
		newEnglish.OnChange(change)
		newNative.OnChange(change)

		vbox.Append(hbox, false)
	}

	for _, w := range t.Words {
		addEntry(w.English, w.Native)
	}
	addEntry("", "")

	return vbox
}

func (e *languageEditor) addWord(w *language.Word) gui.Control {
	vbox := e.addTagBase("Word:", &w.ID, language.Tag{Word: w})

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

func (e *languageEditor) onChange() {
	e.baseVisualEditor.onChange("LANGUAGE", func(w *raws.Writer) error {
		return w.SerializeAll(e.tags)
	})
}
