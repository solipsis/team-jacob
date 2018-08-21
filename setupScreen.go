package main

import (
	"strconv"
	"strings"

	ui "github.com/gizak/termui"
)

// Screen for setting up order details
type SetupScreen struct {
	amtEntry  *ui.Par
	addrEntry *ui.Par
	retEntry  *ui.Par
	help      *ui.Par
	legend    *legend
	stats     *pairStats

	precise  bool
	editing  bool
	selected int
	fields   []*ui.Par
}

type field struct {
	par    *ui.Par
	active bool
}

func newSetupScreen(precise bool, stats *pairStats) *SetupScreen {
	Log.Println("NewSetupScreen")

	screen := &SetupScreen{}
	screen.fields = make([]*ui.Par, 0)
	screen.stats = stats
	screen.precise = precise

	l := new(legend)
	l.entries = append(l.entries, entry{key: "Q", text: "Quit"})
	l.entries = append(l.entries, entry{key: "S", text: "Save these settings for future orders"})
	screen.legend = l

	var c *setupConfig
	if precise {
		c = setupPreciseConfig
	} else {
		c = setupQuickConfig
	}

	// only show the amount field on precise orders
	if precise {
		amtEntry := ui.NewPar("")
		amtEntry.BorderLabel = "Amount"
		amtEntry.SetX(c.entryX)
		amtEntry.SetY(c.amtY)
		amtEntry.Width = c.entryWidth
		amtEntry.Height = c.entryHeight
		screen.amtEntry = amtEntry
		screen.fields = append(screen.fields, amtEntry)
	}

	addrEntry := ui.NewPar("")
	addrEntry.BorderLabel = "Receive Address"
	addrEntry.SetX(c.entryX)
	addrEntry.SetY(c.addrY)
	addrEntry.Width = c.entryWidth
	addrEntry.Height = c.entryHeight
	screen.addrEntry = addrEntry
	screen.fields = append(screen.fields, addrEntry)

	retEntry := ui.NewPar("")
	retEntry.BorderLabel = "Return Address (optional)"
	retEntry.SetX(c.entryX)
	retEntry.SetY(c.retY)
	retEntry.Width = c.entryWidth
	retEntry.Height = c.entryHeight
	screen.retEntry = retEntry
	screen.fields = append(screen.fields, retEntry)

	help := ui.NewPar(" Use the <arrow keys> to select a field and <space> to edit/confirm.  Press <enter> to confirm your order")
	help.SetX(c.helpX)
	help.SetY(c.helpY)
	help.Height = c.helpHeight
	help.Width = c.helpWidth
	screen.help = help

	screen.changeSelection(0)
	return screen
}

func (s *SetupScreen) receiveAddress() string {
	return s.addrEntry.Text
}

func (s *SetupScreen) returnAddress() string {
	return s.retEntry.Text
}

func (s *SetupScreen) amount() (float64, error) {
	if s.precise {
		return strconv.ParseFloat(s.amtEntry.Text, 64)
	}
	// return 0 if this is a quick order
	return 0, nil
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
	// visually deselect current selection
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

	// visually select new item
	current = s.fields[s.selected]
	current.BorderFg = ui.ColorYellow
}

// activate currently selected field for editing
func (s *SetupScreen) activate() {
	s.editing = true
	current := s.fields[s.selected]
	current.BorderFg = ui.ColorRed
	current.BorderBg = ui.ColorBlack
}

// stop editing currently selected field
func (s *SetupScreen) deactivate() {
	s.editing = false
	current := s.fields[s.selected]
	current.BorderFg = ui.ColorYellow
	current.BorderBg = ui.ColorDefault
}

func (s *SetupScreen) toggle() {
	if s.editing {
		s.deactivate()
	} else {
		s.activate()
	}
}

// Deal with user input on this screen
func (s *SetupScreen) Handle(e string) {

	// if arrow keys and no selected item, change the selected item
	if !s.editing {
		if e == "/sys/kbd/<up>" || e == "/sys/kbd/k" {
			s.changeSelection(-1)
			return
		}
		if e == "/sys/kbd/<down>" || e == "/sys/kbd/j" {
			s.changeSelection(1)
			return
		}
	}

	if e == "/sys/kbd/<space>" {
		s.toggle()
		return
	}
	current := s.fields[s.selected]

	if s.editing {
		// All the keys that could be used to "undo"
		if strings.HasSuffix(e, "<backspace>") || strings.HasSuffix(e, "<delete>") || strings.HasSuffix(e, "C-8") {
			if len(current.Text) > 0 {
				current.Text = current.Text[:len(current.Text)-1]
			}
			return
		}

		arr := strings.Split(e, "/")
		if len(arr) < 4 || len(arr[3]) > 1 {
			return
		}

		// append the character to the text
		current.Text += arr[3]
	} else {
		if e == "/sys/kbd/q" {
			ui.StopLoop()
		}
	}

}
