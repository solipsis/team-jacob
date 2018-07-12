package main

import ui "github.com/gizak/termui"

// represents a circular list of ShapeShift coins
type coinRing interface {
	Next() coinRing
	Prev() coinRing
	Value() *Coin
}

type ringSelector struct {
	items              []ringItem
	visibleItems       int
	index              int
	background, active *ui.List
}

func (r *ringSelector) Next() {
	r.index = (r.index + 1) % len(r.items)
}

func (r *ringSelector) Prev() {
	r.index = (r.index - 1 + len(r.items)) % len(r.items)
}

func (r *ringSelector) Selected() ringItem {
	return r.items[r.index]
}

func (r *ringSelector) Buffers() []ui.Bufferer {
	r.background.Items = r.backgroundItems()
	r.active.Items = []string{r.background.Items[len(r.background.Items)/2]}
	return []ui.Bufferer{r.background, r.active}
}

func (r *ringSelector) backgroundItems() []string {
	mid := r.visibleItems / 2
	start := (r.index - r.visibleItems/2 + len(r.items)) % len(r.items)
	names := make([]string, 0)

	// append all names adding a padding item above and below the center item
	for i := 0; i < r.visibleItems; i++ {
		if i == mid+1 {
			names = append(names, "")
		}
		names = append(names, r.items[(start+i)%len(r.items)].Text())
		if i == mid-1 {
			names = append(names, "")
		}
	}
	return names
}

func NewRingSelector(items []ringItem, label string, x, y, visibleItems int) *ringSelector {

	// TODO: clean up height calculations
	back := ui.NewList()
	back.Height = visibleItems + 4 // account for borders and active overlay
	back.Width = 21
	back.X = x
	back.Y = y
	back.BorderLabel = label
	back.BorderLabelFg = ui.ColorMagenta

	active := ui.NewList()
	active.Items = []string{items[0].Text()}
	active.Width = 21
	active.Height = 3
	active.X = x
	active.Y = y + (visibleItems / 2) + 1
	active.ItemFgColor = ui.ColorRed

	return &ringSelector{
		items:        items, // Should this be read-only copy
		index:        0,
		background:   back,
		active:       active,
		visibleItems: visibleItems,
	}

}

// linked list based implementation of coinRing
type coinNode struct {
	// next, previous
	next, prev *coinNode
	coin       *Coin
}

type ringItem interface {
	Text() string
}

/*
type item struct {
	ringItem
	next, prev ringItem
}
*/

// implements coinRing.Value()
func (c *coinNode) Value() *Coin {
	return c.coin
}

// implements coinRing.Prev()
func (c *coinNode) Prev() coinRing {
	return c.prev
}

// implements coinList.Next()
func (c *coinNode) Next() coinRing {
	return c.next
}

// Take a list of coins and convert them to a circular coinRing
func toCoinRing(coins []*Coin) *coinNode {
	if len(coins) == 0 {
		return nil
	}
	// TODO: fix edge case of 1 element list
	start := &coinNode{coin: coins[0]}
	prev := start
	for i := 1; i < len(coins); i++ {
		cur := coins[i]
		n := &coinNode{coin: cur, prev: prev}
		prev.next = n
		prev = n
	}
	prev.next = start
	start.prev = prev

	return start
}
