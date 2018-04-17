// +build js

package gui // import "github.com/BenLubar/dfide/gui"

import (
	"fmt"
	"strings"

	"github.com/gopherjs/gopherjs/js"

	"honnef.co/go/js/dom"
)

var doc = dom.GetWindow().Document().(dom.HTMLDocument)

func main(f func()) error {
	f()

	doc.QuerySelector(".loading").Class().Add("hidden")

	return nil
}

func queueMain(f func()) {
	f()
}

type jsWindow struct {
	win   *dom.HTMLDivElement
	body  *dom.HTMLDivElement
	title *dom.HTMLButtonElement
	close *dom.HTMLButtonElement
}

func newWindow(title string, width, height int) *jsWindow {
	win := doc.CreateElement("div").(*dom.HTMLDivElement)
	win.Class().Add("window")
	win.Class().Add("hidden")
	titleBar := doc.CreateElement("button").(*dom.HTMLButtonElement)
	titleBar.Class().Add("title")
	titleBar.SetTextContent(title)
	win.AppendChild(titleBar)
	close := doc.CreateElement("button").(*dom.HTMLButtonElement)
	close.Class().Add("close-button")
	win.AppendChild(close)
	body := doc.CreateElement("div").(*dom.HTMLDivElement)
	body.Class().Add("body")
	body.Style().Set("width", fmt.Sprintf("%dpx", width))
	body.Style().Set("height", fmt.Sprintf("%dpx", height))
	win.AppendChild(body)
	doc.Body().AppendChild(win)
	w := &jsWindow{win: win, body: body, title: titleBar, close: close}
	w.init()
	return w
}

func (w *jsWindow) init() {
	w.win.AddEventListener("focus", true, func(dom.Event) {
		if w.win.Style().Get("z-index").Int() == 1 {
			return
		}

		for _, el := range doc.QuerySelectorAll(".window") {
			win := el.(*dom.HTMLDivElement)
			win.Style().Set("z-index", win.Style().Get("z-index").Int()-1)
		}
		w.win.Style().Set("z-index", 1)
	})
	w.win.AddEventListener("mousemove", false, func(event dom.Event) {
		e := event.(*dom.MouseEvent)
		if !e.Target().IsEqualNode(w.win) {
			w.win.Style().Delete("cursor")
			return
		}
		offsetX := e.Get("offsetX").Int()
		offsetY := e.Get("offsetY").Int()
		if offsetX <= 0 {
			if offsetY <= 0 {
				w.win.Style().Set("cursor", "nw-resize")
			} else if offsetY >= int(w.body.OffsetTop()+w.body.OffsetHeight()) {
				w.win.Style().Set("cursor", "sw-resize")
			} else {
				w.win.Style().Set("cursor", "w-resize")
			}
		} else if offsetX >= int(w.body.OffsetWidth()) {
			if offsetY <= 0 {
				w.win.Style().Set("cursor", "ne-resize")
			} else if offsetY >= int(w.body.OffsetTop()+w.body.OffsetHeight()) {
				w.win.Style().Set("cursor", "se-resize")
			} else {
				w.win.Style().Set("cursor", "e-resize")
			}
		} else {
			if offsetY <= 0 {
				w.win.Style().Set("cursor", "n-resize")
			} else {
				w.win.Style().Set("cursor", "s-resize")
			}
		}
	})
	w.win.AddEventListener("mousedown", false, func(event dom.Event) {
		e := event.(*dom.MouseEvent)
		if !e.Target().IsEqualNode(e.CurrentTarget()) {
			return
		}
		e.PreventDefault()
		resize := w.win.Style().Get("cursor").String()
		w.startMouseAction(resize, e.ClientX, e.ClientY)
	})
	w.win.AddEventListener("click", false, func(event dom.Event) {
		if doc.ActiveElement().CompareDocumentPosition(w.win)&dom.DocumentPositionContainedBy == 0 {
			w.title.Focus()
		}
	})
	w.title.AddEventListener("mousedown", false, func(event dom.Event) {
		e := event.(*dom.MouseEvent)
		w.startMouseAction("move", e.ClientX, e.ClientY)
	})
	w.close.AddEventListener("click", false, func(dom.Event) {
		js.Debugger()
	})
}

