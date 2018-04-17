package file // import "github.com/BenLubar/dfide/gui/file"

import (
	"time"

	"github.com/BenLubar/dfide/gui"
)

type File struct {
	Name         string
	ContentType  string
	LastModified time.Time
	Contents     []byte
}

func Open(w gui.Window, f func(*File), accept ...string) {
	open(w, f, accept...)
}

func Save(w gui.Window, f *File) {
	save(w, f)
}
