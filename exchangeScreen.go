package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	"github.com/mdp/qrterminal"
)

// Screen that displays order information including deposit and receive addresses,
// return address if applicable, and a QR code of the deposit address
type ExchangeScreen struct {
	c       *countdown
	qr      *qr
	stats   *pairStats
	precise bool
	depAddr *ui.Par
	recAddr *ui.Par
	retAddr *ui.Par
}

type countdown struct {
	gauge    *ui.Gauge
	start    time.Time
	duration time.Duration
}

type qr struct {
	buf *bytes.Buffer
}

func NewExchangeScreen(resp *FixMeshift, precise bool) *ExchangeScreen {

	cfg := DefaultExchangeConfig

	qr := newQR(resp.SendTo)
	qrWidth := qr.width()

	dep := ui.NewPar(resp.SendTo)
	dep.BorderLabel = "Deposit Address"
	dep.Height = cfg.DepHeight
	dep.Width = cfg.DepWidth
	dep.BorderFg = cfg.DepColor
	dep.X = qrWidth + 4
	dep.Y = cfg.DepY

	rec := ui.NewPar(resp.receiveAddr)
	rec.BorderLabel = "Receive Address"
	rec.Height = cfg.RecHeight
	rec.Width = cfg.RecWidth
	rec.BorderFg = cfg.RecColor
	rec.X = qrWidth + 4
	rec.Y = cfg.RecY

	ret := ui.NewPar(resp.ReturnTo)
	ret.BorderLabel = "Return Address"
	ret.Height = cfg.RetHeight
	ret.Width = cfg.RetWidth
	ret.BorderFg = cfg.RetColor
	ret.X = qrWidth + 4
	ret.Y = cfg.RetY

	c := newCountdown(300)
	return &ExchangeScreen{c: c, qr: qr, depAddr: dep, recAddr: rec, retAddr: ret, precise: precise}
}

func (e *ExchangeScreen) Buffers() []ui.Bufferer {
	e.c.update()
	bufs := append([]ui.Bufferer{}, e.depAddr)
	bufs = append(bufs, e.recAddr)
	bufs = append(bufs, e.retAddr)
	if e.precise {
		bufs = append(bufs, e.c.gauge)
	}
	return bufs
}

// DrawQR must be called seperately because termui does not accept
// some of the characters used for the qr code
func (e *ExchangeScreen) DrawQR() {
	e.qr.draw()
}

// update the countdown based on elapsed time
func (c *countdown) update() {

	// Adjust filled gauge proportion
	diff := time.Since(c.start)
	if diff == 0 {
		c.gauge.Percent = 100
	} else {
		c.gauge.Percent = 100 - int((diff*100)/(c.duration))
	}

	// calculate time remaining
	seconds := int((c.duration - diff) / time.Second)
	if seconds < 0 {
		seconds = 0
	}
	c.gauge.Label = strconv.Itoa(seconds) + "s Remaining"
}

func newCountdown(duration int) *countdown {
	g := ui.NewGauge()
	g.Percent = 100
	g.Width = 46
	g.Height = 3
	g.Y = 27
	g.X = 67
	g.BorderFg = ui.ColorBlue
	g.BorderLabelFg = ui.ColorYellow
	g.BorderLabel = "Time Remaining"

	return &countdown{gauge: g, start: time.Now(), duration: time.Second * time.Duration(duration)}
}

func newQR(data string) *qr {
	buf := new(bytes.Buffer)

	//qrterminal.GenerateHalfBlock(data, qrterminal.L, buf)
	qrterminal.Generate(data, qrterminal.L, buf)
	return &qr{buf}
}

func (q *qr) width() int {

	// Each line of the qr is a sequence of ansi escape sequences so I can't use any normal
	// string length or utf8 rune counting methods.
	line := strings.Split(q.buf.String(), "\n")[0]
	width := len(strings.Split(line, " "))
	return width
}

func (q *qr) draw() {
	// Position terminal cursor with escape codes
	fmt.Printf("\033[10;0H")
	for _, l := range strings.Split(q.buf.String(), "\n") {
		fmt.Printf(l)
		fmt.Printf("\033[B")
		fmt.Printf("\r")
	}
}
