package main

import (
	"fmt"
	"strings"

	ui "github.com/gizak/termui"
)

type legend struct {
	entries []entry
}

type entry struct {
	key  string
	text string
}

func (l *legend) Buffers() []ui.Bufferer {
	txt := " "
	for _, e := range l.entries {
		txt += fmt.Sprintf("[%s] %s | ", e.key, e.text)
	}
	txt = txt[:len(txt)-2]
	txt = strings.Repeat("-", 80) + "\n" + txt
	p := ui.NewPar(txt)
	p.SetX(1)
	p.SetY(30)
	p.Width = len(txt)
	p.Height = 3
	Log.Println("TEXT", txt)
	p.Border = false
	p.TextFgColor = ui.ColorYellow

	return []ui.Bufferer{p}
}
