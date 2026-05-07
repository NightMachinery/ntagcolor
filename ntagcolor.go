/// 2>/dev/null ; exec gorun "$0" "$@"

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	defaultOutputBufferSize = 64 * 1024
	resetANSI               = "\x1B[00m"
)

type resolvedTagStyle struct {
	FG     rgbColor
	BG     rgbColor
	Prefix string
}

func resolveStyle(tag string) (rgbColor, rgbColor, bool) {
	style, ok := resolvedTagStyles[strings.ToLower(tag)]
	if !ok {
		return rgbColor{}, rgbColor{}, false
	}
	return style.FG, style.BG, true
}

func printColored(text string, fg, bg rgbColor) {
	fmt.Print(ansiPrefix(fg, bg), text, resetANSI)
}

func ansiPrefix(fg, bg rgbColor) string {
	return fmt.Sprintf("\x1b[1m\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm", fg.R, fg.G, fg.B, bg.R, bg.G, bg.B)
}

func renderTag(tag string) {
	renderTagTo(os.Stdout, tag)
}

func renderTagTo(w io.StringWriter, tag string) {
	if style, ok := resolvedTagStyles[strings.ToLower(tag)]; ok {
		writeColored(w, tag, style)
		return
	}
	if strings.ContainsRune(tag, '/') {
		w.WriteString(".")
		w.WriteString(tag)
		w.WriteString(".")
		return
	}
	writeColored(w, tag, unknownTagStyle)
}

func writeColored(w io.StringWriter, tag string, style resolvedTagStyle) {
	w.WriteString(style.Prefix)
	w.WriteString(".")
	w.WriteString(tag)
	w.WriteString(".")
	w.WriteString(resetANSI)
}

func renderLine(line string) {
	renderLineTo(os.Stdout, line)
}

func renderLineTo(w io.StringWriter, line string) {
	start := 0
	first := true
	for {
		next := strings.Index(line[start:], "..")
		if next < 0 {
			if first {
				w.WriteString(line[start:])
			} else {
				w.WriteString(".")
				w.WriteString(line[start:])
			}
			return
		}

		next += start
		if first {
			w.WriteString(line[start:next])
			w.WriteString(".")
			first = false
		} else {
			renderTagTo(w, line[start:next])
		}
		start = next + 2
	}
}

func run(args []string, in io.Reader, out io.Writer) error {
	fs := flag.NewFlagSet("ntagcolor", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	lineBuffered := fs.Bool("line-buffered", false, "flush output after each input line")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	writer := bufio.NewWriterSize(out, defaultOutputBufferSize)
	if err := renderInput(in, writer, *lineBuffered); err != nil {
		return err
	}
	return writer.Flush()
}

func renderInput(in io.Reader, out *bufio.Writer, lineBuffered bool) error {
	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		if line == "" {
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			continue
		}
		renderLineTo(out, line)
		out.WriteString("\n")
		if lineBuffered {
			if err := out.Flush(); err != nil {
				return err
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		log.Fatalln(err.Error())
	}
}
