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

// TODO: marketInfoResponse to interface???
// info pane should probably be freed of ss dependencies.
func NewPairStats(dep, rec string, info ss.MarketInfoResponse) *pairStats {
	stats := pairStats{dep: dep, rec: rec, info: info}
	stats.min = uiPar("B", "Deposit Min", 13, 13, 20, 3)
	stats.min.BorderFg = ui.ColorBlue
	stats.max = uiPar("A", "Deposit Max", 33, 13, 20, 3)
	stats.max.BorderFg = ui.ColorMagenta
	stats.rate = uiPar("C", "Rate", 53, 13, 25, 3)
	stats.rate.BorderFg = ui.ColorYellow
	stats.fee = uiPar("D", "Miner Fee", 78, 13, 20, 3)
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
	p.max.Text = fmt.Sprintf("%f %s", p.info.Limit, p.dep)
	p.min.Text = fmt.Sprintf("%f %s", p.info.Min, p.dep)
	p.rate.Text = fmt.Sprintf("1 %s = %f %s", p.dep, p.info.Rate, p.rec)
	p.fee.Text = fmt.Sprintf("%f %s", p.info.MinerFee, p.rec)
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

/*
func (p *pairStats) statsStrings() []string {
	key := p.dep + "_" + p.rec
	info, ok := p.marketInfo[key]

	stats := make([]string, 0)
	if !ok {
		stats = append(stats, "This pair is not available")
		return stats
	}
	fmt.Println(info)
	return stats
}
*/
