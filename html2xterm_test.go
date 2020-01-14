package html2xterm

import (
	"io/ioutil"
	"strings"
	"testing"
)

func loadFile(t *testing.T, filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestConvert(t *testing.T) {
	cases := []struct {
		html       string
		outputType string
		expected   string
	}{
		{
			html:       loadFile(t, `testdata/001_input_div_span.html`),
			outputType: "string",
			expected:   loadFile(t, `testdata/001_output_string.txt`),
		},
		{
			html:       loadFile(t, `testdata/002_input_font_br.html`),
			outputType: "string",
			expected:   loadFile(t, `testdata/002_output_string.txt`),
		},
		{
			html:       loadFile(t, `testdata/003_input_html_span_br.html`),
			outputType: "ansi",
			expected:   loadFile(t, `testdata/003_output_ansi.ans`),
		},
		{
			html:       loadFile(t, `testdata/003_input_html_span_br.html`),
			outputType: "xtermjs",
			expected:   loadFile(t, `testdata/003_output_xterm.js`),
		},
	}

	for _, c := range cases {
		output, err := Convert(c.html)
		if err != nil {
			t.Fatal(err)
		}

		var actual string
		switch c.outputType {
		case "string":
			actual = output.String()
		case "ansi":
			actual = output.ANSI()
		case "xtermjs":
			actual = output.JS()
		default:
			t.Errorf("unknown output type '%s'", c.outputType)
			continue
		}

		if c.expected != actual {
			t.Errorf("\nwant:\n---\n%s\n---\ngot:\n---\n%s\n---", c.expected, actual)

			expectedLines := strings.Split(c.expected, "\n")
			actualLines := strings.Split(actual, "\n")

			for i := 0; i < len(expectedLines) && i < len(actualLines); i++ {
				if expectedLines[i] != actualLines[i] {
					if len(expectedLines[i]) > 200 {
						found := false
						for j := 0; j < len(expectedLines[i]) && j < len(actualLines[i]); j++ {
							if expectedLines[i][j] != actualLines[i][j] {
								t.Errorf("first difference at line %d pos %d: want: '%c' %v got: '%c' %v", i, j,
									expectedLines[i][j], expectedLines[i][j], actualLines[i][j], actualLines[i][j])
								found = true
								break
							}
						}
						if found {
							break
						}
					}

					t.Errorf("first difference:\n\nwant: '%s'\n      %v\ngot:  '%s'\n      %v\n",
						expectedLines[i], []byte(expectedLines[i]), actualLines[i], []byte(actualLines[i]))
					break
				}
			}
		}
	}
}
