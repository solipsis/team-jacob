package main

import (
	"fmt"
	"log"

	ui "github.com/gizak/termui"
	ss "github.com/solipsis/shapeshift"
)

type pairSelectorScreen struct {
	selector   *pairSelector
	stats      *pairStats
	marketInfo map[string]ss.MarketInfoResponse
}

type pairSelector struct {
	deposit, receive, active *coinWheel
}

func (p *pairSelectorScreen) Init() {
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
	pair := NewPairSelector(n)

	p.selector = pair
	p.marketInfo = m
	d, r := pair.deposit.node.coin.Symbol, pair.receive.node.coin.Symbol
	p.stats = NewPairStats(d, r, m[d+"_"+r])
}

func NewPairSelector(n *windowNode) *pairSelector {
	dep := NewCoinWheel(n, 7, "Deposit")
	rec := NewCoinWheel(n.next, 7, "Receive")
	rec.active.X = 70
	rec.background.X = 70
	rec.active.ItemFgColor = ui.ColorGreen

	return &pairSelector{dep, rec, dep}
}

func (p *pairSelectorScreen) SelectedCoins() (dep, rec *Coin) {
	return p.selector.deposit.node.coin, p.selector.receive.node.coin
}

// TODO: remove dependency on ui???
func (p *pairSelector) Buffers() []ui.Bufferer {
	bufs := p.deposit.Buffers()
	bufs = append(bufs, p.receive.Buffers()...)
	return bufs
}

func (p *pairSelectorScreen) Buffers() []ui.Bufferer {
	bufs := p.selector.Buffers()
	// TODO: refactor this
	d, r := p.SelectedCoins()
	p.stats.Update(d.Symbol, r.Symbol, p.marketInfo[d.Symbol+"_"+r.Symbol])
	bufs = append(bufs, p.stats.Buffers()...)
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
