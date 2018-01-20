package main

import (
	ui "github.com/gizak/termui"
)

type coinWheel struct {
	node               *windowNode
	background, active *ui.List
	numItems           int
}

// NewCoinWheel creates a new gui element represeting a slot-machine like
// display of a circular list of items with focus on the center item
// numItems must be odd for the wheel to look decent
func NewCoinWheel(n *windowNode, numItems int) *coinWheel {
	// TODO: clean up height calculations
	back := ui.NewList()
	back.Height = numItems + 4 // account for borders and active overlay
	back.Width = 20
	back.X = 30
	back.Y = 20

	active := ui.NewList()
	active.Items = []string{n.coin.Name}
	active.Width = 20
	active.Height = 3
	active.X = 30
	active.Y = 20 + (numItems / 2) + 1
	active.ItemFgColor = ui.ColorRed

	return &coinWheel{active: active, background: back, node: n, numItems: numItems}
}

// We always want to render the background before the active item
func (w *coinWheel) Buffers() []ui.Bufferer {
	w.background.Items = w.backgroundItems()
	w.active.Items = []string{w.node.coin.Name}
	return []ui.Bufferer{w.background, w.active}
}

func (w *coinWheel) Next() {
	w.node = w.node.next
}

func (w *coinWheel) Prev() {
	w.node = w.node.prev
}

func (w *coinWheel) backgroundItems() []string {
	mid := w.numItems / 2
	viewNodes := w.node.selection(mid, mid)
	names := make([]string, 0)

	// Get the coin name for each coin appending a padding item above and below the center item
	for i := 0; i < w.numItems; i++ {
		if i == mid+1 {
			names = append(names, "")
		}
		names = append(names, viewNodes[i].coin.Name)
		if i == mid-1 {
			names = append(names, "")
		}
	}
	return names
}