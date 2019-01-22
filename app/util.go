package app

func getRuneForByte(b byte) (rune, bool) { // true if literal, false for symbols representing runes e.g. new line
	if b < 32 || b >= 127 {
		return '.', false
	}

	return rune(b), true
}
