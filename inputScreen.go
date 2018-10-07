package main

import (
	ui "github.com/gizak/termui"
)

type InputScreen struct {
	prompt string
	input  *ui.Par
	stats  *pairStats
}

func NewInputScreen(prompt string) *InputScreen {
	in := ui.NewPar("")
	in.Border = true
	in.X = 20
	in.Y = 13
	in.Height = 3
	in.Width = 60
	in.BorderLabel = prompt

	return &InputScreen{prompt: prompt, input: in}
}

func (i *InputScreen) Buffers() []ui.Bufferer {
	bufs := []ui.Bufferer{i.input}
	bufs = append(bufs, i.stats.Buffers()...)
	return bufs
}

func (i *InputScreen) Handle(e string) {

	Log.Println(e)

	// All the keys that could be used to "undo" - as of right now github.com/gizak/termui is not
	// working correctly and backspaces are coming through as "C-8>" - when this is fixed/PR'ed
	// these conditions as well as what we allow through anyKey in main, can be updated.
	if e == "<Backspace>" || e == "<Delete>" || e == "C-8>" {
		if len(i.input.Text) > 0 {
			i.input.Text = i.input.Text[:len(i.input.Text)-1]
		}
		return
	}

	if len(e) >= 4 {
		return
	}

	// append the character to the text
	i.input.Text += e
}

func (i *InputScreen) Text() string {
	return i.input.Text
}

type LoadingScreen struct {
	load *ui.Par
}

func NewLoadingScreen(text string) *LoadingScreen {
	load := ui.NewPar(text)
	load.X = DefaultLoadingConfig.X + 10
	load.Y = DefaultLoadingConfig.Y
	load.Width = DefaultLoadingConfig.Width
	load.Height = DefaultLoadingConfig.Height
	load.TextFgColor = ui.ColorYellow
	load.BorderFg = ui.ColorRed
	return &LoadingScreen{load: load}
}

func (l *LoadingScreen) Buffers() []ui.Bufferer {
	return []ui.Bufferer{l.load}
}

type ErrorScreen struct {
	err *ui.Par
}

func NewErrorScreen(text string) *ErrorScreen {
	err := ui.NewPar(text)
	err.X = DefaultLoadingConfig.X - 15
	err.Y = DefaultLoadingConfig.Y
	err.Width = len(text) + 4
	err.Height = DefaultLoadingConfig.Height
	err.TextFgColor = ui.ColorYellow
	err.BorderFg = ui.ColorRed
	err.BorderLabel = "Error"
	return &ErrorScreen{err: err}
}

func (e *ErrorScreen) Buffers() []ui.Bufferer {
	return []ui.Bufferer{e.err}
}
