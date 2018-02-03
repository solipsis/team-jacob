package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	ui "github.com/gizak/termui"
	"github.com/mdp/qrterminal"
)

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

func (c *countdown) draw() {
	diff := time.Since(c.start)
	if diff == 0 {
		c.gauge.Percent = 100
	} else {
		c.gauge.Percent = 100 - int((diff*100)/(c.duration))
	}
	seconds := int((c.duration - diff) / time.Second)
	c.gauge.Label = strconv.Itoa(seconds) + "s Remaining"
}

func NewCountdown(duration int) *countdown {
	g := ui.NewGauge()
	g = ui.NewGauge()
	g.Percent = 80
	g.Width = 50
	g.Height = 5
	g.Y = 30

	return &countdown{gauge: g, start: time.Now(), duration: time.Second * time.Duration(duration)}

}

//TODO: change to enum
func NewQR(format string) *qr {
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
	fmt.Printf("\033[10;0H")
	fmt.Println(q.buf.String())
	//io.Copy(os.Stdout, q)
}
