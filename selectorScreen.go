package main

import (
	"log"

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
	wheelX:     21,
	wheelWidth: 21,
	wheelY:     12,
}

type PairSelectorScreen struct {
	selector   *pairSelector
	stats      *pairStats
	divider    *ui.Par
	info       *ui.Par
	help       *ui.Par
	layout     *SelectLayout
	marketInfo map[string]ss.MarketInfoResponse
}

type pairSelector struct {
	deposit, receive, active *coinWheel
}

func NewPairSelectorScreen(l *SelectLayout) *PairSelectorScreen {
	return &PairSelectorScreen{layout: l}
}

func (p *PairSelectorScreen) Init() {
	coins, err := activeCoins()
	if err != nil {
		log.Println("Unableto contact shapeshift")
	}

	rates, err := ss.MarketInfo()
	if err != nil {
		//fmt.Println(err)
	}

	m := make(map[string]ss.MarketInfoResponse)
	for _, v := range rates {
		m[v.Pair] = v
	}
	n := initWindow(coins)
	pair := addPairSelector(n)
	formatSelector(pair, p.layout)

	p.selector = pair
	p.marketInfo = m
	d, r := pair.deposit.ring.Value().Symbol, pair.receive.ring.Value().Symbol
	p.stats = NewPairStats(d, r, m[d+"_"+r])

	div := ui.NewPar(" < --- > ")
	div.Border = false
	div.Y = p.layout.wheelY + 5
	div.X = p.layout.wheelX + p.layout.wheelWidth
	div.Height = 3
	div.Width = 8
	p.divider = div

	msg := " Use <arrow keys> to select 2 coins and <Enter> to initiate a Shift "
	help := ui.NewPar(msg)
	help.X = 13
	help.Height = 3
	help.Width = len(msg) + 2
	help.Y = 24
	p.help = help
}

func addPairSelector(r coinRing) *pairSelector {
	dep := NewCoinWheel(r, 7, "Deposit")
	//positionWheel(dep, x, y)
	rec := NewCoinWheel(r.Next(), 7, "Receive")
	//positionWheel(rec, x+30, y)
	//rec.active.X = 70
	//rec.background.X = 70
	dep.background.BorderLabelFg = ui.ColorGreen
	rec.background.BorderLabelFg = ui.ColorGreen
	rec.active.ItemFgColor = ui.ColorGreen

	return &pairSelector{dep, rec, dep}
}

func formatSelector(pair *pairSelector, layout *SelectLayout) {
	// TODO: better layout system
	// TODO: set position of pairSelector together i.e. withjust one x,y pair

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

func positionWheel(c *coinWheel, x, y int) {
	c.background.X = x
	c.active.X = x
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
	return bufs
}

// Handle responds to select UI events
func (p *pairSelector) Handle(e string) {
	if e == "/sys/kbd/<up>" {
		p.active.Prev()
	}
	if e == "/sys/kbd/<down>" {
		p.active.Next()
	}
	if e == "/sys/kbd/<right>" {
		p.active.background.BorderFg = ui.ColorWhite
		p.active = p.receive
		p.active.background.BorderFg = ui.ColorRed
	}
	if e == "/sys/kbd/<left>" {
		p.active.background.BorderFg = ui.ColorWhite
		p.active = p.deposit
		p.active.background.BorderFg = ui.ColorRed
	}
}
