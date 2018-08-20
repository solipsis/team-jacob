package main

import (
	"math/rand"
	"strings"
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
	wheelX:     25,
	wheelWidth: 21,
	wheelY:     15,
}

type PairSelectorScreen struct {
	selector       *pairSelector
	stats          *pairStats
	typePar        *ui.Par
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

// TODO: edge cases
func centerText(p *ui.Par, t string) {
	l := p.Width - 2 // border lines
	Log.Println("l:", l, "width:", p.Width)
	pad := (l - len(t)) / 2
	p.Text = strings.Repeat(" ", pad) + t
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

	p.marketInfo = make(map[string]ss.MarketInfoResponse)
	for _, v := range rates {
		p.marketInfo[v.Pair] = v
	}
	n := toCoinRing(coins)

	pair := addPairSelector(n)
	formatSelector(pair, p.layout)
	p.selector = pair
	p.stats = newPairStats()

	typePar := ui.NewPar("")
	typePar.BorderLabel = "Order Type"
	typePar.BorderFg = ui.ColorMagenta
	typePar.Height = 3
	typePar.Width = 13
	centerText(typePar, "Quick")
	typePar.Y = p.layout.wheelY - 3
	typePar.X = p.layout.wheelX + 19
	p.typePar = typePar

	div := ui.NewPar(" < --- > ")
	div.Border = false
	div.Y = p.layout.wheelY + 5
	div.X = p.layout.wheelX + p.layout.wheelWidth
	div.Height = 3
	div.Width = 8
	p.divider = div

	l := new(legend)
	l.entries = append(l.entries, entry{key: "Q", text: "Quit"})
	l.entries = append(l.entries, entry{key: "T", text: "Toggle Order Type"})
	l.entries = append(l.entries, entry{key: "W", text: "Keepkey Mode (Disabled)"})
	l.entries = append(l.entries, entry{key: "Y", text: "I'm feeling Lucky"})
	p.legend = l

	msg := " Use <arrow keys> or <hjkl> to select 2 coins and <Enter> to initiate a Shift "
	help := ui.NewPar(msg)
	help.X = 10
	help.Height = 3
	help.Width = len(msg) + 2
	help.Y = 27
	p.help = help
}

func (p *PairSelectorScreen) activePair() string {
	return p.selector.receive.SelectedCoin().Symbol + "_" + p.selector.deposit.SelectedCoin().Symbol
}

func (p *PairSelectorScreen) isPreciseOrder() bool {
	return strings.Contains(p.typePar.Text, "Precise")
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
	info := p.marketInfo[d.Symbol+"_"+r.Symbol]
	p.stats.update(d.Symbol, r.Symbol, info.Min, info.Limit, info.Rate, info.MinerFee)
	bufs = append(bufs, p.stats.Buffers()...)
	bufs = append(bufs, p.divider)
	bufs = append(bufs, p.help)
	bufs = append(bufs, p.legend.Buffers()...)
	bufs = append(bufs, p.typePar)
	return bufs
}

// Handle responds to select UI events
func (s *PairSelectorScreen) Handle(e string) {
	Log.Println("Select Input", e)

	// Toggling orderType between quick and precise
	if e == "/sys/kbd/t" {
		if strings.Contains(s.typePar.Text, "Quick") {
			centerText(s.typePar, "Precise")
		} else {
			centerText(s.typePar, "Quick")
		}
	}

	// Deposit and recieve coin selection
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
		p.active.background.BorderFg = ui.ColorWhite
		p.active = p.deposit
		p.active.background.BorderFg = ui.ColorRed
	}

	// I'm feeling lucky toggle
	if e == "/sys/kbd/y" {
		if s.luckyTicker == nil {
			s.luckyTicker = time.NewTicker(100 * time.Millisecond)
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
		} else {
			s.luckyTicker.Stop()
			s.luckyTicker = nil
		}
	}
}
