/// 2>/dev/null ; exec gorun "$0" "$@"

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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
	writeColoredRestoring(w, tag, style, "")
}

func writeColoredRestoring(w io.StringWriter, tag string, style resolvedTagStyle, restore string) {
	w.WriteString(style.Prefix)
	w.WriteString(".")
	w.WriteString(tag)
	w.WriteString(".")
	w.WriteString(resetANSI)
	w.WriteString(restore)
}

func renderLine(line string) {
	renderLineTo(os.Stdout, line)
}

func renderLineTo(w io.StringWriter, line string) {
	renderLineToWithRenderer(w, line, func(tag string) {
		renderTagTo(w, tag)
	}, func(text string) {
		w.WriteString(text)
	})
}

func renderLineToPreservingANSI(w io.StringWriter, line string, state *ansiState) {
	renderLineToWithRenderer(w, line, func(tag string) {
		renderTagPreservingANSI(w, tag, state)
	}, func(text string) {
		writeTrackingANSI(w, text, state)
	})
}

func renderLineToWithRenderer(w io.StringWriter, line string, renderTag func(string), writeText func(string)) {
	start := 0
	first := true
	for {
		next := strings.Index(line[start:], "..")
		if next < 0 {
			if first {
				writeText(line[start:])
			} else {
				writeText(".")
				writeText(line[start:])
			}
			return
		}

		next += start
		if first {
			writeText(line[start:next])
			writeText(".")
			first = false
		} else {
			renderTag(line[start:next])
		}
		start = next + 2
	}
}

func renderTagPreservingANSI(w io.StringWriter, tag string, state *ansiState) {
	if style, ok := resolvedTagStyles[strings.ToLower(tag)]; ok {
		writeColoredRestoring(w, tag, style, state.restoreSequence())
		return
	}
	if strings.ContainsRune(tag, '/') {
		writeTrackingANSI(w, ".", state)
		writeTrackingANSI(w, tag, state)
		writeTrackingANSI(w, ".", state)
		return
	}
	writeColoredRestoring(w, tag, unknownTagStyle, state.restoreSequence())
}

func run(args []string, in io.Reader, out io.Writer) error {
	fs := flag.NewFlagSet("ntagcolor", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	lineBuffered := fs.Bool("line-buffered", false, "flush output after each input line")
	preserveANSI := fs.Bool("preserve-ansi", true, "preserve active ANSI SGR styles after rendered tags")
	noPreserveANSI := fs.Bool("no-preserve-ansi", false, "disable ANSI SGR style preservation")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}
	if *noPreserveANSI {
		*preserveANSI = false
	}

	writer := bufio.NewWriterSize(out, defaultOutputBufferSize)
	if err := renderInput(in, writer, *lineBuffered, *preserveANSI); err != nil {
		return err
	}
	return writer.Flush()
}

