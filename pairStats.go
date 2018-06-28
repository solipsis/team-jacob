package main

import (
	"fmt"

	ui "github.com/gizak/termui"
	ss "github.com/solipsis/shapeshift"
)

type pairStats struct {
	dep, rec            string
	min, max, rate, fee *ui.Par
	info                ss.MarketInfoResponse
}

type statConfig struct {
	x, y                                    int
	minWidth, maxWidth, rateWidth, feeWidth int
}

var defaultConfig = statConfig{
	x:         1,
	y:         8,
	minWidth:  25,
	maxWidth:  25,
	rateWidth: 25,
	feeWidth:  25,
}

// TODO: marketInfoResponse to interface???
// info pane should probably be freed of ss dependencies.
func NewPairStats(dep, rec string, info ss.MarketInfoResponse) *pairStats {
	stats := pairStats{dep: dep, rec: rec, info: info}
	c := defaultConfig
	// TODO: rework layout setup
	stats.min = uiPar("B", "Deposit Min", c.x, c.y, c.minWidth, 3)
	stats.min.BorderFg = ui.ColorBlue
	stats.max = uiPar("A", "Deposit Max", c.x+c.minWidth, c.y, c.maxWidth, 3)
	stats.max.BorderFg = ui.ColorMagenta
	stats.rate = uiPar("C", "Rate", c.x+c.minWidth+c.maxWidth, c.y, c.rateWidth, 3)
	stats.rate.BorderFg = ui.ColorYellow
	stats.fee = uiPar("D", "Miner Fee", c.x+c.minWidth+c.maxWidth+c.rateWidth, c.y, c.feeWidth, 3)
	stats.fee.BorderFg = ui.ColorRed
	//stats.marketInfo = m

	return &stats
}

// TODO: clean up buffers and constructor
func (p *pairStats) Update(dep, rec string, info ss.MarketInfoResponse) {
	p.dep = dep
	p.rec = rec
	p.info = info

}

func (p *pairStats) Buffers() []ui.Bufferer {
	//info := p.marketInfo[p.dep.Symbol+"_"+p.rec.Symbol]
	if p.dep == "" || p.rec == "" || p.info.Rate == 0 {
		p.max.Text = "Pair Unavailable"
		p.min.Text = "Pair Unavailable"
		p.rate.Text = "Pair Unavailable"
		p.fee.Text = "Pair Unavailable"
		return []ui.Bufferer{p.min, p.max, p.rate, p.fee}
	}
	if p.rec == "XRP" || p.dep == "XRP" {
		p.max.Text = "        LOL "
		p.min.Text = " NOPE NOPE NOPE!!!"
		p.rate.Text = "        NOT"
		p.fee.Text = "       Happening"
	} else {
		p.max.Text = fmt.Sprintf("%f %s", p.info.Limit, p.dep)
		p.min.Text = fmt.Sprintf("%f %s", p.info.Min, p.dep)
		p.rate.Text = fmt.Sprintf("1 %s = %f %s", p.dep, p.info.Rate, p.rec)
		p.fee.Text = fmt.Sprintf("%f %s", p.info.MinerFee, p.rec)
	}
	return []ui.Bufferer{p.min, p.max, p.rate, p.fee}
}

func uiPar(text, bLabel string, x, y, width, height int) *ui.Par {
	par := ui.NewPar(text)
	par.BorderLabel = bLabel
	par.X = x
	par.Y = y
	par.Width = width
	par.Height = height
	return par
}
