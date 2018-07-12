package main

import (
	"fmt"

	ui "github.com/gizak/termui"
)

// UI element that displays the min/max/rate/fee between
// a pair of shapeshift coins
type pairStats struct {
	dep, rec            string
	min, max, rate, fee *ui.Par
}

// contruct new pairStats ui element. Use update() for adjusting the displayed data
func newPairStats() *pairStats {
	stats := pairStats{}
	c := defaultStatsConfig

	stats.min = pairInfoBox("Deposit Min", c.x, c.y, c.minWidth, ui.ColorBlue)
	stats.max = pairInfoBox("Deposit Max", c.x+c.minWidth, c.y, c.maxWidth, ui.ColorMagenta)
	stats.rate = pairInfoBox("Rate", c.x+c.minWidth+c.maxWidth, c.y, c.rateWidth, ui.ColorYellow)
	stats.fee = pairInfoBox("Miner Fee", c.x+c.minWidth+c.maxWidth+c.rateWidth, c.y, c.feeWidth, ui.ColorRed)

	return &stats
}

func (p *pairStats) Buffers() []ui.Bufferer {
	return []ui.Bufferer{p.min, p.max, p.rate, p.fee}
}

// construct one of the box subelements of the pairStats element
func pairInfoBox(bLabel string, x, y, width int, borderColor ui.Attribute) *ui.Par {
	par := ui.NewPar("")
	par.BorderLabel = bLabel
	par.X = x
	par.Y = y
	par.Width = width
	par.Height = 3 // 1 + 2 rows of border characters
	par.BorderFg = borderColor
	return par
}

// adjust the text shown in the pairStats boxes to updated values
func (p *pairStats) update(dep, rec string, depMin, depMax, rate, fee float64) {
	if dep == "" || rec == "" || rate == 0 {
		p.max.Text = "Pair Unavailable"
		p.min.Text = "Pair Unavailable"
		p.rate.Text = "Pair Unavailable"
		p.fee.Text = "Pair Unavailable"
	} else if rec == "XRP" {
		p.max.Text = "        LOL "
		p.min.Text = "        NOPE"
		p.rate.Text = "       Wrong"
		p.fee.Text = "        Answer"
	} else {
		p.max.Text = fmt.Sprintf("%f %s", depMax, dep)
		p.min.Text = fmt.Sprintf("%f %s", depMin, dep)
		p.rate.Text = fmt.Sprintf("1 %s = %f %s", dep, rate, rec)
		p.fee.Text = fmt.Sprintf("%f %s", fee, rec)
	}
}
