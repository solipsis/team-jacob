package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	ui "github.com/gizak/termui"
	ss "github.com/solipsis/shapeshift"
)

type windowNode struct {
	next, prev *windowNode
	coin       *ss.Coin
}

type coinStats struct {
	node  *windowNode
	panel *ui.List
}

type pairStats struct {
	dep, rec *windowNode
	panel    *ui.List
	pairInfo map[string]ss.MarketInfoResponse
}

type pairSelector struct {
	deposit, receive, active *coinWheel
}

func NewCoinStats(n *windowNode) *coinStats {
	panel := ui.NewList()
	panel.Height = 6
	panel.Width = 30
	panel.X = 10
	panel.Y = 34

	return &coinStats{node: n, panel: panel}
}

func (p *coinStats) panelCoinStats() []string {
	c := p.node.coin
	items := make([]string, 0)
	items = append(items, c.Name+"("+c.Symbol+")")
	items = append(items, "status: "+c.Status)
	items = append(items, "rate: ")

	return items
}

func (p *coinStats) Buffer() ui.Buffer {
	p.panel.Items = p.panelCoinStats()
	return p.panel.Buffer()
}

func (p *pairStats) statsStrings() []string {
	key := p.dep.coin.Name + "_" + p.rec.coin.Name
	info, ok := p.pairInfo[key]

	stats := make([]string, 0)
	if !ok {
		stats = append(stats, "This pair is not available")
		return stats
	}
	fmt.Println(info)
	return stats
}

// TODO: decide if i want dupes in the list if less than range size
func (n *windowNode) selection(back, forward int) []*windowNode {
	// walk the starting pointer back and ending pointer forward
	start, end := n, n
	for ; back > 0; back-- {
		start = start.prev
	}
	for ; forward > 0; forward-- {
		end = end.next
	}

	// append nodes until after we append the end node
	arr := make([]*windowNode, 0)
	for {
		arr = append(arr, start)
		if start == end {
			break
		}
		start = start.next
	}
	return arr
}

func activeCoins() ([]ss.Coin, error) {
	coins, err := ss.CoinsAsList()
	if err != nil {
		coins = append(coins, ss.Coin{Name: "Unable to contact shapeshift"})
		return coins, err
	}
	active := make([]ss.Coin, 0)
	for _, c := range coins {
		if c.Status == "available" {
			active = append(active, c)
		}
	}
	sort.Slice(coins, func(i, j int) bool {
		return strings.ToLower(coins[i].Name) < strings.ToLower(coins[j].Name)
	})
	return active, nil
}

func initWindow(coins []ss.Coin) *windowNode {
	if len(coins) == 0 {
		return nil
	}
	// TODO: fix edge case of 1 element list
	start := &windowNode{coin: &coins[0]}
	prev := start
	for i := 1; i < len(coins); i++ {
		cur := &coins[i]
		n := &windowNode{coin: cur, prev: prev}
		prev.next = n
		prev = n
	}
	prev.next = start
	start.prev = prev

	return start
}

func main() {

	coins, err := activeCoins()
	if err != nil {
		log.Println("UNable to contact shapeshift")
	}

	depositCoins, recieveCoins := make([]string, 0), make([]string, 0)

	for _, c := range coins {
		depositCoins = append(depositCoins, c.Name)
		recieveCoins = append(recieveCoins, c.Name)
	}

	rates, err := ss.MarketInfo()
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range rates {
		fmt.Println("rate: ", v.Rate)
	}

	m := make(map[string]ss.MarketInfoResponse)
	for _, v := range rates {
		m[v.Pair] = v
	}

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	n := initWindow(coins)
	wheel := NewCoinWheel(n, 7)
	recWheel := NewCoinWheel(n, 7)
	recWheel.active.X = 50
	recWheel.background.X = 50
	recWheel.active.ItemFgColor = ui.ColorGreen
	stats := NewCoinStats(n)
	pair := &pairSelector{wheel, recWheel, wheel}

	p := ui.NewPar(SHAPESHIFT)
	p.Height = 10
	p.Width = 100
	p.TextFgColor = ui.ColorCyan
	p.BorderLabel = "Butt"
	p.BorderFg = ui.ColorWhite

	max := ui.NewPar("1.00345 BTC")
	max.Height = 4
	max.Width = 20
	max.BorderLabel = "Deposit Max"
	max.X = 13
	max.Y = 13

	min := ui.NewPar(".000234 BTC")
	min.Height = 3
	min.Width = 20
	min.BorderLabel = "Min Deposit"
	min.BorderLabelFg = ui.ColorRed
	min.X = 33
	min.Y = 13

	rate := ui.NewPar("1 BTC = 10.2032 ETH")
	rate.Height = 3
	rate.Width = 30
	rate.BorderLabel = "Rate"
	rate.BorderLabelFg = ui.ColorYellow
	rate.X = 53
	rate.Y = 13

	/*
		data := [][]string{
			{"Deposit Min", "Deposit Max", "Miner Fee", "Rate"},
			{"Some Item #1", "AAA", "123", "CCCCC"},
			{"Some Item #2", "BBB", "456", "DDDDD"},
		}
		table := ui.NewTable()
		table.Rows = data // type [][]string
		table.FgColor = ui.ColorWhite
		table.BgColor = ui.ColorDefault
		table.Height = 7
		table.Width = 62
		table.Y = 13
		table.X = 20
		table.Border = true
	*/

	draw := func(t int) {
		//pair.active.background.BorderBg = ui.ColorMagenta
		ui.Render(p, stats, max, min, rate)
		ui.Render(wheel.Buffers()...)
		ui.Render(recWheel.Buffers()...)
	}
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/<up>", func(e ui.Event) {
		pair.active.Prev()
		//recWheel.Next()
		stats.node = wheel.node
		draw(0)
	})
	ui.Handle("/sys/kbd/<left>", func(e ui.Event) {
		pair.active.background.BorderFg = ui.ColorWhite
		pair.active = pair.deposit
		pair.active.background.BorderFg = ui.ColorCyan
	})
	ui.Handle("/sys/kbd/<right>", func(e ui.Event) {
		pair.active.background.BorderFg = ui.ColorWhite
		pair.active = pair.receive
		pair.active.background.BorderFg = ui.ColorCyan
	})

	ui.Handle("/sys/kbd/<down>", func(e ui.Event) {
		pair.active.Next()
		stats.node = wheel.node
		//recWheel.Prev()
		info, ok := m[wheel.node.coin.Symbol+"_"+recWheel.node.coin.Symbol]
		if ok {
			max.Text = fmt.Sprintf("%f", info.Limit)
			min.Text = fmt.Sprintf("%f", info.Min)
			rate.Text = fmt.Sprintf("1 %s = %f %s", wheel.node.coin.Symbol, info.Rate, recWheel.node.coin.Symbol)
			//	fmt.Println(info.Rate)
		} else {
			max.Text = "Pair unavailable"
		}
		draw(0)
	})
	draw(0)

	ui.Loop()
	fmt.Println("done")
}

