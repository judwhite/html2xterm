package html2xterm

import (
	"fmt"
	"strings"
)

type Output struct {
	Lines []Line
}

func (o *Output) AdjustWidth(width int) {
	fmt.Printf("// maxlen = %d\n", o.MaxLength())

	padding := (width - o.MaxLength()) / 2
	fmt.Printf("// padding = %d\n", padding)
	if padding <= 0 {
		return
	}

	for i := 0; i < len(o.Lines); i++ {
		o.Lines[i] = o.Lines[i].Pad(padding)
	}
}

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

func (o Output) JS() string {
	text := strings.ReplaceAll(o.ANSI(), "\n", "\r\n")
	text = strings.ReplaceAll(text, `\`, `\\`)
	text = strings.ReplaceAll(text, `"`, `\"`)
	text = strings.ReplaceAll(text, "\r", `\r`)
	text = strings.ReplaceAll(text, "\n", `\n`)
	text = strings.ReplaceAll(text, "\x1b", `\x1b`)
	return text
}

func (o Output) ANSI() string {
	var s strings.Builder
	for _, line := range o.Lines {
		s.WriteString(line.ANSI())
		s.WriteRune('\n')
	}
	s.WriteString("\x1b[0m")
	return s.String()
}

func (o Output) String() string {
	var s strings.Builder
	for _, line := range o.Lines {
		s.WriteString(line.String())
		s.WriteRune('\n')
	}
	return s.String()
}

type Line struct {
	Segments []Segment
}

func (l Line) Pad(pad int) Line {
	seg := make([]Segment, len(l.Segments)+1)
	copy(seg[1:], l.Segments)
	seg[0].Text = strings.Repeat(" ", pad)
	l.Segments = seg
	return l
}

func (l Line) Length() int {
	length := 0
	for _, segment := range l.Segments {
		length += len(segment.Text)
	}
	return length
}

func (l Line) ANSI() string {
	var s strings.Builder
	for _, segment := range l.Segments {
		s.WriteString(segment.Color.ANSI())
		s.WriteString(segment.Text)
	}
	return s.String()
}

func (l Line) JS() string {
	text := l.ANSI() + "\x1b[0m"
	text = strings.ReplaceAll(text, `\`, `\\`)
	text = strings.ReplaceAll(text, `"`, `\"`)
	text = strings.ReplaceAll(text, "\x1b", `\x1b`)
	return text
}

func (l Line) String() string {
	var s strings.Builder
	for _, segment := range l.Segments {
		s.WriteString(segment.Text)
	}
	return s.String()
}

type Segment struct {
	Text  string
	Color Color
}
