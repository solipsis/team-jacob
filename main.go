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

type coinWheel struct {
	node               *windowNode
	background, active *ui.List
}

func NewCoinWheel(n *windowNode) *coinWheel {
	back := ui.NewList()
	back.Items = n.windowStart(7).windowStrings(7)
	back.Height = 11
	back.Width = 20
	back.X = 30
	back.Y = 20

	active := ui.NewList()
	active.Items = []string{n.coin.Name}
	active.Width = 20
	active.Height = 3
	active.X = 30
	active.Y = 24
	active.ItemFgColor = ui.ColorRed

	return &coinWheel{active: active, background: back, node: n}
}

// We always want to render the background before the active item
func (w *coinWheel) Buffers() []ui.Bufferer {
	return []ui.Bufferer{w.background, w.active}
}

func (w *coinWheel) Next() {
	w.node = w.node.next
	w.background.Items = w.node.windowStart(7).windowStrings(7)
	w.active.Items = []string{w.node.coin.Name}
}

func (w *coinWheel) Prev() {
	w.node = w.node.prev
	w.background.Items = w.node.windowStart(7).windowStrings(7)
	w.active.Items = []string{w.node.coin.Name}
}

/*
func (n *windowNode) window(size int) []*windowNode {
	target := size / 2
	window := make([]*windowNode, 0)
	// walk list backward till we are half the target size
	for i := target; i > 0; i-- {
		n = n.prev
	}
	// walk forward appending all nodes
	for i := 0; i < size; i++ {
		window = append(window, n)
		n = n.next
	}
	return window
}
*/

func (n *windowNode) windowStart(size int) *windowNode {
	target := size / 2
	start := n
	// walk list backward till we are half the target size
	for i := target; i > 0; i-- {
		start = start.prev
	}

	return start
}

func (n *windowNode) windowStrings(size int) []string {
	ptr := n
	strs := make([]string, 0)
	mid := size / 2 // TODO: fix for non odd
	for i := 0; i < size; i++ {
		if i == mid+1 {
			strs = append(strs, "")
		}
		strs = append(strs, ptr.coin.Name)
		if i == mid-1 {
			strs = append(strs, "")
		}

		ptr = ptr.next
	}
	return strs
}

func activeCoins() ([]ss.Coin, error) {
	coins, err := ss.CoinsAsList()
	if err != nil {
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
	return coins, nil
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

/*
func activeCoins() ([]ss.Coin, error) {
	coins := make([]ss.Coin, 0)
	coinResp, err := ss.CoinsAsList()
	if err != nil {
		return coins, err
	}

	val := reflect.ValueOf(coinResp).Elem()
	for i := 0; i < val.NumField(); i++ {
		if coin, ok := val.Field(i).Field(0).Interface().(ss.Coin); ok && coin.Status == "available" {
			coins = append(coins, val.Field(i).Field(0).Interface().(ss.Coin))
		}
	}
	return coins, nil
}
*/

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

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()
	/*
		list := ui.NewList()
		list.Items = depositCoins
		list.Height = 10
		list.Width = 20
		list.X = 30
		list.Y = 20
	*/

	n := initWindow(coins)
	//start := n.windowStart(5)
	wheel := NewCoinWheel(n)
	recWheel := NewCoinWheel(n)
	recWheel.active.X = 50
	recWheel.background.X = 50
	recWheel.active.ItemFgColor = ui.ColorGreen

	/*
		numItems := 7
		list := ui.NewList()
		list.Items = start.windowStrings(numItems)
		list.Height = 11
		list.Width = 20
		list.X = 30
		list.Y = 20

		lc := ui.NewList()
		lc.Items = []string{n.coin.Name}
		lc.Width = 20
		lc.Height = 3
		lc.X = 30
		lc.Y = 24
		lc.ItemFgColor = ui.ColorRed
	*/
	p := ui.NewPar(SHAPESHIFT)
	p.Height = 10
	p.Width = 100
	p.TextFgColor = ui.ColorCyan
	p.BorderLabel = "Butt"
	p.BorderFg = ui.ColorWhite

	draw := func(t int) {
		//list.Items = n.windowStart(numItems).windowStrings(numItems)
		//lc.Items = []string{n.coin.Name}
		ui.Render(p)
		ui.Render(wheel.Buffers()...)
		ui.Render(recWheel.Buffers()...)
	}
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	//ui.Render(p, list)
	//ui.Handle("/timer/1s", func(e ui.Event) {
	//t := e.Data.(ui.EvtTimer)
	////n = n.next
	//draw(int(t.Count))
	//})
	ui.Handle("/sys/kbd/<up>", func(e ui.Event) {
		wheel.Prev()
		recWheel.Next()
		draw(0)
	})
	ui.Handle("/sys/kbd/<down>", func(e ui.Event) {
		wheel.Next()
		recWheel.Prev()
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
