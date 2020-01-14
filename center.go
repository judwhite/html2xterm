package html2xterm

func Center(width int, outputs ...Output) []Output {
	maxWidth := width

	for i := 0; i < len(outputs); i++ {
		outputWidth := outputs[i].MaxLength()
		if outputWidth > maxWidth {
			maxWidth = outputWidth
		}
	}

	for i := 0; i < len(outputs); i++ {
		outputs[i].AdjustWidth(maxWidth)
	}

	return outputs
}