func (w *jsWindow) startMouseAction(action string, startX, startY int) {
	go func() {
		close := make(chan struct{}, 2)
		cover := doc.CreateElement("div").(*dom.HTMLDivElement)
		cover.Class().Add("cover")
		titleCursor := dom.GetWindow().GetComputedStyle(w.title, "").GetPropertyValue("cursor")
		if action == "move" && strings.HasSuffix(titleCursor, "grab") {
			cover.Style().Set("cursor", titleCursor+"bing")
		} else {
			cover.Style().Set("cursor", action)
		}

		cover.AddEventListener("mousemove", false, func(event dom.Event) {
			e := event.(*dom.MouseEvent)
			e.StopPropagation()
			e.PreventDefault()

			dx := e.ClientX - startX
			dy := e.ClientY - startY
			startX, startY = e.ClientX, e.ClientY
			w.mouseAction(action, dx, dy)
		})

		cover.AddEventListener("mouseup", false, func(event dom.Event) {
			e := event.(*dom.MouseEvent)
			e.StopPropagation()
			e.PreventDefault()

			close <- struct{}{}
		})

		reg := dom.GetWindow().AddEventListener("blur", false, func(dom.Event) {
			close <- struct{}{}
		})
		defer dom.GetWindow().RemoveEventListener("blur", false, reg)

		doc.Body().AppendChild(cover)
		defer doc.Body().RemoveChild(cover)

		<-close
	}()
}

func (w *jsWindow) mouseAction(action string, dx, dy int) {
	switch action {
	case "move":
		w.win.Style().Set("left", fmt.Sprintf("%dpx", int(w.win.OffsetLeft())+dx))
		w.win.Style().Set("top", fmt.Sprintf("%dpx", int(w.win.OffsetTop())+dy))
	case "n-resize":
		w.win.Style().Set("top", fmt.Sprintf("%dpx", int(w.win.OffsetTop())+dy))
		w.body.Style().Set("height", fmt.Sprintf("%dpx", int(w.body.OffsetHeight())-dy))
	case "s-resize":
		w.body.Style().Set("height", fmt.Sprintf("%dpx", int(w.body.OffsetHeight())+dy))
	case "w-resize":
		w.win.Style().Set("left", fmt.Sprintf("%dpx", int(w.win.OffsetLeft())+dx))
		w.body.Style().Set("width", fmt.Sprintf("%dpx", int(w.body.OffsetWidth())-dx))
	case "e-resize":
		w.body.Style().Set("width", fmt.Sprintf("%dpx", int(w.body.OffsetWidth())+dx))
	case "nw-resize":
		w.mouseAction("n-resize", dx, dy)
		w.mouseAction("w-resize", dx, dy)
	case "ne-resize":
		w.mouseAction("n-resize", dx, dy)
		w.mouseAction("e-resize", dx, dy)
	case "sw-resize":
		w.mouseAction("s-resize", dx, dy)
		w.mouseAction("w-resize", dx, dy)
	case "se-resize":
		w.mouseAction("s-resize", dx, dy)
		w.mouseAction("e-resize", dx, dy)
	default:
		js.Debugger()
		panic("gui: unhandled mouse action: " + action)
	}
}

func (w *jsWindow) SetChild(c Control) {
	for {
		child := w.body.FirstChild()
		if child == nil {
			break
		}
		w.body.RemoveChild(child)
	}

	w.body.AppendChild(c.control())
}

func (w *jsWindow) Show() {
	w.win.Class().Remove("hidden")
	w.title.Focus()
}

func (w *jsWindow) SetMargined(margined bool) {
	if margined {
		w.win.Class().Add("margined")
	} else {
		w.win.Class().Remove("margined")
	}
}

type control dom.HTMLElement

func destroy(c control) {
	// no-op
}

func newButton(text string) *jsButton {
	btn := doc.CreateElement("button").(*dom.HTMLButtonElement)
	btn.SetTextContent(text)
	return &jsButton{btn}
}

type jsButton struct {
	btn *dom.HTMLButtonElement
}

func (b *jsButton) OnClick(f func()) {
	b.btn.AddEventListener("click", false, func(dom.Event) {
		f()
	})
}

func (b *jsButton) control() control {
	return b.btn
}

var tabIndex int