/*
func layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-50, maxY/2-16, maxX/2+20, maxY/2+16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		config := qrterminal.Config{
			Level:          qrterminal.M,
			Writer:         v,
			BlackChar:      qrterminal.BLACK,
			WhiteChar:      qrterminal.WHITE,
			WhiteBlackChar: qrterminal.WHITE_BLACK,
			BlackWhiteChar: qrterminal.BLACK_WHITE,
			//HalfBlocks:     true,
			//BlackChar:  qrterminal.WHITE,
			//WhiteChar:  qrterminal.BLACK,
		}
		fmt.Println(config)
		//		qrterminal.GenerateHalfBlock("hello potatoes", qrterminal.L, v)
		blue := color.New(color.FgGreen).SprintFunc()
		fmt.Println(blue("https://github.com/mdp/qrterminal"))
		qrterminal.GenerateWithConfig(blue("https://github.com/mdp/qrterminal"), config)
		//qrterminal.Generate("blah", qrterminal.L, v)
		//		fmt.Fprintln(v, "Hello world!")
	}

	if v, err := g.SetView("shapeshift", 1, 1, 100, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		blue := color.New(color.FgGreen).SprintFunc()
		fmt.Fprintf(v, blue(SHAPESHIFT))
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
*/

/*
const SHAPESHIFT = `
   _____ __                    _____ __    _ ______
  / ___// /_  ____ _____  ___ / ___// /_  (_) __/ /_
  \__ \/ __ \/ __ '/ __ \/ _ \\__ \/ __ \/ / /_/ __/
 ___/ / / / / /_/ / /_/ /  __/__/ / / / / / __/ /_
/____/_/ /_/\__,_/ .___/\___/____/_/ /_/_/_/  \__/
                /_/
`
*/
/*
const SHAPESHIFT = `
███████╗██╗  ██╗ █████╗ ██████╗ ███████╗███████╗██╗  ██╗██╗███████╗████████╗
██╔════╝██║  ██║██╔══██╗██╔══██╗██╔════╝██╔════╝██║  ██║██║██╔════╝╚══██╔══╝
███████╗███████║███████║██████╔╝█████╗  ███████╗███████║██║█████╗     ██║
╚════██║██╔══██║██╔══██║██╔═══╝ ██╔══╝  ╚════██║██╔══██║██║██╔══╝     ██║
███████║██║  ██║██║  ██║██║     ███████╗███████║██║  ██║██║██║        ██║
╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚══════╝╚══════╝╚═╝  ╚═╝╚═╝╚═╝        ╚═╝
`
*/

const SHAPESHIFT = "" +
	"  ____  _                      ____  _     _  __ _   \n" +
	" / ___|| |__   __ _ _ __   ___/ ___|| |__ (_)/ _| |_ \n" +
	" \\___ \\| '_ \\ / _` | '_ \\ / _ \\___ \\| '_ \\| | |_| __|\n" +
	"  ___) | | | | (_| | |_) |  __/___) | | | | |  _| |_ \n" +
	" |____/|_| |_|\\__,_| .__/ \\___|____/|_| |_|_|_|  \\__|\n" +
	"                   |_|                               "
