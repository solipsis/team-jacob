package main

import (
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui"
)

// ui state
type state int

const (
	loading state = iota
	encounteredError
	selection
	addressInput
	exchange
)

var activeState = loading
var Log *log.Logger

// ui elements
var (
	loadingScreen  *LoadingScreen
	errorScreen    *ErrorScreen
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

	// Begin by loading the selection screen
	header = newHeader(DefaultHeaderConfig)
	activeState = activeState.transitionLoading("Loading...")
	draw(0)
	activeState = activeState.transitionSelect()
	draw(0)

	// Loop until ui exits
	listenForEvents()
	ui.Loop()
}

// Screen drawing state machine
func draw(t int) {
	Log.Println("Current State: ", activeState)

	ui.Render(header.draw()...)

	switch activeState {
	case loading:
		ui.Render(loadingScreen.Buffers()...)

	case encounteredError:
		ui.Render(errorScreen.Buffers()...)

	case selection:
		ui.Render(selectScreen.Buffers()...)

	case addressInput:
		ui.Render(inputScreen.Buffers()...)

	case exchange:
		// Delays are to ensure QR buffer gets flushed as it
		// is drawn separately from the rest of the ui elements
		ui.Render(exchangeScreen.Buffers()...)
		time.Sleep(100 * time.Millisecond)
		exchangeScreen.DrawQR()
	}
}

// State transitions
func (s *state) transitionLoading(text string) state {
	loadingScreen = NewLoadingScreen(text)
	ui.Clear()
	return loading
}

func (s *state) transitionError(err error) state {
	errorScreen = NewErrorScreen(err.Error())
	ui.Clear()
	return encounteredError
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
	// if not prompt the user for an address
	if recAddr == "" {
		return s.transitionInput("Please enter an address")
	}

	shift, err := newShift(selectScreen.activePair(), recAddr)
	if err != nil {
		return s.transitionError(err)
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
		case encounteredError:
			ui.StopLoop()
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
