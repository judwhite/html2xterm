package html2xterm

// Center returns a slice of Output's centered together.
// If no outputs are at least minWidth in length, then they are centered using minWidth.
// It is valid to call this function with midWidth set to 0.
func Center(minWidth int, outputs ...Output) []Output {
	maxWidth := minWidth

	for i := 0; i < len(outputs); i++ {
		outputWidth := outputs[i].MaxLength()
		if outputWidth > maxWidth {
			maxWidth = outputWidth
		}
	}

	for i := 0; i < len(outputs); i++ {
		outputs[i].Center(maxWidth)
	}

	return outputs
}
