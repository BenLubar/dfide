// +build !js

package gui // import "github.com/BenLubar/dfide/gui"

import "github.com/andlabs/ui"

func main(f func()) error {
	return ui.Main(f)
}

func queueMain(f func()) {
	ui.QueueMain(f)
}

func newWindow(title string, width, height int) *uiWindow {
	return &uiWindow{ui.NewWindow(title, width, height, false)}
}

type uiWindow struct {
	win *ui.Window
}

func (w *uiWindow) SetChild(child Control) {
	w.win.SetChild(child.control())
}

func (w *uiWindow) Show() {
	w.win.Show()
}

func (w *uiWindow) SetMargined(margined bool) {
	w.win.SetMargined(margined)
}

func (w *uiWindow) UIWindow() *ui.Window { return w.win }

type control ui.Control

func destroy(c control) {
	c.Destroy()
}

func newButton(text string) *uiButton {
	return &uiButton{ui.NewButton(text)}
}

type uiButton struct {
	btn *ui.Button
}

func (b *uiButton) OnClick(f func()) {
	b.btn.OnClicked(func(*ui.Button) {
		f()
	})
}

func (b *uiButton) control() control {
	return b.btn
}

type uiTab struct {
	tab *ui.Tab
}

func newTab() *uiTab {
	return &uiTab{ui.NewTab()}
}

func (t *uiTab) Append(name string, child Control) {
	t.tab.Append(name, child.control())
}

func (t *uiTab) control() control {
	return t.tab
}

func newHorizontalBox() *uiBox {
	return &uiBox{ui.NewHorizontalBox()}
}

func newVerticalBox() *uiBox {
	return &uiBox{ui.NewVerticalBox()}
}

type uiBox struct {
	box *ui.Box
}

func (b *uiBox) Append(child Control, stretchy bool) {
	b.box.Append(child.control(), stretchy)
}

func (b *uiBox) RemoveAt(n int) {
	b.box.Delete(n)
}

func (b *uiBox) SetPadded(padded bool) {
	b.box.SetPadded(padded)
}

func (b *uiBox) SetScrollable(scrollable Scrollable) {
	panic("gui: TODO: https://github.com/andlabs/libui/issues/178")
}

func (b *uiBox) control() control {
	return b.box
}

func newLabel(text string) *uiLabel {
	return &uiLabel{label: ui.NewLabel(text)}
}

type uiLabel struct {
	label *ui.Label
}

func (l *uiLabel) control() control {
	return l.label
}

func newEntry() *uiEntry {
	return &uiEntry{entry: ui.NewEntry()}
}

type uiEntry struct {
	entry *ui.Entry
}

func (e *uiEntry) Text() string {
	return e.entry.Text()
}

func (e *uiEntry) SetText(text string) {
	e.entry.SetText(text)
}

func (e *uiEntry) OnChange(f func()) {
	e.entry.OnChanged(func(*ui.Entry) {
		f()
	})
}

func (e *uiEntry) control() control {
	return e.entry
}
