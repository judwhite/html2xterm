package html2xterm

import "strings"

// Output represents the parsed HTML. See the Convert function.
type Output struct {
	Lines []Line
}

// Center centers all lines using the specified width.
// If the Output's MaxLength is greater than or equal to
// width no change is made.
func (o *Output) Center(width int) {
	padding := (width - o.MaxLength()) / 2
	if padding <= 0 {
		return
	}

	for i := 0; i < len(o.Lines); i++ {
		o.Lines[i] = o.Lines[i].LeftPad(padding)
	}
}

// MaxLength returns the maximum line length in number of characters.
func (o Output) MaxLength() int {
	maxLength := 0
	for _, line := range o.Lines {
		length := line.Length()
		if length > maxLength {
			maxLength = length
		}
	}
	return maxLength
}

// ANSI returns the ANSI representation of the output
// written in true-color notation (24-bit).
func (o Output) ANSI() string {
	var s strings.Builder
	for _, line := range o.Lines {
		s.WriteString(line.ANSI())
		s.WriteRune('\n')
	}
	s.WriteString("\x1b[0m")
	return s.String()
}

// JS returns an unquoted JavaScript string of the ANSI output.
// The string is escaped to work inside either double or single quotes.
func (o Output) JS() string {
	text := strings.ReplaceAll(o.ANSI(), "\n", "\r\n")
	text = strings.ReplaceAll(text, `\`, `\\`)
	text = strings.ReplaceAll(text, `'`, `\x27`)
	text = strings.ReplaceAll(text, `"`, `\x22`)
	text = strings.ReplaceAll(text, "\r", `\r`)
	text = strings.ReplaceAll(text, "\n", `\n`)
	text = strings.ReplaceAll(text, "\x1b", `\x1b`)
	return text
}

// String returns the text-only output without ANSI codes.
func (o Output) String() string {
	var s strings.Builder
	for _, line := range o.Lines {
		s.WriteString(line.String())
		s.WriteRune('\n')
	}
	return s.String()
}

// Line represents a line of Output.
type Line struct {
	Segments []Segment
}

// LeftPad returns a new Line left-padded with n spaces.
func (l Line) LeftPad(n int) Line {
	seg := make([]Segment, len(l.Segments)+1)
	copy(seg[1:], l.Segments)
	seg[0].Text = strings.Repeat(" ", n)
	l.Segments = seg
	return l
}

// Length returns the line length in number of characters.
func (l Line) Length() int {
	length := 0
	for _, segment := range l.Segments {
		length += len(segment.Text)
	}
	return length
}

// ANSI returns the ANSI representation of the line
// written in true-color notation (24-bit).
func (l Line) ANSI() string {
	var s strings.Builder
	for _, segment := range l.Segments {
		s.WriteString(segment.Color.ANSI())
		s.WriteString(segment.Text)
	}
	return s.String()
}

// JS returns an unquoted JavaScript string of the ANSI output.
// The string is escaped to work inside either double or single quotes.
func (l Line) JS() string {
	text := l.ANSI() + "\x1b[0m"
	text = strings.ReplaceAll(text, `\`, `\\`)
	text = strings.ReplaceAll(text, `'`, `\x27`)
	text = strings.ReplaceAll(text, `"`, `\x22`)
	text = strings.ReplaceAll(text, "\x1b", `\x1b`)
	return text
}

// String returns the text-only output without ANSI codes.
func (l Line) String() string {
	var s strings.Builder
	for _, segment := range l.Segments {
		s.WriteString(segment.Text)
	}
	return s.String()
}

// Segment represents a text segment of a line and its color.
type Segment struct {
	Text  string
	Color Color
}
