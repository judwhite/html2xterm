package html2xterm

import (
	"fmt"
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
		html     string
		expected string
	}{
		{
			html:     loadFile(t, `testdata/001_input_div_span.html`),
			expected: loadFile(t, `testdata/001_output_string.txt`),
		},
		{
			html:     loadFile(t, `testdata/002_input_font_br.html`),
			expected: loadFile(t, `testdata/002_output_string.txt`),
		},
	}

	for _, c := range cases {
		actual, err := Convert(c.html)
		if err != nil {
			t.Fatal(err)
		}

		if c.expected != actual.String() {
			t.Errorf("\nwant:\n---\n%s\n---\ngot:\n---\n%s\n---", c.expected, actual)

			expectedLines := strings.Split(c.expected, "\n")
			actualLines := strings.Split(actual.String(), "\n")

			for i := 0; i < len(expectedLines) && i < len(actualLines); i++ {
				if expectedLines[i] != actualLines[i] {
					t.Errorf("first difference:\n\nwant: '%s'\n      %v\ngot:  '%s'\n      %v\n",
						expectedLines[i], []byte(expectedLines[i]), actualLines[i], []byte(actualLines[i]))
					break
				}
			}
		}

		fmt.Println(actual.ANSI())
	}
}
