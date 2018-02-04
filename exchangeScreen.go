package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	ui "github.com/gizak/termui"
	"github.com/mdp/qrterminal"
)

type ExchangeScreen struct {
	c  *countdown
	qr *qr
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

func NewExchangeScreen() *ExchangeScreen {
	c := newCountdown(120)
	qr := newQR("test")
	return &ExchangeScreen{c, qr}

}

func (e *ExchangeScreen) Buffers() []ui.Bufferer {
	e.c.update()
	return append([]ui.Bufferer{}, e.c.gauge)
}

// DrawQR muust be called seperately because termui does not accept
// some of the characters used for the qr code
func (e *ExchangeScreen) DrawQR() {
	e.qr.draw()
}

// update the countdown based on elapsed time
func (c *countdown) update() {
	diff := time.Since(c.start)
	if diff == 0 {
		c.gauge.Percent = 100
	} else {
		c.gauge.Percent = 100 - int((diff*100)/(c.duration))
	}
	seconds := int((c.duration - diff) / time.Second)
	c.gauge.Label = strconv.Itoa(seconds) + "s Remaining"
}

func newCountdown(duration int) *countdown {
	g := ui.NewGauge()
	g.Percent = 80
	g.Width = 50
	g.Height = 5
	g.Y = 40
	g.BorderLabel = "Time Remaining"

	return &countdown{gauge: g, start: time.Now(), duration: time.Second * time.Duration(duration)}

}

//TODO: change to enum
func newQR(format string) *qr {
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
	qrterminal.Generate("blah", qrterminal.L, buf)
	return &qr{buf}
}

func (q *qr) draw() {
	buf := new(bytes.Buffer)
	i := rand.Intn(5000000)
	qrterminal.Generate(strconv.Itoa(i)+"butt"+strconv.Itoa(i), qrterminal.L, buf)
	q.buf = buf
	fmt.Printf("\033[10;0H")
	fmt.Println(q.buf.String())
	//io.Copy(os.Stdout, q)
}
