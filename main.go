package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui"
	kk "github.com/solipsis/go-keepkey/pkg/keepkey"
	ss "github.com/solipsis/shapeshift"
)

// ui state
type state int

const (
	loading state = iota
	encounteredError
	selection
	exchange
	setup
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
	setupScreen    *SetupScreen
	header         *Header
)

var kkMode = flag.Bool("kk", false, "keepkey mode")
var kkDevice *kk.Keepkey

func main() {

	flag.Parse()

	// debug logging
	f, err := os.OpenFile("debugLog", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
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

	case setup:
		ui.Render(setupScreen.Buffers()...)

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

func (s *state) transitionSetup(precise bool) state {
	Log.Println("selectStats", selectScreen.stats)
	setupScreen = newSetupScreen(precise, selectScreen.stats)
	ui.Clear()
	return setup
}

func (s *state) transitionExchange() state {

	precise := selectScreen.isPreciseOrder()

	// Parse the amount if this is a precise order
	var amount float64
	if precise {
		amt, err := setupScreen.amount()
		if err != nil {
			return s.transitionError(errors.New("Invalid order amount"))
		}
		amount = amt
	}

	// create the order and submit to ShapeShift
	shift := &ss.New{
		ToAddress:   setupScreen.receiveAddress(),
		FromAddress: setupScreen.returnAddress(),
		Amount:      amount,
		Pair:        selectScreen.activePair(),
	}
	nshift, err := newShift(shift)
	if err != nil {
		return s.transitionError(err)
	}

	// if we have just transitioned to this page
	// set up timer to update the time remaining
	if precise {
		ticker := time.NewTicker(1 * time.Second)
		go func() {
			for range ticker.C {
				if activeState == exchange {
					ui.Render(exchangeScreen.Buffers()...)
				}
			}
		}()
	}

	exchangeScreen = NewExchangeScreen(nshift, precise)
	ui.Clear()
	return exchange
}

func listenForEvents() {

	// Subscribe to keyboard event listeners
	ui.Handle("/sys/kbd", func(e ui.Event) {
		Log.Println("ANY KEY", e.Path)
	})

	ui.Handle("/sys/kbd/<enter>", func(e ui.Event) {
		// TODO: move this logic into their respective screens
		switch activeState {
		case selection:
			selectScreen.jankDrawToggle = true
			activeState = activeState.transitionSetup(selectScreen.isPreciseOrder())
		case encounteredError:
			ui.StopLoop()
		case setup:
			activeState = activeState.transitionExchange()
		}
		draw(0)
	})
	ui.Handle("/sys/kbd", func(e ui.Event) {
		switch activeState {
		case selection:
			selectScreen.Handle(e.Path)
		case setup:
			setupScreen.Handle(e.Path)
		}
		draw(0)
	})
	ui.Handle("/sys/kbd/q", func(e ui.Event) {

		switch activeState {
		case setup:
			setupScreen.Handle(e.Path)
		default:
			ui.StopLoop()
		}
		draw(0)
	})

	// Redraw if user resizes gui
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		draw(0)
	})

}
