package html2xterm

import "fmt"

type Color struct {
	R, G, B uint8
}

func (c Color) ANSI() string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
}
