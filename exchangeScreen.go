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

type ExchangeScreen struct {
	c       *countdown
	qr      *qr
	stats   *pairStats
	depAddr *ui.Par
	recAddr *ui.Par
	retAddr *ui.Par
	//txProgress *txProgress

}

type countdown struct {
	gauge    *ui.Gauge
	start    time.Time
	duration time.Duration
}

type qr struct {
	//data string
	//io.Reader
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
	//g.PercentColor = ui.ColorRed
	//g.BarColor = ui.ColorGreen
	g.BorderLabel = "Time Remaining"

	return &countdown{gauge: g, start: time.Now(), duration: time.Second * time.Duration(duration)}

}

//TODO: change to enum
func newQR(data string) *qr {
	buf := new(bytes.Buffer)
	/*
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
	*/
	//qrtermival
	//qrterminal.GenerateHalfBlock(data, qrterminal.L, buf)
	qrterminal.Generate(data, qrterminal.L, buf)
	return &qr{buf}
}

// On exchange screen i need to ...
// 1. show qr code
// 2. show min max rate miner
// 3. show awaiting deposit / awaiting exchange / done
// 4. show order id
// 4.5? where to show time remaining
// 5. log transaction
// 6. how to get recieve address from user ?? require entry in file?

func (q *qr) draw() {
	//buf := new(bytes.Buffer)
	//i := rand.Intn(5000000)
	//qrterminal.Generate(strconv.Itoa(i)+"butt"+strconv.Itoa(i), qrterminal.L, buf)
	//qrterminal.Generate("0x05a30f30ad43faea94d1d3d35e3222375bd9dd21", qrterminal.L, buf)
	//q.buf = buf
	fmt.Printf("\033[10;0H")
	for _, l := range strings.Split(q.buf.String(), "\n") {
		fmt.Printf(l)
		fmt.Printf("\033[B")
		fmt.Printf("\r")
	}
	//fmt.Fprint(os.Stdout, q.buf.String())
	//io.Copy(os.Stdout, q)
}
