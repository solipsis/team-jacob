package main

import (
	"fmt"
	"sort"
	"strings"

	ui "github.com/gizak/termui"
	ss "github.com/solipsis/shapeshift"
)

type Coin struct {
	Name      string
	Symbol    string
	Available bool
}

func toCoin(sc ss.Coin) *Coin {
	return &Coin{
		Name:      sc.Name,
		Symbol:    sc.Symbol,
		Available: sc.Status == "Available",
	}
}

type windowNode struct {
	next, prev *windowNode
	coin       *Coin
}

func wipe() {
	fmt.Printf("\003[0;0H]")
	fmt.Println(strings.Repeat("\n", 100))
}

func newShift() (*ss.NewTransactionResponse, error) {
	return &ss.NewTransactionResponse{
		SendTo:     "0xa6bd216e8e5f463742f37aaab169cabce601835c",
		SendType:   "ETH",
		ReturnTo:   "16FdfRFVPUwiKAceRSqgEfn1tmB4sVUmLh",
		ReturnType: "BTC",
	}, nil
}

func activeCoins() ([]*Coin, error) {
	ssCoins, err := ss.CoinsAsList()
	active := make([]*Coin, 0)
	if err != nil {
		active = append(active, &Coin{Name: "Unable to contact Shapeshift"})
		active = append(active, &Coin{Name: "potato"})
		return active, err
	}
	for _, c := range ssCoins {
		if c.Status == "available" {
			active = append(active, toCoin(c))
		}
	}

	sort.Slice(active, func(i, j int) bool {
		return strings.ToLower(active[i].Name) < strings.ToLower(active[j].Name)
	})
	return active, nil
}

func initWindow(coins []*Coin) *windowNode {
	if len(coins) == 0 {
		return nil
	}
	// TODO: fix edge case of 1 element list
	start := &windowNode{coin: coins[0]}
	prev := start
	for i := 1; i < len(coins); i++ {
		cur := coins[i]
		n := &windowNode{coin: cur, prev: prev}
		prev.next = n
		prev = n
	}
	prev.next = start
	start.prev = prev

	return start
}

func main() {

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	selectScreen := new(pairSelectorScreen)
	selectScreen.Init()

	pair := selectScreen.selector

	p := ui.NewPar(SHAPESHIFT)
	p.Height = 10
	p.Width = 100
	p.TextFgColor = ui.ColorCyan
	p.BorderLabel = "Butt"
	p.BorderFg = ui.ColorWhite

	//pairStats := NewPairStats(pair.deposit.node.coin, pair.receive.node.coin, m)

	qr := NewQR("bloop")
	fmt.Println(qr)

	count := NewCountdown(200)

	draw := func(t int) {
		//wipe()
		//ui.Clear()
		qr.draw()
		ui.Clear()
		ui.Render(pair.Buffers()...)
		count.draw()
		ui.Render(count.gauge)
		//qr.draw()
		//fmt.Printf("\033[10;0H")
		//fmt.Print(buf.String())
		//ui.Render(pairStats.Buffers()...)
		ui.Render(p)
		//count.draw()
		//ui.Render(count.gauge)
	}
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/<up>", func(e ui.Event) {
		pair.Handle(e.Path)
		//pairStats.dep = pair.deposit.node.coin
		//pairStats.rec = pair.receive.node.coin

		draw(0)
	})
	ui.Handle("/sys/kbd/<left>", func(e ui.Event) {
		pair.Handle(e.Path)
		draw(0)
	})
	ui.Handle("/sys/kbd/<right>", func(e ui.Event) {
		pair.Handle(e.Path)
		draw(0)
	})

	ui.Handle("/sys/kbd/<down>", func(e ui.Event) {
		pair.Handle(e.Path)
		//info, ok := m[pair.deposit.node.coin.Symbol+"_"+pair.receive.node.coin.Symbol]
		//if ok {
		//max.Text = fmt.Sprintf("%f", info.Limit)
		//min.Text = fmt.Sprintf("%f", info.Min)
		//rate.Text = fmt.Sprintf("1 %s = %f %s", pair.deposit.node.coin.Symbol, info.Rate, pair.receive.node.coin.Symbol)
		//	fmt.Println(info.Rate)
		//} else {
		//max.Text = "Pair unavailable"
		//}
		draw(0)
	})
	draw(0)

	ui.Loop()
	fmt.Println("done")
}

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
