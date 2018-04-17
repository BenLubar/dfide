// +build js

package file // import "github.com/BenLubar/dfide/gui/file"

import (
	"strings"
	"time"

	"github.com/BenLubar/dfide/gui"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var doc = dom.GetWindow().Document().(dom.HTMLDocument)
var body = doc.Body().(*dom.HTMLBodyElement)

var fileInput = func() *dom.HTMLInputElement {
	input := doc.CreateElement("input").(*dom.HTMLInputElement)
	input.Type = "file"
	input.Multiple = true
	input.Class().Add("hidden")
	body.AppendChild(input)
	return input
}()

var lastFileCallback func(*js.Object)

func open(w gui.Window, f func(*File), accept ...string) {
	fileInput.RemoveEventListener("change", false, lastFileCallback)

	fileInput.Accept = strings.Join(accept, ",")

	lastFileCallback = fileInput.AddEventListener("change", false, func(dom.Event) {
		for _, file := range fileInput.Files() {
			lastModified := js.Global.Get("Date").New(file.Get("lastModified"))
			obj := &File{
				Name:         file.Get("name").String(),
				ContentType:  file.Get("type").String(),
				LastModified: lastModified.Interface().(time.Time).Round(time.Millisecond),
			}
			fileReader := js.Global.Get("FileReader").New()
			fileReader.Set("onload", func() {
				buffer := fileReader.Get("result")
				contents := js.Global.Get("Uint8Array").New(buffer)
				obj.Contents = contents.Interface().([]byte)
				go f(obj)
			})
			fileReader.Call("readAsArrayBuffer", file)
		}
		fileInput.Value = ""
	})

	fileInput.Click()
}

func save(w gui.Window, f *File) {
	blob := js.Global.Get("Blob").New(f.Contents, map[string]string{
		"type": f.ContentType,
	})

	url := js.Global.Get("URL").Call("createObjectURL", blob).String()
	defer js.Global.Get("URL").Call("revokeObjectURL", url)

	a := doc.CreateElement("a").(*dom.HTMLAnchorElement)
	a.Class().Add("hidden")
	a.SetAttribute("download", f.Name)
	a.Href = url

	body.AppendChild(a)
	defer body.RemoveChild(a)

	a.Click()
}
