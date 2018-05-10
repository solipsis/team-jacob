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

	arr := strings.Split(e, "/")
	if len(arr) < 4 || len(arr[3]) > 1 {
		return
	}

	// append the character to the text
	i.input.Text += arr[3]
}
