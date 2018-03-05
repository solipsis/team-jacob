package main

// represents a circular list of ShapeShift coins
type coinRing interface {
	Next() coinRing
	Prev() coinRing
	Value() *Coin
}

// linked list based implementation of coinRing
type coinNode struct {
	// next, previous
	next, prev *coinNode
	coin       *Coin
}

// ipmlements coinRing.Value()
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
