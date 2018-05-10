package main

import ui "github.com/gizak/termui"

type Header struct {
	logo, fox *ui.Par
}

func newHeader(c *HeaderConfig) *Header {
	logo := ui.NewPar(SHAPESHIFT)
	logo.X = c.LogoX
	logo.Y = c.LogoY
	logo.Width = c.LogoWidth
	logo.Height = c.LogoHeight
	logo.TextFgColor = c.LogoTextFgColor
	logo.Border = false

	fox := ui.NewPar(FOX)
	fox.X = c.FoxX
	fox.Y = c.FoxY
	fox.Width = c.FoxWidth
	fox.Height = c.FoxHeight
	fox.TextFgColor = c.FoxTextFgColor
	fox.Border = false

	return &Header{logo: logo, fox: fox}
}

func (h *Header) draw() []ui.Bufferer {
	return []ui.Bufferer{h.logo, h.fox}
}
