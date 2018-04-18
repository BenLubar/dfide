// +build !js

package file // import "github.com/BenLubar/dfide/gui/file"

import (
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/BenLubar/dfide/gui"
	"github.com/andlabs/ui"
)

// TODO: handle errors

func uiWindow(w gui.Window) *ui.Window {
	return w.(interface{ UIWindow() *ui.Window }).UIWindow()
}

func open(w gui.Window, f func(*File), accept ...string) {
	// TODO: allow multiple files, check "accept"

	name := ui.OpenFile(uiWindow(w))
	if name == "" {
		return
	}

	fi, err := os.Stat(name)
	if err != nil {
		log.Println(err)
		return
	}

	content, err := ioutil.ReadFile(name)
	if err != nil {
		log.Println(err)
		return
	}

	f(&File{
		Name:         filepath.Base(name),
		ContentType:  mime.TypeByExtension(filepath.Ext(name)),
		LastModified: fi.ModTime(),
		Contents:     content,
	})
}

func save(w gui.Window, f *File) {
	name := ui.SaveFile(uiWindow(w))

	err := ioutil.WriteFile(name, f.Contents, 0644)
	if err != nil {
		log.Println(err)
		return
	}
}
