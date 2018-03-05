package main

import ui "github.com/gizak/termui"

// Configuration for the header which includes the logo and fox image
type HeaderConfig struct {
	FoxX, FoxY          int
	FoxHeight, FoxWidth int
	FoxTextFgColor      ui.Attribute

	LogoX, LogoY          int
	LogoHeight, LogoWidth int
	LogoTextFgColor       ui.Attribute
}

var DefaultHeaderConfig = &HeaderConfig{
	FoxX:            70,
	FoxY:            0,
	FoxHeight:       8,
	FoxWidth:        29,
	FoxTextFgColor:  ui.ColorCyan,
	LogoX:           5,
	LogoY:           1,
	LogoWidth:       70,
	LogoHeight:      7,
	LogoTextFgColor: ui.ColorCyan,
}

type LoadingConfig struct {
	X, Y          int
	Width, Height int
	TextFgColor   ui.Attribute
}

var DefaultLoadingConfig = &LoadingConfig{
	X:      30,
	Y:      10,
	Width:  10,
	Height: 3,
}

const FOX = `            ,^
           ;  ;
\'.,'/      ; ;
/_  _\'-----';

  \/' ,,,,,, ;
    )//     \))`

const SHAPESHIFT = "" +
	"  ____  _                      ____  _     _  __ _   \n" +
	" / ___|| |__   __ _ _ __   ___/ ___|| |__ (_)/ _| |_ \n" +
	" \\___ \\| '_ \\ / _` | '_ \\ / _ \\___ \\| '_ \\| | |_| __|\n" +
	"  ___) | | | | (_| | |_) |  __/___) | | | | |  _| |_ \n" +
	" |____/|_| |_|\\__,_| .__/ \\___|____/|_| |_|_|_|  \\__|\n" +
	"                   |_|                               "
	/*
	   const SHAPESHIFT = `
	      _____ __                    _____ __    _ ______
	     / ___// /_  ____ _____  ___ / ___// /_  (_) __/ /_
	     \__ \/ __ \/ __ '/ __ \/ _ \\__ \/ __ \/ / /_/ __/
	    ___/ / / / / /_/ / /_/ /  __/__/ / / / / / __/ /_
	   /____/_/ /_/\__,_/ .___/\___/____/_/ /_/_/_/  \__/
	                   /_/
	   `
	*/
	/*
	   const SHAPESHIFT = `
	   ███████╗██╗  ██╗ █████╗ ██████╗ ███████╗███████╗██╗  ██╗██╗███████╗████████╗
	   ██╔════╝██║  ██║██╔══██╗██╔══██╗██╔════╝██╔════╝██║  ██║██║██╔════╝╚══██╔══╝
	   ███████╗███████║███████║██████╔╝█████╗  ███████╗███████║██║█████╗     ██║
	   ╚════██║██╔══██║██╔══██║██╔═══╝ ██╔══╝  ╚════██║██╔══██║██║██╔══╝     ██║
	   ███████║██║  ██║██║  ██║██║     ███████╗███████║██║  ██║██║██║        ██║
	   ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚══════╝╚══════╝╚═╝  ╚═╝╚═╝╚═╝        ╚═╝
	   `
	*/
