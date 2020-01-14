package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/judwhite/html2xterm"
)

func main() {
	var flagSet flag.FlagSet
	flagXTermJS := flagSet.Bool("js", false, "output xterm.js code")
	flagWidth := flagSet.Int("width", 0, "center using width")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	filenames := flagSet.Args()

	if len(filenames) == 0 {
		flagSet.PrintDefaults()
		log.Fatal("usage: html2xterm [OPTIONS] <filename1> <filename2> ...")
	}

	centerWidth := *flagWidth
	js := *flagXTermJS

	outputs := make([]html2xterm.Output, len(filenames))

	for i := 0; i < len(filenames); i++ {
		filename := filenames[i]
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

	if centerWidth != 0 {
		outputs = html2xterm.Center(centerWidth, outputs...)
	}

	var sb strings.Builder
	if js {
		sb.WriteString("const ansi = \"")
	}

	for _, output := range outputs {
		if sb.Len() != 0 {
			sb.WriteString(`\r\n`)
		}

		if js {
			sb.WriteString(output.JS())
		} else {
			fmt.Print(output.ANSI())
		}
	}

	if js {
		sb.WriteString("\";\n\nterm.writeln(ansi);\n")
	}

	fmt.Printf("%s", sb.String())
}
