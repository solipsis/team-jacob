package main

import (
	"strconv"
	"strings"

	ui "github.com/gizak/termui"
)

type SetupScreen struct {
	amtEntry  *ui.Par
	addrEntry *ui.Par
	retEntry  *ui.Par
	help      *ui.Par
	legend    *legend
	stats     *pairStats

	editing  bool
	selected int
	fields   []*ui.Par
}

type field struct {
	par    *ui.Par
	active bool
}

func newSetupScreen(precise bool, stats *pairStats) *SetupScreen {

	l := new(legend)
	l.entries = append(l.entries, entry{key: "Q", text: "Quit"})
	l.entries = append(l.entries, entry{key: "S", text: "Save these settings for future orders"})
	// cycle between receive address, return address (optional), amount (if precise)
	// use arrow keys to select
	// enter + border color change to start editing

	amtEntry := ui.NewPar("")
	amtEntry.BorderLabel = "Amount"
	amtEntry.SetX(20)
	amtEntry.SetY(12)
	amtEntry.Width = 60
	amtEntry.Height = 3
	amtEntry.BorderFg = ui.ColorYellow

	addrEntry := ui.NewPar("")
	addrEntry.BorderLabel = "Receive Address"
	addrEntry.SetX(20)
	addrEntry.SetY(17)
	addrEntry.Width = 60
	addrEntry.Height = 3

	retEntry := ui.NewPar("")
	retEntry.BorderLabel = "Return Address (optional)"
	retEntry.SetX(20)
	retEntry.SetY(22)
	retEntry.Width = 60
	retEntry.Height = 3

	fields := []*ui.Par{amtEntry, addrEntry, retEntry}

	help := ui.NewPar(" Use the <arrow keys> to select a field and <space> to edit/confirm.  Press <enter> to confirm your order")
	help.SetX(15)
	help.SetY(26)
	help.Height = 4
	help.Width = 71

	return &SetupScreen{
		amtEntry:  amtEntry,
		addrEntry: addrEntry,
		retEntry:  retEntry,
		stats:     stats,
		fields:    fields,
		help:      help,
		legend:    l,
	}
}

func (s *SetupScreen) receiveAddress() string {
	return s.addrEntry.Text
}

func (s *SetupScreen) returnAddress() string {
	return s.retEntry.Text
}

func (s *SetupScreen) amount() (float64, error) {
	return strconv.ParseFloat(s.amtEntry.Text, 64)
}

func (s *SetupScreen) Buffers() []ui.Bufferer {
	bufs := []ui.Bufferer{}
	bufs = append(bufs, s.help)
	for _, f := range s.fields {
		bufs = append(bufs, f)
	}
	bufs = append(bufs, s.legend.Buffers()...)
	bufs = append(bufs, s.stats.Buffers()...)
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
	current.BorderBg = ui.ColorBlack
}

func (s *SetupScreen) deactivate() {
	s.editing = false
	current := s.fields[s.selected]
	current.BorderFg = ui.ColorYellow
	current.BorderBg = ui.ColorDefault
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
	current := s.fields[s.selected]
	if strings.HasSuffix(e, "<backspace>") || strings.HasSuffix(e, "<delete>") || strings.HasSuffix(e, "C-8") {
		if len(current.Text) > 0 {
			current.Text = current.Text[:len(current.Text)-1]
		}
		return
	}

	if s.editing {
		arr := strings.Split(e, "/")
		if len(arr) < 4 || len(arr[3]) > 1 {
			return
		}

		// append the character to the text
		current.Text += arr[3]
	}

}
