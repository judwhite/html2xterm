package html2xterm

import "fmt"

// Color represents a 24-bit RGB color.
type Color struct {
	R, G, B uint8
}

// ANSI returns the Color's ANSI code using true-color notation (24-bit).
func (c Color) ANSI() string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
}
