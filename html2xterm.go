package html2xterm

import (
	"encoding/hex"
	"fmt"
	"html"
	"strings"
)

func Convert(html string) (Output, error) {
	var result Output

	html = strings.ReplaceAll(html, "\n", "")
	html = strings.ReplaceAll(html, "\r", "")
	html = strings.ReplaceAll(html, "</body>", "")
	html = strings.ReplaceAll(html, "</html>", "")
	html = strings.ReplaceAll(html, "<pre>", "")
	html = strings.ReplaceAll(html, "</pre>", "")

	html = strings.ReplaceAll(html, "&nbsp;", " ")

	html = strings.ReplaceAll(html, "<DIV>", "<div>")
	html = strings.ReplaceAll(html, "</DIV>", "</div>")
	html = strings.ReplaceAll(html, "<SPAN", "<span")
	html = strings.ReplaceAll(html, "</SPAN", "</span>")
	html = strings.ReplaceAll(html, "<FONT", "<font")
	html = strings.ReplaceAll(html, "</FONT>", "</font>")
	html = strings.ReplaceAll(html, "<br/>", "<br>")
	html = strings.ReplaceAll(html, "<br />", "<br>")
	html = strings.ReplaceAll(html, "<BR/>", "<br>")
	html = strings.ReplaceAll(html, "<BR />", "<br>")

	bodyStart := strings.Index(html, "<body")
	if bodyStart != -1 {
		end := strings.Index(html[bodyStart:], ">")
		if end == -1 {
			return Output{}, fmt.Errorf("found '<body' but could not find closing '>'")
		}
		html = html[bodyStart+end+1:]
	}

	const beginComment = "<!-- IMAGE BEGINS HERE -->"
	const endComment = "<!-- IMAGE ENDS HERE -->"
	commentStart := strings.Index(html, beginComment)
	if commentStart != -1 {
		html = html[commentStart+len(beginComment):]
		commentEnd := strings.Index(html, endComment)
		if commentEnd != -1 {
			html = html[:commentEnd]
		}
	}

	if strings.HasPrefix(html, "<font size=") {
		start := strings.Index(html, ">")
		html = html[start+1:]
		html = strings.TrimSuffix(html, "</font>")
	}

	for {
		start := strings.Index(html, "<div>")
		useDiv := start != -1
		var end int

		if useDiv {
			// TODO: should probably do something about nested divs
			html = html[start+len("<div>"):]
			end = strings.Index(html, "</div>")
			if end == -1 {
				end = len(html)
			}
		} else {
			start = 0
			end = strings.Index(html, "<br>")
			if end == -1 {
				end = len(html)
			}
		}

		line, err := parseLine(html[:end])
		if err != nil {
			return result, fmt.Errorf("fragment: '%v': %w", []byte(html[:end]), err)
		}

		if len(line.Segments) != 0 || len(result.Lines) > 0 {
			result.Lines = append(result.Lines, line)
		}

		if useDiv {
			end += len("</div>")
			if end > len(html) {
				break
			}
			html = html[end:]
		} else {
			end += len("<br>")
			if end > len(html) {
				break
			}
			html = html[end:]
		}
	}

	// trim trailing blank lines
	for i := len(result.Lines) - 1; i >= 0; i-- {
		line := result.Lines[i]
		if len(line.Segments) != 0 {
			break
		}
		result.Lines = result.Lines[:i]
	}

	return result, nil
}

func parseLine(text string) (Line, error) {
	var line Line
	if len(text) == 0 {
		return line, nil
	}

	var lines []string
	if strings.Contains(text, "<span") {
		lines = strings.SplitAfter(text, "</span>")
	} else if strings.Contains(text, "<font") {
		lines = strings.SplitAfter(text, "</font>")
	} else {
		// TODO: maybe check if strings.TrimSpace(html) returns a non-empty string
		return line, fmt.Errorf("line: '%s', can't find <span> or <font>", text)
	}

	for _, s := range lines {
		s = strings.TrimSpace(s)
		s = strings.TrimSuffix(s, "</font>")
		s = strings.TrimSuffix(s, "</span>")
		if len(s) == 0 {
			continue
		}

		if !strings.HasPrefix(s, "<font") && !strings.HasPrefix(s, "<span") {
			return line, fmt.Errorf("fragment: '%s ...', unhandled html tag", s)
		}

		textStart := strings.Index(s, ">")
		if textStart == -1 {
			return line, fmt.Errorf("fragment: '%s ...', missing '>'", s)
		}

		text := html.UnescapeString(s[textStart+1:])
		if len(text) == 0 {
			continue
		}

		// look for color= or style= ...
		// anything after '>' is text
		var r, g, b uint8
		attrs := s[5:textStart]
		colorIndex := strings.Index(attrs, "color")
		if colorIndex != -1 {
			attrs = attrs[colorIndex:]
			attrs = strings.ReplaceAll(attrs, `'`, `"`)
			attrs = strings.ReplaceAll(attrs, `="`, `:`)
			semiEnd := strings.Index(attrs, `;`)
			quotEnd := strings.Index(attrs, `"`)
			if semiEnd < quotEnd && semiEnd != -1 {
				attrs = attrs[:semiEnd]
			} else if quotEnd != -1 {
				attrs = attrs[:quotEnd]
			} else {
				return line, fmt.Errorf("fragment: '%s ...', missing ';' or '\"'", s)
			}
			attrs = strings.TrimSpace(strings.TrimPrefix(attrs, "color:"))

			if strings.HasPrefix(attrs, "#") {
				attrs = attrs[1:] // trim #
				if len(attrs) == 3 {
					attrs = fmt.Sprintf("%[1]c%[1]c%[2]c%[2]c%[3]c%[3]c", attrs[0], attrs[1], attrs[2])
				}
				if len(attrs) != 6 {
					return line, fmt.Errorf("fragment: '%s ...', unknown color '%s'", s, attrs)
				}
				col, err := hex.DecodeString(attrs)
				if err != nil {
					return line, fmt.Errorf("fragment: '%s ...', error decoding color '%s': %w", s, attrs, err)
				}
				r, g, b = col[0], col[1], col[2]
			} else {
				switch strings.ToLower(attrs) {
				case "black":
					r, g, b = 0, 0, 0
				case "white":
					r, g, b = 255, 255, 255
				default:
					return line, fmt.Errorf("fragment: '%s ...', unknown color '%s'", s, attrs)
				}
			}
		} else {
			return line, fmt.Errorf("fragment: '%s ...', can't find color in '%s'", s, attrs)
		}

		segColor := Color{R: r, G: g, B: b}

		if len(line.Segments) > 0 {
			i := len(line.Segments) - 1
			prev := line.Segments[i]
			if prev.Color == segColor {
				line.Segments[i].Text += text
				continue
			}
			if strings.TrimSpace(prev.Text) == "" && strings.TrimSpace(text) == "" {
				line.Segments[i].Text += text
				continue
			}
		}

		segment := Segment{
			Text:  text,
			Color: segColor,
		}

		line.Segments = append(line.Segments, segment)
	}

	// trim trailing whitespace
	for i := len(line.Segments) - 1; i >= 0; i-- {
		seg := line.Segments[i]
		if strings.TrimSpace(seg.Text) != "" {
			break
		}
		line.Segments = line.Segments[:i]
	}

	return line, nil
}
