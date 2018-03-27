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

// extract the fields we need from a shapeshift coin response object
func toCoin(sc ss.Coin) *Coin {
	return &Coin{
		Name:      sc.Name,
		Symbol:    sc.Symbol,
		Available: sc.Status == "Available",
	}
}

func wipe() {
	//fmt.Printf("\003[0;0H]")
	fmt.Println(strings.Repeat("\n", 100))
}

// initiate a new shift with Shapeshift
func newShift() (*ss.NewTransactionResponse, error) {

	//s := ss.New{
	//TODO; check other similar method on select screen
	//Pair:      selectScreen.activePair(),
	//ToAddress: "0xa6bd216e8e5f463742f37aaab169cabce601835c",
	//}

	//response, err := s.Shift()
	//if err != nil {
	//panic(err)
	//}

	//if response.ErrorMsg() != "" {
	//panic(response.ErrorMsg())
	//}

	// TODO; setup send and re
	return &ss.NewTransactionResponse{
		SendTo:     "0xa6bd216e8e5f463742f37aaab169cabce601835c",
		SendType:   "ETH",
		ReturnTo:   "16FdfRFVPUwiKAceRSqgEfn1tmB4sVUmLh",
		ReturnType: "BTC",
	}, nil
	/*
		return &ss.NewTransactionResponse{
			SendTo:     "0xa6bd216e8e5f463742f37aaab169cabce601835c",
			SendType:   "ETH",
			ReturnTo:   "16FdfRFVPUwiKAceRSqgEfn1tmB4sVUmLh",
			ReturnType: "BTC",
		}, nil
	*/
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

type state int

const (
	loading state = iota
	selection
	exchange
)

func (s *state) transitionSelect() {
	selectScreen.Init()
	*s = selection
}

type header struct {
	logo, fox *ui.Par
}

func newHeader(c *HeaderConfig) *header {
	logo := ui.NewPar(SHAPESHIFT)
	logo.X = c.LogoX
	logo.Y = c.LogoY
	logo.Width = c.LogoWidth
	logo.Height = c.LogoHeight
	logo.TextFgColor = c.LogoTextFgColor
	logo.Border = false

	fox := ui.NewPar(FOX)
	fox.X = c.FoxX
	fox.Y = c.FoxY
	fox.Width = c.FoxWidth
	fox.Height = c.FoxHeight
	fox.TextFgColor = c.FoxTextFgColor
	fox.Border = false

	return &header{logo: logo, fox: fox}
}

func (h *header) draw() []ui.Bufferer {
	return []ui.Bufferer{h.logo, h.fox}
}

var (
	selectScreen   *PairSelectorScreen
	exchangeScreen *ExchangeScreen
)

func main() {

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	selectScreen = NewPairSelectorScreen(DefaultSelectLayout)

	header := newHeader(DefaultHeaderConfig)
	shift, _ := newShift()
	exchangeScreen = NewExchangeScreen(shift)
	var curState = loading
	first := true

	draw := func(t int) {
		ui.Clear()

		switch curState {
		case loading:
			load := ui.NewPar("Loading...")
			load.X = DefaultLoadingConfig.X
			load.Y = DefaultLoadingConfig.Y
			load.Width = DefaultLoadingConfig.Width
			load.Height = DefaultLoadingConfig.Height
			load.TextFgColor = ui.ColorYellow
			load.BorderFg = ui.ColorRed

			ui.Render(header.draw()...)
			ui.Render(load)
		case selection:
			ui.Render(selectScreen.Buffers()...)
			ui.Render(header.draw()...)
		case exchange:

			// if we have just transitioned to this page
			// set up timer to update the time remaining
			if first {
				first = false
				ticker := time.NewTicker(1 * time.Second)

				go func() {
					for range ticker.C {
						ui.Clear()
						ui.Render(header.draw()...)
						ui.Render(exchangeScreen.Buffers()...)
						time.Sleep(100 * time.Millisecond)
						exchangeScreen.DrawQR()
					}
				}()
			}

			// Delays are to ensure QR buffer gets flushed as it
			// is drawn separately from the rest of the ui elements
			ui.Clear()
			ui.Render(header.draw()...)
			ui.Render(exchangeScreen.Buffers()...)
			time.Sleep(100 * time.Millisecond)
			exchangeScreen.DrawQR()
		}
	}
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {

		//wnd, ok := e.Data.(ui.EvtWnd)
		//if !ok {
		//fmt.Println("HEEEEELLLLPPPP")
		//}
		//type EvtWnd struct {
		//Width  int
		//Height int
		//}
	})
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/<enter>", func(e ui.Event) {
		curState = exchange
		draw(0)

	})

	ui.Handle("/sys/kbd/<up>", func(e ui.Event) {
		selectScreen.Handle(e.Path)
		draw(0)
	})
	ui.Handle("/sys/kbd/<left>", func(e ui.Event) {
		selectScreen.Handle(e.Path)
		draw(0)
	})
	ui.Handle("/sys/kbd/<right>", func(e ui.Event) {
		selectScreen.Handle(e.Path)
		draw(0)
	})

	ui.Handle("/sys/kbd/<down>", func(e ui.Event) {
		selectScreen.Handle(e.Path)
		draw(0)
	})
	draw(0)
	curState.transitionSelect()
	draw(0)
	ui.Loop()
	fmt.Println("done")
}