func newTab() *jsTab {
	name := fmt.Sprintf("tabs%d", tabIndex)
	tabIndex++
	form := doc.CreateElement("form").(*dom.HTMLFormElement)
	form.Class().Add("tab-container")
	form.AddEventListener("submit", false, func(e dom.Event) {
		e.PreventDefault()
	})
	tabs := doc.CreateElement("ul").(*dom.HTMLUListElement)
	tabs.Class().Add("tabs")
	form.AppendChild(tabs)
	content := doc.CreateElement("div").(*dom.HTMLDivElement)
	content.Class().Add("content")
	form.AppendChild(content)
	return &jsTab{name: name, form: form, tabs: tabs, content: content}
}

type jsTab struct {
	name    string
	form    *dom.HTMLFormElement
	tabs    *dom.HTMLUListElement
	content *dom.HTMLDivElement
}

func (t *jsTab) Append(name string, child Control) {
	id := fmt.Sprintf("tab%d", tabIndex)
	tabIndex++

	li := doc.CreateElement("li").(*dom.HTMLLIElement)
	label := doc.CreateElement("label").(*dom.HTMLLabelElement)
	label.SetTextContent(name)
	label.For = id
	label.SetTabIndex(0)
	li.AppendChild(label)
	t.tabs.AppendChild(li)

	radio := doc.CreateElement("input").(*dom.HTMLInputElement)
	radio.Class().Add("hidden")
	radio.Type = "radio"
	radio.Name = t.name
	radio.SetID(id)
	radio.AddEventListener("change", false, func(dom.Event) {
		for _, el := range t.tabs.QuerySelectorAll("li>label") {
			el.Class().Remove("active")
		}
		label.Class().Add("active")
	})
	if t.content.FirstChild() == nil {
		label.Class().Add("active")
		radio.Checked = true
	}
	t.content.AppendChild(radio)

	t.content.AppendChild(child.control())
}

func (t *jsTab) control() control {
	return t.form
}

func newHorizontalBox() *jsBox {
	return newJSBox("horizontal")
}

func newVerticalBox() *jsBox {
	return newJSBox("vertical")
}

func newJSBox(direction string) *jsBox {
	box := doc.CreateElement("div").(*dom.HTMLDivElement)
	box.Class().Add(direction)
	box.Class().Add("box")
	return &jsBox{box: box}
}

type jsBox struct {
	box *dom.HTMLDivElement
}

func (b *jsBox) Append(child Control, stretchy bool) {
	if stretchy {
		child.control().Class().Add("stretchy")
	}

	b.box.AppendChild(child.control())
}

func (b *jsBox) RemoveAt(n int) {
	child := b.box.ChildNodes()[n]
	b.box.RemoveChild(child)
	child.(dom.HTMLElement).Class().Remove("stretchy")
}

func (b *jsBox) SetPadded(padded bool) {
	if padded {
		b.box.Class().Add("padded")
	} else {
		b.box.Class().Remove("padded")
	}
}

func (b *jsBox) SetScrollable(scrollable Scrollable) {
	b.box.Class().Remove("scroll-overflow")
	b.box.Class().Remove("scroll-always")
	switch scrollable {
	case ScrollableNever:
		break
	case ScrollableOverflow:
		b.box.Class().Add("scroll-overflow")
	case ScrollableAlways:
		b.box.Class().Add("scroll-always")
	}
}

func (b *jsBox) control() control {
	return b.box
}

func newLabel(text string) *jsLabel {
	label := doc.CreateElement("span").(*dom.HTMLSpanElement)
	label.SetTextContent(text)
	label.Class().Add("label")
	return &jsLabel{label: label}
}

type jsLabel struct {
	label *dom.HTMLSpanElement
}

func (l *jsLabel) control() control {
	return l.label
}

func newEntry() *jsEntry {
	input := doc.CreateElement("input").(*dom.HTMLInputElement)
	input.Type = "text"
	input.Class().Add("entry")
	return &jsEntry{input: input}
}

type jsEntry struct {
	input *dom.HTMLInputElement
}

func (e *jsEntry) Text() string {
	return e.input.Value
}

func (e *jsEntry) SetText(text string) {
	e.input.Value = text
}

func (e *jsEntry) OnChange(f func()) {
	lastValue := e.input.Value
	listener := func(dom.Event) {
		if e.input.Value == lastValue {
			return
		}
		lastValue = e.input.Value
		f()
	}
	e.input.AddEventListener("input", false, listener)
	e.input.AddEventListener("change", false, listener)
}

func (e *jsEntry) control() control {
	return e.input
}
