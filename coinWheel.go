package main

import (
	ui "github.com/gizak/termui"
)

type coinWheel struct {
	node               *windowNode
	background, active *ui.List
	numItems           int
}

type wheelConfig struct {
	x, y, width int
}

// NewCoinWheel creates a new gui element represeting a slot-machine like
// display of a circular list of items with focus on the center item
// numItems must be odd for the wheel to look decent
func NewCoinWheel(n *windowNode, numItems int, label string) *coinWheel {
	// TODO: clean up height calculations
	back := ui.NewList()
	back.Height = numItems + 4 // account for borders and active overlay
	back.Width = 21
	back.X = 50
	back.Y = 20
	back.BorderLabel = label
	back.BorderLabelFg = ui.ColorMagenta

	active := ui.NewList()
	active.Items = []string{n.coin.Name}
	active.Width = 21
	active.Height = 3
	active.X = 50
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

// TODO: decide if i want dupes in the list if less than range size
func (n *windowNode) selection(back, forward int) []*windowNode {
	// walk the starting pointer back and ending pointer forward
	start := n
	for i := 0; i < back; i++ {
		start = start.prev
	}

	// append nodes until we have the size we want
	arr := make([]*windowNode, 0)
	for len(arr) <= forward+back {
		arr = append(arr, start)
		start = start.next
	}
	return arr
}
