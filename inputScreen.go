package main

import (
	"strings"

	ui "github.com/gizak/termui"
)

type InputScreen struct {
	prompt string
	input  *ui.Par
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
	return []ui.Bufferer{i.input}
}

func (i *InputScreen) Handle(e string) {

	if i == nil {
		return
	}

	Log.Println(e)
	if strings.HasSuffix(e, "<backspace>") || strings.HasSuffix(e, "<delete>") {
		i.input.Text = i.input.Text[:len(i.input.Text)-1]
		return
	}

	arr := strings.Split(e, "/")
	if len(arr) < 4 || len(arr[3]) > 1 {
		return
	}

	// append the character to the text
	i.input.Text += arr[3]
}

type LoadingScreen struct {
	load *ui.Par
}

func NewLoadingScreen(text string) *LoadingScreen {
	load := ui.NewPar(text)
	load.X = DefaultLoadingConfig.X
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
	err.X = DefaultLoadingConfig.X
	err.Y = DefaultLoadingConfig.Y
	err.Width = len(text) + 4
	err.Height = DefaultLoadingConfig.Height
	err.TextFgColor = ui.ColorYellow
	err.BorderFg = ui.ColorRed
	return &ErrorScreen{err: err}
}

func (e *ErrorScreen) Buffers() []ui.Bufferer {
	return []ui.Bufferer{e.err}
}
