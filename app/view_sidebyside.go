package app

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
)

type SideBySide struct{}

// Size returns the width and height in bytes that can be shown at one time by the view
func (view *SideBySide) Size(cellWidth uint, cellHeight uint) (uint, uint) {
	return (cellWidth - 14) / 4, cellHeight - 1
}

func (view *SideBySide) Draw(data []byte, cursor uint64, offset uint64, x uint, y uint, w uint, h uint, colourScheme ColourScheme) error {

	bw, bh := view.Size(w, h)

	// first lets draw the addresses down the left hand side
	addressSize := uint(12) // 10 with padding both sides
	format := fmt.Sprintf("%%0%dX", addressSize-2)
	for ry := y + 1; ry < bh+y+1; ry++ {
		termbox.SetCell(0, int(ry), ' ', colourScheme.AddressFg, colourScheme.AddressBg)
		termbox.SetCell(int(addressSize-1), int(ry), ' ', colourScheme.AddressFg, colourScheme.AddressBg)
		addr := fmt.Sprintf(format, offset+uint64((ry-1-y)*bw))
		for rx := uint(1); rx < addressSize-1; rx++ {
			termbox.SetCell(int(rx), int(ry), rune(addr[rx-1]), colourScheme.AddressFg, colourScheme.AddressBg)
		}
	}

	// now the address columns across the top
	sub := 0
	for rx := x + addressSize; rx < (bw*3)+x+addressSize; rx += 3 {
		str := fmt.Sprintf("%02X", sub)
		termbox.SetCell(int(rx), int(y), rune(str[0]), colourScheme.AddressFg, colourScheme.AddressBg)
		termbox.SetCell(int(rx+1), int(y), rune(str[1]), colourScheme.AddressFg, colourScheme.AddressBg)
		sub++
	}

	relativeCursor := cursor - offset

	// now lets draw the data as hex
	var pointer int
	var raw byte
	var fg, bg termbox.Attribute
	for ry := y + 1; ry < bh+y+1; ry++ {
		for rx := x + addressSize; rx < (bw*3)+x+addressSize; rx += 3 {

			if pointer < len(data) {
				raw = data[pointer]
				str := fmt.Sprintf("%02X", raw)
				if uint64(pointer) == relativeCursor {
					fg = colourScheme.CursorFg
					bg = colourScheme.CursorBg
				} else {
					fg = colourScheme.HexFg
					bg = colourScheme.HexBg
				}
				termbox.SetCell(int(rx), int(ry), rune(str[0]), fg, bg)
				termbox.SetCell(int(rx+1), int(ry), rune(str[1]), fg, bg)
			} else if pointer == len(data) && uint64(pointer) == relativeCursor {
				termbox.SetCell(int(rx), int(ry), ' ', colourScheme.CursorFg, colourScheme.CursorBg)
				termbox.SetCell(int(rx+1), int(ry), ' ', colourScheme.CursorFg, colourScheme.CursorBg)
			}

			pointer++
		}
	}

	// and finally lets draw the data as literal/raw on thr right hand side
	literalStart := (bw * 3) + x + addressSize + 1
	pointer = 0
	for ry := y + 1; ry < bh+y+1; ry++ {
		for rx := literalStart; rx < literalStart+bw; rx++ {
			if pointer > len(data) {
				break
			} else if pointer == len(data) {
				if uint64(pointer) == relativeCursor {
					termbox.SetCell(int(rx), int(ry), ' ', colourScheme.CursorFg, colourScheme.CursorBg)
				}
			} else {
				r, literal := getRuneForByte(data[pointer])
				if uint64(pointer) == relativeCursor {
					fg = colourScheme.CursorFg
					bg = colourScheme.CursorBg
				} else {
					fg = colourScheme.PreviewFg
					bg = colourScheme.PreviewBg
					if !literal {
						fg = fg | colourScheme.AddressFg
					}
				}
				termbox.SetCell(int(rx), int(ry), r, fg, bg)
			}
			pointer++
		}
	}
	return nil
}
