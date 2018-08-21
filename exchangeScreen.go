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

func NewExchangeScreen(resp *FixMeshift) *ExchangeScreen {

	cfg := DefaultExchangeConfig

	dep := ui.NewPar(resp.SendTo)
	dep.BorderLabel = "Deposit Address"
	dep.Height = cfg.DepHeight
	dep.Width = cfg.DepWidth
	dep.BorderFg = cfg.DepColor
	dep.X = cfg.DepX
	dep.Y = cfg.DepY

	rec := ui.NewPar(resp.receiveAddr)
	rec.BorderLabel = "Receive Address"
	rec.Height = cfg.RecHeight
	rec.Width = cfg.RecWidth
	rec.BorderFg = cfg.RecColor
	rec.X = cfg.RecX
	rec.Y = cfg.RecY

	ret := ui.NewPar(resp.ReturnTo)
	ret.BorderLabel = "Return Address"
	ret.Height = cfg.RetHeight
	ret.Width = cfg.RetWidth
	ret.BorderFg = cfg.RetColor
	ret.X = cfg.RetX
	ret.Y = cfg.RetY

	c := newCountdown(300)
	qr := newQR(resp.SendTo)
	return &ExchangeScreen{c: c, qr: qr, depAddr: dep, recAddr: rec, retAddr: ret}
}

func (e *ExchangeScreen) Buffers() []ui.Bufferer {
	e.c.update()
	bufs := append([]ui.Bufferer{}, e.depAddr)
	bufs = append(bufs, e.recAddr)
	bufs = append(bufs, e.retAddr)
	return append(bufs, e.c.gauge)
}

// DrawQR muust be called seperately because termui does not accept
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

func (q *qr) draw() {
	// Position terminal cursor with escape codes
	fmt.Printf("\033[10;0H")
	for _, l := range strings.Split(q.buf.String(), "\n") {
		fmt.Printf(l)
		fmt.Printf("\033[B")
		fmt.Printf("\r")
	}
}
