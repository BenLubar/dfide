package raws // import "github.com/BenLubar/dfide/gui/raws"

import (
	"bytes"

	"golang.org/x/text/encoding/charmap"

	"github.com/BenLubar/dfide/gui/file"
	"github.com/BenLubar/dfide/raws"
)

var cp437Dec = charmap.CodePage437.NewDecoder()

func OpenFile(f *file.File) {
	contents, err := cp437Dec.Bytes(f.Contents)
	if err != nil {
		return
	}

	r := raws.NewReader(bytes.NewReader(contents))

	showEditor(getVisualEditor(r), f, contents)
}

func getVisualEditor(r *raws.Reader) visualEditor {
	objectType, err := r.ObjectType()
	if err != nil {
		return &errorEditor{err: err}
	}

	switch objectType {
	case "LANGUAGE":
		l := &languageEditor{}
		if err := r.ParseAll(&l.tags); err != nil {
			return &errorEditor{err: err}
		}
		return l
	default:
		return nil
	}
}
