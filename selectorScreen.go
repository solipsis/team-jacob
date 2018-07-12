package main

import (
	"math/rand"
	"time"

	ui "github.com/gizak/termui"
	ss "github.com/solipsis/shapeshift"
)

type SelectLayout struct {
	infoX, infoY, infoHeight, infoWidth     int
	wheelX, wheelY, wheelHeight, wheelWidth int
}

var DefaultSelectLayout = &SelectLayout{
	infoX:      20,
	infoY:      40,
	infoHeight: 3,
	infoWidth:  40,
	wheelX:     37,
	wheelWidth: 21,
	wheelY:     12,
}

type PairSelectorScreen struct {
	selector       *pairSelector
	stats          *pairStats
	typeSelector   *ringSelector
	divider        *ui.Par
	info           *ui.Par
	help           *ui.Par
	legend         *legend
	layout         *SelectLayout
	luckyTicker    *time.Ticker
	marketInfo     map[string]ss.MarketInfoResponse
	jankDrawToggle bool
}

type pairSelector struct {
	deposit, receive, active *coinWheel
}

func NewPairSelectorScreen(l *SelectLayout) *PairSelectorScreen {
	return &PairSelectorScreen{layout: l}
}

type test struct {
	str string
}

func (t *test) Text() string {
	return t.str
}
func (p *PairSelectorScreen) Init() {
	// TODO: extract out all pairInfo code
	coins, err := activeCoins()
	if err != nil {
		Log.Println("Unable to contact shapeshift", err)
	}

	rates, err := ss.MarketInfo()
	if err != nil {
		Log.Println("Unable to get market info:", err)
	}

	m := make(map[string]ss.MarketInfoResponse)
	for _, v := range rates {
		m[v.Pair] = v
	}
	n := toCoinRing(coins)
	pair := addPairSelector(n)
	formatSelector(pair, p.layout)

	p.selector = pair
	p.marketInfo = m
	d, r := pair.deposit.ring.Value().Symbol, pair.receive.ring.Value().Symbol
	p.stats = NewPairStats(m[d+"_"+r])

	arr := []ringItem{&test{"Quick"}, &test{"Precise"}, &test{"I'm Feeling Lucky"}}
	p.typeSelector = NewRingSelector(arr, "Order Type", 6, 14, 3)

	div := ui.NewPar(" ----- > ")
	div.Border = false
	div.Y = p.layout.wheelY + 5
	div.X = p.layout.wheelX + p.layout.wheelWidth
	div.Height = 3
	div.Width = 8
	p.divider = div

	l := new(legend)
	l.entries = append(l.entries, entry{key: "Q", text: "Quit"})
	l.entries = append(l.entries, entry{key: "P", text: "Precise"})
	l.entries = append(l.entries, entry{key: "K", text: "Keepkey Mode"})
	p.legend = l

	msg := " Use <arrow keys> or <hjkl> to select 2 coins and <Enter> to initiate a Shift "
	help := ui.NewPar(msg)
	help.X = 7
	help.Height = 3
	help.Width = len(msg) + 2
	help.Y = 24
	p.help = help
}

func (p *PairSelectorScreen) activePair() string {
	return p.selector.receive.SelectedCoin().Symbol + "_" + p.selector.deposit.SelectedCoin().Symbol
}

func addPairSelector(r coinRing) *pairSelector {
	dep := NewCoinWheel(r, 7, "Deposit")
	rec := NewCoinWheel(r.Next(), 7, "Receive")
	dep.background.BorderLabelFg = ui.ColorGreen
	rec.background.BorderLabelFg = ui.ColorGreen
	rec.active.ItemFgColor = ui.ColorGreen

	return &pairSelector{dep, rec, dep}
}

func formatSelector(pair *pairSelector, layout *SelectLayout) {

	pair.deposit.active.X = layout.wheelX
	pair.deposit.active.Y = layout.wheelY + 4
	pair.deposit.background.X = layout.wheelX
	pair.deposit.background.Y = layout.wheelY
	pair.active.background.BorderFg = ui.ColorRed

	pair.receive.active.X = layout.wheelX + layout.wheelWidth + 9
	pair.receive.active.Y = layout.wheelY + 4
	pair.receive.background.X = layout.wheelX + layout.wheelWidth + 9
	pair.receive.background.Y = layout.wheelY
}

