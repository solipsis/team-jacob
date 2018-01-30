package main

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"

	ui "github.com/gizak/termui"
	"github.com/mdp/qrterminal"
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

type pairSelector struct {
	deposit, receive, active *coinWheel
}

func wipe() {
	fmt.Printf("\003[0;0H]")
	fmt.Println(strings.Repeat("\n", 100))
}

func NewPairSelector(n *windowNode) *pairSelector {
	dep := NewCoinWheel(n, 7)
	rec := NewCoinWheel(n.next, 7)
	rec.active.X = 50
	rec.background.X = 50
	rec.active.ItemFgColor = ui.ColorGreen

	return &pairSelector{dep, rec, dep}
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
		active = append(active, &Coin{Name: "Unable to contart Shapeshift"})
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

type pairSelectorScreen struct {
	selector pairSelector
	stats pairStats
	marketInfo map[string]ss.MarketInfoResponse
}

func (*p pairSelectorScreen) Init() {
	
	
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

	m := make(map[string]ss.MarketInfoResponse)
	for _, v := range rates {
		m[v.Pair] = v
	}

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	n := initWindow(coins)
	pair := NewPairSelector(n)

	p := ui.NewPar(SHAPESHIFT)
	p.Height = 10
	p.Width = 100
	p.TextFgColor = ui.ColorCyan
	p.BorderLabel = "Butt"
	p.BorderFg = ui.ColorWhite

	pairStats := NewPairStats(pair.deposit.node.coin, pair.receive.node.coin, m)
	buf := new(bytes.Buffer)
	config := qrterminal.Config{
		Level:          qrterminal.M,
		Writer:         buf,
		BlackChar:      qrterminal.BLACK,
		WhiteChar:      qrterminal.WHITE,
		WhiteBlackChar: qrterminal.WHITE_BLACK,
		BlackWhiteChar: qrterminal.BLACK_WHITE,
		HalfBlocks:     true,
		//BlackChar:  qrterminal.WHITE,
		//WhiteChar:  qrterminal.BLACK,
	}
	//	fmt.Println(config)
	//qrterminal.GenerateWithConfig("butt", config)

	qrterminal.Generate("blah", qrterminal.L, buf)
	//	fmt.Printf(buf.String())
	s := fmt.Sprintf("%v", config)
	qr := ui.NewPar(s)
	qr.Height = 40
	qr.Width = 100
	qr.TextFgColor = ui.ColorDefault
	qr.TextBgColor = ui.ColorDefault
	p.BorderLabel = "HELP"

	draw := func(t int) {
		//pair.active.background.BorderBg = ui.ColorMagenta
		//ui.Render(p, stats, max, min, rate)
		//ui.Clear()
		wipe()
		ui.Clear()
		ui.Render(p)
		ui.Render(pair.deposit.Buffers()...)
		ui.Render(pair.receive.Buffers()...)
		//ui.Render(qr)
		fmt.Printf("\033[10;0H")
		fmt.Print(buf.String())
		//ui.Render(pairStats.Buffers()...)
	}
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/<up>", func(e ui.Event) {
		pair.active.Prev()
		//recWheel.Next()
		pairStats.dep = pair.deposit.node.coin
		pairStats.rec = pair.receive.node.coin

		draw(0)
	})
	ui.Handle("/sys/kbd/<left>", func(e ui.Event) {
		pair.active.background.BorderFg = ui.ColorWhite
		pair.active = pair.deposit
		pair.active.background.BorderFg = ui.ColorYellow
	})
	ui.Handle("/sys/kbd/<right>", func(e ui.Event) {
		pair.active.background.BorderFg = ui.ColorWhite
		pair.active = pair.receive
		pair.active.background.BorderFg = ui.ColorYellow
	})

	ui.Handle("/sys/kbd/<down>", func(e ui.Event) {
		pair.active.Next()
		pairStats.dep = pair.deposit.node.coin
		pairStats.rec = pair.receive.node.coin
		//recWheel.Prev()
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
