package main

import (
	"fmt"
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui"
)

var (
	Log *log.Logger
)

type state int

const (
	loading state = iota
	selection
	addressInput
	exchange
)

var activeState = loading

// ui elements
var (
	selectScreen   *PairSelectorScreen
	exchangeScreen *ExchangeScreen
	inputScreen    *InputScreen
	header         *Header
)

func main() {
	// debug logging
	f, err := os.OpenFile("debugLog", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Log.Printf("error opening log file: %v\n", err)
		panic(err)
	}
	defer f.Close()

	Log = log.New(f, "", 0)
	Log.SetOutput(f)

	// start ui thread
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	header = newHeader(DefaultHeaderConfig)
	listenForEvents()

	draw(0)
	activeState = activeState.transitionSelect()
	draw(0)
	ui.Loop()
	fmt.Println("done")
}

// Screen drawing state machine
func draw(t int) {
	Log.Println("Current State: ", activeState)

	//ui.Clear()
	ui.Render(header.draw()...)

	switch activeState {
	case loading:
		load := ui.NewPar("Loading...")
		load.X = DefaultLoadingConfig.X
		load.Y = DefaultLoadingConfig.Y
		load.Width = DefaultLoadingConfig.Width
		load.Height = DefaultLoadingConfig.Height
		load.TextFgColor = ui.ColorYellow
		load.BorderFg = ui.ColorRed

		//ui.Render(header.draw()...)
		ui.Render(load)

	case selection:
		ui.Render(selectScreen.Buffers()...)
		//ui.Render(header.draw()...)

	case addressInput:
		ui.Render(inputScreen.Buffers()...)
		//ui.Render(header.draw()...)

	case exchange:
		// Delays are to ensure QR buffer gets flushed as it
		// is drawn separately from the rest of the ui elements
		//ui.Render(header.draw()...)
		ui.Render(exchangeScreen.Buffers()...)
		time.Sleep(100 * time.Millisecond)
		exchangeScreen.DrawQR()
	}
}

func (s *state) transitionSelect() state {
	selectScreen = NewPairSelectorScreen(DefaultSelectLayout)
	selectScreen.Init()
	ui.Clear()
	return selection
}

func (s *state) transitionInput(prompt string) state {
	inputScreen = NewInputScreen(prompt)
	ui.Clear()
	return addressInput
}

func (s *state) transitionExchange(recAddr string) state {
	Log.Println("Transition Exchange. recAddr: ", recAddr)

	// if destination Address set go to exchange
	// if not prompt
	if recAddr == "" {
		return s.transitionInput("Please enter an address")
	}

	shift, err := newShift(selectScreen.activePair(), recAddr)
	if err != nil {
		Log.Println(err)
		panic(err)
	}

	// if we have just transitioned to this page
	// set up timer to update the time remaining
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			if activeState == exchange {
				ui.Render(exchangeScreen.Buffers()...)
			}
		}
	}()

	exchangeScreen = NewExchangeScreen(shift)
	ui.Clear()
	return exchange
}

type Header struct {
	logo, fox *ui.Par
}

func newHeader(c *HeaderConfig) *Header {
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

	return &Header{logo: logo, fox: fox}
}

func (h *Header) draw() []ui.Bufferer {
	return []ui.Bufferer{h.logo, h.fox}
}

type eventHandler interface {
	Handle() string
}

func listenForEvents() {

	// Subscribe to keyboard event listeners
	ui.Handle("/sys/kbd", func(e ui.Event) {
		Log.Println("ANY KEY", e.Path)
	})
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/<enter>", func(e ui.Event) {
		switch activeState {
		case selection:
			_, rec := selectScreen.SelectedCoins()
			activeState = activeState.transitionExchange(loadDepositAddresses()[rec.Symbol])
		case addressInput:
			activeState = activeState.transitionExchange(inputScreen.input.Text)
		}
		draw(0)
	})
	ui.Handle("/sys/kbd", func(e ui.Event) {
		// TODO; addressInput backspace support
		switch activeState {
		case selection:
			selectScreen.Handle(e.Path)
		case addressInput:
			inputScreen.Handle(e.Path)
		}
		draw(0)
	})
	// Redraw if user resizes gui
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		draw(0)
	})

}
