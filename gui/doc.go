package gui // import "github.com/BenLubar/dfide/gui"

func Main(f func()) error {
	return main(func() {
		MainWindow = newWindow("Dwarf Fortress Integrated Development Environment", 300, 200)

		f()
	})
}

func QueueMain(f func()) {
	queueMain(f)
}

type Window interface {
	SetChild(Control)
	Show()
	SetMargined(bool)
	SetTitle(string)
}

func NewWindow(title string, width, height int) Window {
	return newWindow(title, width, height)
}

var MainWindow Window

type Control interface {
	control() control
}

func Destroy(c Control) {
	destroy(c.control())
}

type Button interface {
	Control
	OnClick(func())
}

func NewButton(text string) Button {
	return newButton(text)
}

type Tab interface {
	Control
	Append(string, Control)
	InsertAt(string, int, Control)
	RemoveAt(int)
	SetMargined(int, bool)
}

func NewTab() Tab {
	return newTab()
}

type Scrollable int

const (
	ScrollableNever Scrollable = iota
	ScrollableOverflow
	ScrollableAlways
)

type Box interface {
	Control
	Append(Control, bool)
	RemoveAt(int)
	SetPadded(bool)
	SetScrollable(Scrollable)
}

func NewHorizontalBox() Box {
	return newHorizontalBox()
}

func NewVerticalBox() Box {
	return newVerticalBox()
}

type Label interface {
	Control
}

func NewLabel(text string) Label {
	return newLabel(text)
}

type Entry interface {
	Control
	Text() string
	SetText(string)
	OnChange(func())
}

func NewEntry() Entry {
	return newEntry()
}

func NewMultiLineEntry() Entry {
	return newMultiLineEntry()
}
