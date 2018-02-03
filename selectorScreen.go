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
	p.stats = NewPairStats(pair.deposit.node.coin, pair.receive.node.coin, m)
}

func NewPairSelector(n *windowNode) *pairSelector {
	dep := NewCoinWheel(n, 7, "Deposit")
	rec := NewCoinWheel(n.next, 7, "Receive")
	rec.active.X = 50
	rec.background.X = 50
	rec.active.ItemFgColor = ui.ColorGreen

	return &pairSelector{dep, rec, dep}
}

// TODO: remove dependency on ui???
func (p *pairSelector) Buffers() []ui.Bufferer {
	bufs := p.deposit.Buffers()
	return append(bufs, p.receive.Buffers()...)
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
