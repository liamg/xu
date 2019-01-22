package app

// View is a system of rendering the app
type View interface {
	Draw(data []byte, cursor uint64, offset uint64, x uint, y uint, w uint, h uint, colourScheme ColourScheme) error // Draw renders the app to the CLI
	Size(cellWidth uint, cellHeight uint) (uint, uint)                                                               // Size returns the dimensions of the view in renderable bytes
}
