package main

import (
	"fmt"
	"log"

	ss "github.com/solipsis/shapeshift"
)

type pairSelectorScreen struct {
	selector   *pairSelector
	stats      pairStats
	marketInfo map[string]ss.MarketInfoResponse
}

func (p *pairSelectorScreen) Init() {
	coins, err := activeCoins()
	if err != nil {
		log.Fatal("Unableto contact shapeshift")
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
}
