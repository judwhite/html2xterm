# html2xterm
[![GoDoc](https://godoc.org/github.com/judwhite/html2xterm?status.svg)](https://godoc.org/github.com/judwhite/html2xterm) [![MIT License](http://img.shields.io/:license-mit-blue.svg)](https://github.com/judwhite/html2xterm/blob/develop/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/judwhite/html2xterm)](https://goreportcard.com/report/github.com/judwhite/html2xterm)
[![CircleCI](https://circleci.com/gh/judwhite/html2xterm.svg?style=shield)](https://circleci.com/gh/judwhite/html2xterm) [![codecov](https://codecov.io/gh/judwhite/html2xterm/branch/develop/graph/badge.svg)](https://codecov.io/gh/judwhite/html2xterm)

Convert (certain) colorized HTML to ANSI and xterm.js

Works well with HTML generated by:
- http://patorjk.com/text-color-fader/
- https://asciiart.club/
- https://www.text-image.com/convert/

Usage:

```
$ ./html2xterm saved.html > saved.ans
$ ./html2xterm -js saved.html > saved.js
```
