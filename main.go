package main

import (
	"errors"
	"flag"
	"log"
	"math/big"
	"os"
	"strconv"
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
	addressInput
	amountInput
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
var shift *ss.New

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
	shift = new(ss.New)
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

	case addressInput, amountInput:
		ui.Render(inputScreen.Buffers()...)

	case setup:
		ui.Render(setupScreen.Buffers()...)

	case exchange:
		// Delays are to ensure QR buffer gets flushed as it
		// is drawn separately from the rest of the ui elements
		ui.Render(exchangeScreen.Buffers()...)
		time.Sleep(100 * time.Millisecond)
		exchangeScreen.DrawQR()

		Log.Println("mode", *kkMode)
		Log.Println("device", kkDevice)
		// connect to keepkey
		if *kkMode && kkDevice == nil {
			Log.Println("Connecting to kk")
			devices, err := kk.GetDevices()
			Log.Println("devices", devices, "err", err)
			if err != nil {
				activeState = activeState.transitionError(err)
				return
			}
			kkDevice = devices[0]

			nonce := uint64(20)
			recipient := exchangeScreen.depAddr.Text
			amount := big.NewInt(1337000000000000000)
			gasLimit := big.NewInt(80000)
			gasPrice := big.NewInt(22000000000)
			data := []byte{}
			tx := kk.NewTransaction(nonce, recipient, amount, gasLimit, gasPrice, data)
			tx, err = kkDevice.EthereumSignTx([]uint32{0}, tx)
			// TODO: publish tx using etherscan?
			if err != nil {
				activeState = activeState.transitionError(err)
				return
			}
			ui.StopLoop()
		}
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

func (s *state) transitionAddressInput(prompt string) state {
	inputScreen = NewInputScreen(prompt)
	inputScreen.stats = selectScreen.stats // TODO: cleaner data transfer
	ui.Clear()
	return addressInput
}

func (s *state) transitionAmountInput(prompt string) state {
	inputScreen = NewInputScreen(prompt)
	inputScreen.stats = selectScreen.stats // TODO: cleaner data transfer
	ui.Clear()
	return amountInput
}

func (s *state) transitionSetup(precise bool) state {
	setupScreen = newSetupScreen(precise)
	ui.Clear()
	return setup
}

func (s *state) transitionExchange() state {
	Log.Println("Transition Exchange:", *shift)

	// TODO: Refator the hunterlong ss library
	if shift == nil {
		shift = new(ss.New)
	}

	// if destination Address set go to exchange
	// if not prompt the user for an address
	if shift.ToAddress == "" {
		return s.transitionAddressInput("Please enter an address")
	}

	// If its a quick order make sure the user has specified a deposit amount
	if shift.Amount <= 0 && selectScreen.isPreciseOrder() {
		return s.transitionAmountInput("Please enter a deposit amount")
	}

	nshift, err := newShift(shift, selectScreen.activePair())
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

	exchangeScreen = NewExchangeScreen(nshift)
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
		// TODO: move this logic into their respective screens
		switch activeState {
		case selection:
			_, rec := selectScreen.SelectedCoins()
			selectScreen.jankDrawToggle = true

			shift.ToAddress = loadDepositAddresses()[rec.Symbol]
			activeState = activeState.transitionExchange()
		case addressInput:
			shift.ToAddress = inputScreen.Text()
			activeState = activeState.transitionExchange()
		case amountInput:
			amt, err := strconv.Atoi(inputScreen.Text())
			if err != nil {
				activeState = activeState.transitionError(errors.New("Invalid send amount"))
			}
			shift.Amount = float64(amt)
			activeState = activeState.transitionExchange()
		case encounteredError:
			ui.StopLoop()
		case setup:
			setupScreen.Handle(e.Path)
		}
		draw(0)
	})
	ui.Handle("/sys/kbd", func(e ui.Event) {
		switch activeState {
		case selection:
			selectScreen.Handle(e.Path)
		case addressInput, amountInput:
			inputScreen.Handle(e.Path)
		case setup:
			setupScreen.Handle(e.Path)
		}
		draw(0)
	})
	// Redraw if user resizes gui
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		draw(0)
	})

}
