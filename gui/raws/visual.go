package raws // import "github.com/BenLubar/dfide/gui/raws"

import (
	"bytes"
	"log"

	"github.com/BenLubar/dfide/gui"
	"github.com/BenLubar/dfide/raws"
)

type visualEditor interface {
	control() gui.Control
	setName(string)
	OnChange(func([]byte))
}

type baseVisualEditor struct {
	listeners []func([]byte)
	name      string
}

func (e *baseVisualEditor) OnChange(f func([]byte)) {
	e.listeners = append(e.listeners, f)
}

func (e *baseVisualEditor) setName(name string) {
	e.name = name
}

func (e *baseVisualEditor) onChange(objectType string, write func(*raws.Writer) error) {
	if len(e.listeners) == 0 {
		return
	}

	// TODO: handle errors

	var buf bytes.Buffer

	w, err := raws.NewWriter(&buf, e.name, objectType)
	if err != nil {
		log.Println(err)
		return
	}

	err = write(w)
	if err != nil {
		log.Println(err)
		return
	}

	err = w.Flush()
	if err != nil {
		log.Println(err)
		return
	}

	for _, l := range e.listeners {
		l(buf.Bytes())
	}
}

type errorEditor struct {
	baseVisualEditor
	err error
}

func (e *errorEditor) control() gui.Control {
	return gui.NewLabel("Error: " + e.err.Error())
}
