package app

import (
	"fmt"
	"io"
	"os"

	"github.com/h2non/filetype"
	"github.com/nsf/termbox-go"
)

type Editor struct {
	seeker   io.ReadSeeker
	view     View
	cursor   uint64
	offset   uint64
	width    uint
	height   uint
	title    string
	subtitle string
	filesize uint64
}

func NewEditor(file *os.File) (*Editor, error) {

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	peek := make([]byte, 262)
	if _, err := file.Read(peek); err != nil {
		return nil, err
	}

	subtitle := "unknown file type"
	if kind, err := filetype.Match(peek); err == nil {
		subtitle = fmt.Sprintf("%s: %s", kind.Extension, kind.MIME.Value)
	}

	return &Editor{
		filesize: uint64(info.Size()),
		seeker:   file,
		title:    file.Name(),
		subtitle: subtitle,
		view:     &SideBySide{},
	}, nil
}

func (editor *Editor) Init() error {
	return termbox.Init()
}

func (editor *Editor) Close() {
	termbox.Close()
}

// Run displays the editor and handles all input. Blocks.
func (editor *Editor) Run() error {

	drawRequired := true

	// set initial size
	{
		w, h := termbox.Size()
		editor.width = uint(uint(w))
		editor.height = uint(uint(h))
	}

	for {
		if drawRequired {
			if err := editor.draw(); err != nil {
				return err
			}
		}
		drawRequired = editor.handleEvent(termbox.PollEvent())
	}

}

func (editor *Editor) viewPosition() (uint, uint) {
	return 1, 2
}

func (editor *Editor) viewArea() (uint, uint) {
	return editor.width - 2, editor.height - 5
}

func (editor *Editor) draw() error {

	editor.calculateOffset()

	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)

	if editor.width < 20 || editor.height < 5 {
		return nil
	}

	if err := editor.drawMenu(); err != nil {
		return err
	}

	posX, posY := editor.viewPosition()
	viewW, viewH := editor.viewArea()

	bw, bh := editor.view.Size(viewW, viewH)
	data := make([]byte, bw*bh)
	if _, err := editor.seeker.Seek(int64(editor.offset), 0); err != nil {
		return err
	}
	bytesRead, err := editor.seeker.Read(data)
	if err != nil && err != io.EOF {
		return err
	}
	// remove empty bytes after actual length of file
	if bytesRead < len(data) {
		data = data[:bytesRead]

	}

	if err := editor.view.Draw(data, editor.cursor, editor.offset, posX, posY, viewW, viewH, defaultColourScheme); err != nil {
		return err
	}
	// @todo: draw overlays
	return termbox.Sync()
}

func (editor *Editor) calculateOffset() {

	viewW, viewH := editor.viewArea()

	// first let's make sure the offset maps to the block size so we don;t get weird behaviour
	w, h := editor.view.Size(viewW, viewH)
	editor.offset = editor.offset - (editor.offset % uint64(w))

	// now let's "scroll up", until we find the cursor
	for editor.cursor < editor.offset {
		editor.offset -= uint64(w)
	}

	// ... or scroll down if we need to
	for editor.cursor-editor.offset >= uint64(w)*uint64(h) {
		editor.offset += uint64(w)
	}
}

func (editor *Editor) drawMenu() error {
	text := fmt.Sprintf(" xu: %s", editor.title)

	for len(text)+1 < int(editor.width)-len(editor.subtitle) {
		text = text + " "
	}

	text = fmt.Sprintf("%s %s", text, editor.subtitle)

	for i := 0; i < int(editor.width); i++ {
		r := ' '
		if i < len(text) {
			r = rune(text[i])
		}
		termbox.SetCell(i, 0, r, defaultColourScheme.MenuFg, defaultColourScheme.MenuBg)
	}

	for i := 0; i < int(editor.width); i++ {
		termbox.SetCell(i, int(editor.height-2), ' ', defaultColourScheme.MenuFg, defaultColourScheme.MenuBg)
	}
	for i := 0; i < int(editor.width); i++ {
		termbox.SetCell(i, int(editor.height-1), ' ', defaultColourScheme.MenuFg, defaultColourScheme.MenuBg)
	}

	return nil
}

// return true if event causes redraw to be required
func (editor *Editor) handleEvent(event termbox.Event) bool {
	switch event.Type {
	case termbox.EventKey:
		return editor.handleKey(event.Key)
	case termbox.EventResize:
		editor.width = uint(event.Width)
		editor.height = uint(event.Height)
		viewW, viewH := editor.viewArea()
		w, _ := editor.view.Size(viewW, viewH)
		editor.offset = editor.offset - (editor.offset % uint64(w))
		return true
	case termbox.EventMouse:
	case termbox.EventError:
	case termbox.EventInterrupt:
	}

	return false
}

func (editor *Editor) handleKey(key termbox.Key) bool {

	viewW, viewH := editor.viewArea()

	switch key {
	case termbox.KeyArrowLeft:
		editor.setCursor(int64(editor.cursor) - 1)
	case termbox.KeyArrowRight:
		editor.setCursor(int64(editor.cursor) + 1)
	case termbox.KeyArrowUp:
		w, _ := editor.view.Size(viewW, viewH)
		editor.setCursor(int64(editor.cursor) - int64(w))
	case termbox.KeyArrowDown:
		w, _ := editor.view.Size(viewW, viewH)
		editor.setCursor(int64(editor.cursor) + int64(w))
	case termbox.KeyPgdn:
		w, h := editor.view.Size(viewW, viewH)
		editor.setCursor(int64(editor.cursor) + int64(w*h))
	case termbox.KeyPgup:
		w, h := editor.view.Size(viewW, viewH)
		editor.setCursor(int64(editor.cursor) - int64(w*h))
	case termbox.KeyHome:
		editor.setCursor(0)
	case termbox.KeyEnd:
		editor.setCursor(int64(editor.filesize))
	case termbox.KeyCtrlC:
		editor.Close()
		os.Exit(0)
	default:
		return false
	}

	return true
}

func (editor *Editor) setCursor(cursor int64) {
	if cursor < 0 {
		editor.cursor = 0
	} else if cursor > int64(editor.filesize) {
		editor.cursor = editor.filesize
	} else { // <= because we allow the cursor to position one place after the bytes to append data
		editor.cursor = uint64(cursor)
	}

}
