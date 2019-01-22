package app

import termbox "github.com/nsf/termbox-go"

type ColourScheme struct {
	MenuBg    termbox.Attribute
	MenuFg    termbox.Attribute
	AddressBg termbox.Attribute
	AddressFg termbox.Attribute
	HexBg     termbox.Attribute
	HexFg     termbox.Attribute
	PreviewFg termbox.Attribute
	PreviewBg termbox.Attribute
	CursorFg  termbox.Attribute
	CursorBg  termbox.Attribute
}

var defaultColourScheme = ColourScheme{
	MenuBg:    termbox.ColorWhite,
	MenuFg:    termbox.ColorBlack,
	AddressBg: termbox.ColorDefault,
	AddressFg: termbox.ColorBlue,
	HexBg:     termbox.ColorDefault,
	HexFg:     termbox.ColorDefault,
	PreviewBg: termbox.ColorDefault,
	PreviewFg: termbox.ColorDefault,
	CursorBg:  termbox.ColorRed,
	CursorFg:  termbox.ColorWhite,
}
