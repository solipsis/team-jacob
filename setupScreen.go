package main

import (
	"strings"

	ui "github.com/gizak/termui"
)

type SetupScreen struct {
	//amtEntry  *ui.Par
	//addrEntry *ui.Par
	//retEntry  *ui.Par
	help   *ui.Par
	legend *legend

	editing  bool
	selected int
	fields   []*ui.Par
}

type field struct {
	par    *ui.Par
	active bool
}

func newSetupScreen(precise bool) *SetupScreen {

	l := new(legend)
	l.entries = append(l.entries, entry{key: "Q", text: "Quit"})
	l.entries = append(l.entries, entry{key: "S", text: "Save these settings for future orders"})
	// cycle between receive address, return address (optional), amount (if precise)
	// use arrow keys to select
	// enter + border color change to start editing

	amtEntry := ui.NewPar("")
	amtEntry.BorderLabel = "Amount"
	amtEntry.SetX(10)
	amtEntry.SetY(10)
	amtEntry.Width = 40
	amtEntry.Height = 3
	amtEntry.BorderFg = ui.ColorYellow

	addrEntry := ui.NewPar("")
	addrEntry.BorderLabel = "Receive Address"
	addrEntry.SetX(10)
	addrEntry.SetY(15)
	addrEntry.Width = 40
	addrEntry.Height = 3

	retEntry := ui.NewPar("")
	retEntry.BorderLabel = "Return Address (optional)"
	retEntry.SetX(10)
	retEntry.SetY(20)
	retEntry.Width = 40
	retEntry.Height = 3

	fields := []*ui.Par{amtEntry, addrEntry, retEntry}

	help := ui.NewPar("Use the <arrow keys> to select a field and <space> to edit/confirm. Press <enter> to confirm your order")
	help.SetX(5)
	help.SetY(25)
	help.Height = 3
	help.Width = 80

	return &SetupScreen{
		//amtEntry:  amtEntry,
		//addrEntry: addrEntry,
		//retEntry:  retEntry,
		fields: fields,
		help:   help,
		legend: l,
	}
}

func (s *SetupScreen) Buffers() []ui.Bufferer {
	//bufs := append([]ui.Bufferer{}, s.amtEntry)
	//bufs = append(bufs, s.addrEntry)
	//bufs = append(bufs, s.retEntry)
	bufs := []ui.Bufferer{}
	bufs = append(bufs, s.help)
	for _, f := range s.fields {
		bufs = append(bufs, f)
	}
	bufs = append(bufs, s.legend.Buffers()...)
	return bufs
}

func (s *SetupScreen) changeSelection(i int) {
	// deselect current selection
	current := s.fields[s.selected]
	current.BorderFg = ui.ColorWhite

	index := s.selected + i
	if index < 0 {
		index = 0
	}
	if index >= len(s.fields) {
		index = len(s.fields) - 1
	}
	s.selected = index

	current = s.fields[s.selected]
	current.BorderFg = ui.ColorYellow
}

func (s *SetupScreen) activate() {
	s.editing = true
	current := s.fields[s.selected]
	current.BorderFg = ui.ColorRed
	current.BorderBg = ui.ColorYellow
}

func (s *SetupScreen) deactivate() {
	s.editing = false
	current := s.fields[s.selected]
	current.BorderFg = ui.ColorYellow
	current.BorderBg = ui.ColorBlack
}

func (s *SetupScreen) Handle(e string) {

	// if arrow keys and no selected item
	//	change selected item
	if !s.editing {
		if e == "/sys/kbd/<up>" || e == "/sys/kbd/k" {
			s.changeSelection(-1)
			return
		}
		if e == "/sys/kbd/<down>" || e == "/sys/kbd/j" {
			s.changeSelection(1)
			return
		}
		if e == "/sys/kbd/<space>" {
			s.activate()
			return
		}
	} else {
		if e == "/sys/kbd/<space>" {
			s.deactivate()
			return
		}
	}
	// if enter and no selected item
	// 	change active item

	// if enter and selected item
	// 	stop editing item

	// if enter on confirm button
	//	transitio	n

	// All the keys that could be used to "undo"
	if strings.HasSuffix(e, "<backspace>") || strings.HasSuffix(e, "<delete>") || strings.HasSuffix(e, "C-8") {
		if len(s.fields[0].Text) > 0 {
			s.fields[0].Text = s.fields[0].Text[:len(s.fields[0].Text)-1]
		}
		return
	}

	if s.editing {
		arr := strings.Split(e, "/")
		if len(arr) < 4 || len(arr[3]) > 1 {
			return
		}

		// append the character to the text
		s.fields[0].Text += arr[3]
	}

}
