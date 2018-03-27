package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	"github.com/mdp/qrterminal"
	ss "github.com/solipsis/shapeshift"
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

func NewExchangeScreen(resp *ss.NewTransactionResponse) *ExchangeScreen {
	c := newCountdown(330)
	qr := newQR(resp.SendTo)

	//dep := ui.NewPar("0x05a30f30ad43faea94d1d3d35e3222375bd9dd21")
	dep := ui.NewPar(resp.SendTo)
	dep.Height = 3
	dep.Width = 46
	dep.BorderLabel = "Deposit Address"
	dep.BorderFg = ui.ColorRed
	dep.X = 67
	dep.Y = 15

	rec := ui.NewPar("1F1tAaz5x1HUXrCNLbtMDqcw6o5GNn4xqX")
	//rec := ui.NewPar(resp.Withdrawal)
	rec.Height = 3
	rec.Width = 46
	rec.BorderLabel = "Receive Address"
	rec.BorderFg = ui.ColorGreen
	rec.X = 67
	rec.Y = 19

	//ret := ui.NewPar("0x6b67c94fc31510707F9c0f1281AaD5ec9a2EEFF0")
	ret := ui.NewPar(resp.ReturnTo)
	ret.Height = 3
	ret.Width = 46
	ret.BorderLabel = "Return Address"
	ret.BorderFg = ui.ColorYellow
	ret.X = 67
	ret.Y = 23

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
	g.Percent = 80
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
