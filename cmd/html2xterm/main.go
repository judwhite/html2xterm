package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/judwhite/html2xterm"
)

type options struct {
	js          bool
	centerWidth int
	filenames   []string
}

func main() {
	opts, err := parseArgs(os.Args[1:])
	if err != nil {
		log.Fatalf("%v\nusage: html2xterm [OPTIONS] <filename1> <filename2> ...", err)
	}

	// parse files
	outputs := make([]html2xterm.Output, len(opts.filenames))
	for i := 0; i < len(opts.filenames); i++ {
		filename := opts.filenames[i]
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		output, err := html2xterm.Convert(string(b))
		if err != nil {
			log.Fatal(err)
		}

		outputs[i] = output
	}

	if opts.centerWidth != 0 {
		outputs = html2xterm.Center(opts.centerWidth, outputs...)
	}

	// create output
	var sb strings.Builder
	if opts.js {
		sb.WriteString("const ansi = \"")
	}

	for _, output := range outputs {
		if sb.Len() != 0 {
			sb.WriteString(`\r\n`)
		}

		if opts.js {
			sb.WriteString(output.JS())
		} else {
			fmt.Print(output.ANSI())
		}
	}

	if opts.js {
		sb.WriteString("\";\n\nterm.writeln(ansi);\n")
	}

	fmt.Printf("%s", sb.String())
}

func parseArgs(args []string) (options, error) {
	var flagSet flag.FlagSet
	flagXTermJS := flagSet.Bool("js", false, "output xterm.js code")
	flagWidth := flagSet.Int("width", 0, "center using width")
	err := flagSet.Parse(args)
	if err != nil {
		flagSet.PrintDefaults()
		return options{}, err
	}

	filenames := flagSet.Args()
	if len(filenames) == 0 {
		flagSet.PrintDefaults()
		return options{}, errors.New("no filename(s) specified")
	}

	return options{
		js:          *flagXTermJS,
		centerWidth: *flagWidth,
		filenames:   flagSet.Args(),
	}, nil
}