func renderInput(in io.Reader, out *bufio.Writer, lineBuffered bool, preserveANSI bool) error {
	reader := bufio.NewReader(in)
	var state ansiState
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
		if preserveANSI {
			renderLineToPreservingANSI(out, line, &state)
		} else {
			renderLineTo(out, line)
		}
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

type ansiState struct {
	bold      bool
	dim       bool
	italic    bool
	underline bool
	blink     bool
	reverse   bool
	hidden    bool
	strike    bool
	fg        sgrColor
	bg        sgrColor
}

type sgrColor struct {
	set    bool
	params []string
}

func writeTrackingANSI(w io.StringWriter, text string, state *ansiState) {
	for len(text) > 0 {
		esc := strings.IndexByte(text, '\x1b')
		if esc < 0 {
			w.WriteString(text)
			return
		}
		if esc > 0 {
			w.WriteString(text[:esc])
			text = text[esc:]
		}
		end, ok := ansiCSIEnd(text)
		if !ok {
			w.WriteString(text)
			return
		}
		seq := text[:end]
		w.WriteString(seq)
		if seq[len(seq)-1] == 'm' {
			state.applySGR(seq[2 : len(seq)-1])
		}
		text = text[end:]
	}
}

func ansiCSIEnd(text string) (int, bool) {
	if len(text) < 3 || text[0] != '\x1b' || text[1] != '[' {
		return 1, true
	}
	for i := 2; i < len(text); i++ {
		if text[i] >= 0x40 && text[i] <= 0x7e {
			return i + 1, true
		}
	}
	return 0, false
}

func (s *ansiState) applySGR(raw string) {
	params, ok := parseSGRParams(raw)
	if !ok {
		return
	}
	for i := 0; i < len(params); i++ {
		switch params[i] {
		case 0:
			*s = ansiState{}
		case 1:
			s.bold = true
		case 2:
			s.dim = true
		case 3:
			s.italic = true
		case 4:
			s.underline = true
		case 5, 6:
			s.blink = true
		case 7:
			s.reverse = true
		case 8:
			s.hidden = true
		case 9:
			s.strike = true
		case 21, 22:
			s.bold = false
			s.dim = false
		case 23:
			s.italic = false
		case 24:
			s.underline = false
		case 25:
			s.blink = false
		case 27:
			s.reverse = false
		case 28:
			s.hidden = false
		case 29:
			s.strike = false
		case 30, 31, 32, 33, 34, 35, 36, 37, 90, 91, 92, 93, 94, 95, 96, 97:
			s.fg = sgrColor{set: true, params: []string{strconv.Itoa(params[i])}}
		case 39:
			s.fg = sgrColor{}
		case 40, 41, 42, 43, 44, 45, 46, 47, 100, 101, 102, 103, 104, 105, 106, 107:
			s.bg = sgrColor{set: true, params: []string{strconv.Itoa(params[i])}}
		case 49:
			s.bg = sgrColor{}
		case 38, 48:
			if color, next, ok := extendedColor(params, i); ok {
				if params[i] == 38 {
					s.fg = color
				} else {
					s.bg = color
				}
				i = next
			}
		}
	}
}

func parseSGRParams(raw string) ([]int, bool) {
	if raw == "" {
		return []int{0}, true
	}
	parts := strings.Split(raw, ";")
	params := make([]int, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			params = append(params, 0)
			continue
		}
		value, err := strconv.Atoi(part)
		if err != nil {
			return nil, false
		}
		params = append(params, value)
	}
	return params, true
}

func extendedColor(params []int, i int) (sgrColor, int, bool) {
	if i+2 < len(params) && params[i+1] == 5 {
		return sgrColor{
			set:    true,
			params: []string{strconv.Itoa(params[i]), "5", strconv.Itoa(params[i+2])},
		}, i + 2, true
	}
	if i+4 < len(params) && params[i+1] == 2 {
		return sgrColor{
			set: true,
			params: []string{
				strconv.Itoa(params[i]),
				"2",
				strconv.Itoa(params[i+2]),
				strconv.Itoa(params[i+3]),
				strconv.Itoa(params[i+4]),
			},
		}, i + 4, true
	}
	return sgrColor{}, i, false
}

func (s ansiState) restoreSequence() string {
	params := make([]string, 0, 16)
	if s.bold {
		params = append(params, "1")
	}
	if s.dim {
		params = append(params, "2")
	}
	if s.italic {
		params = append(params, "3")
	}
	if s.underline {
		params = append(params, "4")
	}
	if s.blink {
		params = append(params, "5")
	}
	if s.reverse {
		params = append(params, "7")
	}
	if s.hidden {
		params = append(params, "8")
	}
	if s.strike {
		params = append(params, "9")
	}
	if s.fg.set {
		params = append(params, s.fg.params...)
	}
	if s.bg.set {
		params = append(params, s.bg.params...)
	}
	if len(params) == 0 {
		return ""
	}
	return "\x1b[" + strings.Join(params, ";") + "m"
}

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		log.Fatalln(err.Error())
	}
}
