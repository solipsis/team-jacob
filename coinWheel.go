package main

import (
	ui "github.com/gizak/termui"
)

// UI wheel element for select coins
type coinWheel struct {
	ring               coinRing
	background, active *ui.List
	numItems           int
}

type wheelConfig struct {
	x, y, width int
}

// NewCoinWheel creates a new gui element represeting a slot-machine like
// display of a circular list of items with focus on the center item
// numItems must be odd for the wheel to look decent
func NewCoinWheel(r coinRing, numItems int, label string) *coinWheel {
	// TODO: clean up height calculations
	back := ui.NewList()
	back.Height = numItems + 4 // account for borders and active overlay
	back.Width = 21
	back.X = 50
	back.Y = 20
	back.BorderLabel = label
	back.BorderLabelFg = ui.ColorMagenta

	active := ui.NewList()
	active.Items = []string{r.Value().Name}
	active.Width = 21
	active.Height = 3
	active.X = 50
	active.Y = 20 + (numItems / 2) + 1
	active.ItemFgColor = ui.ColorRed

	return &coinWheel{active: active, background: back, ring: r, numItems: numItems}
}

// SelectedCoin returns the coin currently highlighted by the wheel
func (w *coinWheel) SelectedCoin() *Coin {
	return w.ring.Value()
}

// We always want to render the background before the active item
func (w *coinWheel) Buffers() []ui.Bufferer {
	w.background.Items = w.backgroundItems()
	w.active.Items = []string{w.ring.Value().Name}
	return []ui.Bufferer{w.background, w.active}
}

func (w *coinWheel) Next() {
	w.ring = w.ring.Next()
}

func (w *coinWheel) Prev() {
	w.ring = w.ring.Prev()
}

func (w *coinWheel) backgroundItems() []string {
	mid := w.numItems / 2
	viewNodes := coinRange(w.ring, mid, mid)
	names := make([]string, 0)

	// Get the coin name for each coin appending a padding item above and below the center item
	for i := 0; i < w.numItems; i++ {
		if i == mid+1 {
			names = append(names, "")
		}
		names = append(names, viewNodes[i].Value().Name)
		if i == mid-1 {
			names = append(names, "")
		}
	}
	return names
}

// return a a range of coins before, after, and including the current coin
func coinRange(r coinRing, back, forward int) []coinRing {
	// walk the starting pointer back and ending pointer forward
	start := r
	for i := 0; i < back; i++ {
		start = start.Prev()
	}

	// append nodes until we have the size we want
	arr := make([]coinRing, 0)
	for len(arr) <= forward+back {
		arr = append(arr, start)
		start = start.Next()
	}
	return arr
}
