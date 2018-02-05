package main

import (
	"fmt"
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
	infoY:      20,
	infoHeight: 3,
	infoWidth:  40,
	wheelX:     5,
	wheelWidth: 21,
	wheelY:     15,
}

type PairSelectorScreen struct {
	selector   *pairSelector
	stats      *pairStats
	divider    *ui.Par
	info       *ui.Par
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
		fmt.Println(err)
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
	d, r := pair.deposit.node.coin.Symbol, pair.receive.node.coin.Symbol
	p.stats = NewPairStats(d, r, m[d+"_"+r])

	div := ui.NewPar(" < --- > ")
	div.Border = false
	div.Y = p.layout.wheelY + 5
	div.X = p.layout.wheelX + p.layout.wheelWidth
	div.Height = 3
	div.Width = 8
	p.divider = div
}

func addPairSelector(n *windowNode) *pairSelector {
	dep := NewCoinWheel(n, 7, "Deposit")
	//positionWheel(dep, x, y)
	rec := NewCoinWheel(n.next, 7, "Receive")
	//positionWheel(rec, x+30, y)
	//rec.active.X = 70
	//rec.background.X = 70
	rec.active.ItemFgColor = ui.ColorGreen

	return &pairSelector{dep, rec, dep}
}

func formatSelector(pair *pairSelector, layout *SelectLayout) {
	// TODO: better layout system
	pair.deposit.active.X = layout.wheelX
	pair.deposit.active.Y = layout.wheelY + 4
	pair.deposit.background.X = layout.wheelX
	pair.deposit.background.Y = layout.wheelY

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
	return p.selector.deposit.node.coin, p.selector.receive.node.coin
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