func (p *PairSelectorScreen) SelectedCoins() (dep, rec *Coin) {
	return p.selector.deposit.ring.Value(), p.selector.receive.ring.Value()
}

// TODO: remove dependency on ui???
func (p *pairSelector) Buffers() []ui.Bufferer {

	bufs := p.deposit.Buffers()
	bufs = append(bufs, p.receive.Buffers()...)
	return bufs
}

func (p *PairSelectorScreen) Buffers() []ui.Bufferer {
	bufs := p.selector.Buffers()
	// TODO: refactor this
	d, r := p.SelectedCoins()
	p.stats.Update(d.Symbol, r.Symbol, p.marketInfo[d.Symbol+"_"+r.Symbol])
	bufs = append(bufs, p.stats.Buffers()...)
	bufs = append(bufs, p.divider)
	bufs = append(bufs, p.help)
	bufs = append(bufs, p.legend.Buffers()...)
	bufs = append(bufs, p.typeSelector.Buffers()...)
	return bufs
}

// TODO: Remove this once I refactor coinwheels to use new list format
// TODO; list of menus. set active colors for active ones
// set inactive colors for others
var activeMenu = deposit

const (
	orderType int = iota
	deposit
	receive
)

// Handle responds to select UI events
func (s *PairSelectorScreen) Handle(e string) {
	Log.Println("Select Input", e)

	// TODO: completely redo this after wheel interface migration
	if activeMenu == deposit {
		p := s.selector
		if e == "/sys/kbd/<up>" || e == "/sys/kbd/k" {
			p.active.Prev()
		}
		if e == "/sys/kbd/<down>" || e == "/sys/kbd/j" {
			p.active.Next()
		}
		if e == "/sys/kbd/<right>" || e == "/sys/kbd/l" {
			activeMenu = receive
			p.active.background.BorderFg = ui.ColorWhite
			p.active = p.receive
			p.active.background.BorderFg = ui.ColorRed
		}
		if e == "/sys/kbd/<left>" || e == "/sys/kbd/h" {
			activeMenu = orderType
			p.active.background.BorderFg = ui.ColorWhite
			s.typeSelector.background.BorderFg = ui.ColorRed
			//p.active.background.BorderFg = ui.ColorWhite
			//p.active = p.deposit
			//p.active.background.BorderFg = ui.ColorRed
		}
	} else if activeMenu == receive {
		p := s.selector
		if e == "/sys/kbd/<up>" || e == "/sys/kbd/k" {
			p.active.Prev()
		}
		if e == "/sys/kbd/<down>" || e == "/sys/kbd/j" {
			p.active.Next()
		}
		if e == "/sys/kbd/<right>" || e == "/sys/kbd/l" {

			p.active.background.BorderFg = ui.ColorWhite
			p.active = p.receive
			p.active.background.BorderFg = ui.ColorRed
		}
		if e == "/sys/kbd/<left>" || e == "/sys/kbd/h" {
			activeMenu = deposit
			p.active.background.BorderFg = ui.ColorWhite
			p.active = p.deposit
			p.active.background.BorderFg = ui.ColorRed
		}
	} else if activeMenu == orderType {
		p := s.selector
		if e == "/sys/kbd/<up>" || e == "/sys/kbd/k" {
			s.typeSelector.Prev()
		}
		if e == "/sys/kbd/<down>" || e == "/sys/kbd/j" {
			s.typeSelector.Next()
		}
		if e == "/sys/kbd/<right>" || e == "/sys/kbd/l" {
			activeMenu = deposit
			s.typeSelector.background.BorderFg = ui.ColorWhite
			p.active = p.deposit
			p.active.background.BorderFg = ui.ColorRed
		}

	}
	// Toggle ticker on and off. TODO: cleaner signalling mechanism
	if s.typeSelector.Selected().Text() == "I'm Feeling Lucky" && s.luckyTicker == nil {
		s.luckyTicker = time.NewTicker(100 * time.Millisecond) // I think this is a giant go-routine leak
		go func() {
			for range s.luckyTicker.C {
				for i := 0; i < rand.Intn(20); i++ {
					s.selector.receive.Next()
				}
				if !s.jankDrawToggle {
					ui.Render(selectScreen.Buffers()...)
				}
			}
		}()
	}
	if s.typeSelector.Selected().Text() != "I'm Feeling Lucky" && s.luckyTicker != nil {
		s.luckyTicker.Stop()
		s.luckyTicker = nil
	}
}
