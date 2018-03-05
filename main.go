package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	ss "github.com/solipsis/shapeshift"
)

type Coin struct {
	Name      string
	Symbol    string
	Available bool
}

// represents a circular list of ShapeShift coins
type coinRing interface {
	Next() coinRing
	Prev() coinRing
	Value() *Coin
}

// extract the fields we need from a shapeshift coin response object
func toCoin(sc ss.Coin) *Coin {
	return &Coin{
		Name:      sc.Name,
		Symbol:    sc.Symbol,
		Available: sc.Status == "Available",
	}
}

type coinNode struct {
	// next, previous
	next, prev *coinNode
	coin       *Coin
}

// ipmlements coinlist
func (c *coinNode) Value() *Coin {
	return c.coin
}

// implements coinList
func (c *coinNode) Prev() coinRing {
	return c.prev
}

// implements coinList
func (c *coinNode) Next() coinRing {
	return c.next
}

func wipe() {
	//fmt.Printf("\003[0;0H]")
	fmt.Println(strings.Repeat("\n", 100))
}

// initiate a new shift with Shapeshift
func newShift() (*ss.NewTransactionResponse, error) {
	return &ss.NewTransactionResponse{
		SendTo:     "0xa6bd216e8e5f463742f37aaab169cabce601835c",
		SendType:   "ETH",
		ReturnTo:   "16FdfRFVPUwiKAceRSqgEfn1tmB4sVUmLh",
		ReturnType: "BTC",
	}, nil
}

// activeCoins returns a slice of all the currently active coins on shapeshift
func activeCoins() ([]*Coin, error) {
	ssCoins, err := ss.CoinsAsList()
	active := make([]*Coin, 0)
	if err != nil {
		// Add 2 dummy coins so the scroll wheels still function
		active = append(active, &Coin{Name: "Unable to contact Shapeshift"})
		active = append(active, &Coin{Name: "Unable to contact Shapeshift"})
		return active, err
	}

	// Ignore any coins that aren't available
	for _, c := range ssCoins {
		if c.Status == "available" {
			active = append(active, toCoin(c))
		}
	}

	// Sort alphabetically
	sort.Slice(active, func(i, j int) bool {
		return strings.ToLower(active[i].Name) < strings.ToLower(active[j].Name)
	})
	return active, nil
}

func initWindow(coins []*Coin) *coinNode {
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

type state int

const (
	selection state = iota
	exchange
)

func main() {

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	selectScreen := NewPairSelectorScreen(DefaultSelectLayout)
	selectScreen.Init()

	pair := selectScreen.selector

	p := ui.NewPar(SHAPESHIFT)
	p.Y = 1
	p.Height = 7
	p.Width = 70
	p.X = 5
	p.TextFgColor = ui.ColorCyan
	p.Border = false

	fox := ui.NewPar(FOX)
	fox.Y = 0
	fox.Height = 8
	fox.Width = 29
	fox.TextFgColor = ui.ColorCyan
	fox.X = 70
	fox.Border = false

	exchangeScreen := NewExchangeScreen()
	var curState = selection
	first := true

	draw := func(t int) {
		//time.Sleep(200)
		ui.Clear()

		switch curState {
		case selection:
			ui.Render(selectScreen.Buffers()...)
			ui.Render(p)
			ui.Render(fox)
		case exchange:
			if first {
				first = false
				ticker := time.NewTicker(500 * time.Millisecond)

				go func() {
					for range ticker.C {
						ui.Render(exchangeScreen.Buffers()...)
					}
				}()
			}
			time.Sleep(100 * time.Millisecond)
			ui.Clear()
			ui.Render(p)
			ui.Render(fox)
			ui.Render(exchangeScreen.Buffers()...)
			time.Sleep(100 * time.Millisecond)
			exchangeScreen.DrawQR()
		}
	}
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {

		// If windowsize < min for barcode then no qrcode

		wnd, ok := e.Data.(ui.EvtWnd)
		if !ok {
			fmt.Println("HEEEEELLLLPPPP")
		}
		fmt.Println("Width:", wnd.Width)
		fmt.Println("Height:", wnd.Height)
		//type EvtWnd struct {
		//Width  int
		//Height int
		//}
	})
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/<up>", func(e ui.Event) {
		pair.Handle(e.Path)
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
		draw(0)
	})
	ui.Handle("/sys/kbd/<enter>", func(e ui.Event) {
		//fmt.Println("Exchange")

		/*
			//wipe()
			time.Sleep(100 * time.Millisecond)
			ui.Clear()
			//ui.Render(p)
			ui.Render(p)
			ui.Render(fox)
			ui.Render(exchangeScreen.Buffers()...)
			//ui.Render(fox)
			//ui.Render(exchangeScreen.Buffers()...)
			time.Sleep(100 * time.Millisecond)
			exchangeScreen.DrawQR()
		*/
		curState = exchange
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

const FOX = `            ,^
           ;  ;
\'.,'/      ; ;
/_  _\'-----';

  \/' ,,,,,, ;
    )//     \))`

const SHAPESHIFT = "" +
	"  ____  _                      ____  _     _  __ _   \n" +
	" / ___|| |__   __ _ _ __   ___/ ___|| |__ (_)/ _| |_ \n" +
	" \\___ \\| '_ \\ / _` | '_ \\ / _ \\___ \\| '_ \\| | |_| __|\n" +
	"  ___) | | | | (_| | |_) |  __/___) | | | | |  _| |_ \n" +
	" |____/|_| |_|\\__,_| .__/ \\___|____/|_| |_|_|_|  \\__|\n" +
	"                   |_|                               "
