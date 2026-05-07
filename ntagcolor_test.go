package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestResolveStylePreservesLegacyColors(t *testing.T) {
	tests := []struct {
		tag string
		fg  rgbColor
		bg  rgbColor
	}{
		{tag: "green", fg: rgbColor{R: 0, G: 0, B: 0}, bg: rgbColor{R: 0, G: 255, B: 0}},
		{tag: "orange", fg: rgbColor{R: 255, G: 255, B: 255}, bg: rgbColor{R: 255, G: 120, B: 0}},
		{tag: "purple", fg: rgbColor{R: 255, G: 255, B: 255}, bg: rgbColor{R: 100, G: 10, B: 255}},
		{tag: "grey", fg: rgbColor{R: 255, G: 255, B: 255}, bg: rgbColor{R: 100, G: 100, B: 100}},
	}

	for _, tt := range tests {
		fg, bg, ok := resolveStyle(tt.tag)
		if !ok {
			t.Fatalf("resolveStyle(%q) was not found", tt.tag)
		}
		if fg != tt.fg {
			t.Errorf("resolveStyle(%q) fg = %#v, want %#v", tt.tag, fg, tt.fg)
		}
		if bg != tt.bg {
			t.Errorf("resolveStyle(%q) bg = %#v, want %#v", tt.tag, bg, tt.bg)
		}
	}
}

func TestResolveStyleSupportsCSSColorsFromEmacs(t *testing.T) {
	tests := []struct {
		tag string
		bg  rgbColor
		fg  rgbColor
	}{
		{tag: "aliceblue", bg: rgbColor{R: 240, G: 248, B: 255}, fg: rgbColor{R: 0, G: 0, B: 0}},
		{tag: "yellowgreen", bg: rgbColor{R: 154, G: 205, B: 50}, fg: rgbColor{R: 0, G: 0, B: 0}},
		{tag: "rebeccapurple", bg: rgbColor{R: 102, G: 51, B: 153}, fg: rgbColor{R: 255, G: 255, B: 255}},
	}

	for _, tt := range tests {
		fg, bg, ok := resolveStyle(tt.tag)
		if !ok {
			t.Fatalf("resolveStyle(%q) was not found", tt.tag)
		}
		if bg != tt.bg {
			t.Errorf("resolveStyle(%q) bg = %#v, want %#v", tt.tag, bg, tt.bg)
		}
		if fg != tt.fg {
			t.Errorf("resolveStyle(%q) fg = %#v, want %#v", tt.tag, fg, tt.fg)
		}
	}
}

func TestResolveStyleIsCaseInsensitive(t *testing.T) {
	fg, bg, ok := resolveStyle("AliceBlue")
	if !ok {
		t.Fatal("resolveStyle(\"AliceBlue\") was not found")
	}
	if fg != (rgbColor{R: 0, G: 0, B: 0}) {
		t.Errorf("fg = %#v, want black", fg)
	}
	if bg != (rgbColor{R: 240, G: 248, B: 255}) {
		t.Errorf("bg = %#v, want aliceblue", bg)
	}
}

func TestColorSpecSupportsRGBAndHex(t *testing.T) {
	tests := []struct {
		name string
		spec colorSpec
		want rgbColor
	}{
		{name: "rgb", spec: RGB(1, 2, 3), want: rgbColor{R: 1, G: 2, B: 3}},
		{name: "hex", spec: Hex("#0a0b0c"), want: rgbColor{R: 10, G: 11, B: 12}},
	}

	for _, tt := range tests {
		got, err := parseColorSpec(tt.spec)
		if err != nil {
			t.Fatalf("%s parseColorSpec returned error: %v", tt.name, err)
		}
		if got != tt.want {
			t.Errorf("%s parseColorSpec = %#v, want %#v", tt.name, got, tt.want)
		}
	}
}

func TestRenderTagFallbackAndSlashBehavior(t *testing.T) {
	unknown := captureStdout(t, func() {
		renderTag("disruptor")
	})
	wantUnknown := "\x1b[1m\x1b[38;2;255;120;0m\x1b[48;2;255;255;255m.disruptor.\x1B[00m"
	if unknown != wantUnknown {
		t.Fatalf("unknown tag render = %q, want %q", unknown, wantUnknown)
	}

	pathTag := captureStdout(t, func() {
		renderTag("path/tag")
	})
	if pathTag != ".path/tag." {
		t.Fatalf("slash tag render = %q, want %q", pathTag, ".path/tag.")
	}
}

func TestRenderLineEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{name: "no tag", line: "plain.txt", want: "plain.txt"},
		{name: "known tag", line: "file..red..txt", want: "file." + styledTag("red", resolvedTagStyles["red"]) + ".txt"},
		{name: "unknown tag", line: "file..disruptor..txt", want: "file." + styledTag("disruptor", unknownTagStyle) + ".txt"},
		{name: "slash tag", line: "file..path/tag..txt", want: "file..path/tag..txt"},
		{name: "multiple tags", line: "rainbow..red..aliceblue..txt", want: "rainbow." + styledTag("red", resolvedTagStyles["red"]) + styledTag("aliceblue", resolvedTagStyles["aliceblue"]) + ".txt"},
		{name: "leading delimiter", line: "..red..txt", want: "." + styledTag("red", resolvedTagStyles["red"]) + ".txt"},
		{name: "trailing delimiter", line: "file..red..", want: "file." + styledTag("red", resolvedTagStyles["red"]) + "."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderLineString(tt.line); got != tt.want {
				t.Fatalf("renderLineTo(%q) = %q, want %q", tt.line, got, tt.want)
			}
		})
	}
}

func TestRunOutputAndFlags(t *testing.T) {
	input := "file..AliceBlue..txt\n\nlast..path/tag..end"
	want := "file." + styledTag("AliceBlue", resolvedTagStyles["aliceblue"]) + ".txt\nlast..path/tag..end\n"

	for _, args := range [][]string{nil, {"--line-buffered"}, {"--line-buffered=true"}} {
		var out bytes.Buffer
		if err := run(args, strings.NewReader(input), &out); err != nil {
			t.Fatalf("run(%v) returned error: %v", args, err)
		}
		if got := out.String(); got != want {
			t.Fatalf("run(%v) output = %q, want %q", args, got, want)
		}
	}

	var out bytes.Buffer
	if err := run([]string{"unexpected"}, strings.NewReader(input), &out); err == nil {
		t.Fatal("run with unexpected positional argument returned nil error")
	}
}

func TestLineBufferedFlushesEachLine(t *testing.T) {
	var writes recordingWriter
	out := bufio.NewWriterSize(&writes, defaultOutputBufferSize)
	if err := renderInput(strings.NewReader("a..red..b\nc..blue..d\n"), out, true); err != nil {
		t.Fatalf("renderInput line-buffered returned error: %v", err)
	}
	if len(writes.writes) != 2 {
		t.Fatalf("line-buffered writes = %d, want 2", len(writes.writes))
	}

	writes = recordingWriter{}
	out = bufio.NewWriterSize(&writes, defaultOutputBufferSize)
	if err := renderInput(strings.NewReader("a..red..b\nc..blue..d\n"), out, false); err != nil {
		t.Fatalf("renderInput block-buffered returned error: %v", err)
	}
	if len(writes.writes) != 0 {
		t.Fatalf("block-buffered writes before final flush = %d, want 0", len(writes.writes))
	}
	if err := out.Flush(); err != nil {
		t.Fatalf("final flush: %v", err)
	}
	if len(writes.writes) != 1 {
		t.Fatalf("block-buffered writes after final flush = %d, want 1", len(writes.writes))
	}
}

func TestGeneratedStylesAreCurrent(t *testing.T) {
	tmp := t.TempDir()
	for _, name := range []string{"generate_styles.go", "styles_decl.go"} {
		content, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if err := os.WriteFile(filepath.Join(tmp, name), content, 0644); err != nil {
			t.Fatalf("copy %s: %v", name, err)
		}
	}

	cmd := exec.Command("go", "run", "-tags", "generate", "generate_styles.go", "styles_decl.go")
	cmd.Dir = tmp
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("run generator: %v\n%s", err, output)
	}

	got, err := os.ReadFile(filepath.Join(tmp, "styles_gen.go"))
	if err != nil {
		t.Fatalf("read generated temp styles: %v", err)
	}
	want, err := os.ReadFile("styles_gen.go")
	if err != nil {
		t.Fatalf("read checked-in styles_gen.go: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("styles_gen.go is stale; run go generate ./...")
	}
}

func BenchmarkResolveStyle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resolveStyle("rebeccapurple")
	}
}

func BenchmarkRenderLineTagged(b *testing.B) {
	out := discardStringWriter{}
	line := "dir/subdir/file..red..orange..yellowgreen..path/tag..txt"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		renderLineTo(out, line)
	}
}

func BenchmarkRenderInputLargeBlockBuffered(b *testing.B) {
	benchmarkRenderInputLarge(b, false)
}

func BenchmarkRenderInputLargeLineBuffered(b *testing.B) {
	benchmarkRenderInputLarge(b, true)
}

func benchmarkRenderInputLarge(b *testing.B, lineBuffered bool) {
	input := benchmarkInput(20000)
	b.SetBytes(int64(len(input)))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		out := bufio.NewWriterSize(io.Discard, defaultOutputBufferSize)
		if err := renderInput(strings.NewReader(input), out, lineBuffered); err != nil {
			b.Fatal(err)
		}
		if err := out.Flush(); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkInput(lines int) string {
	var sb strings.Builder
	tags := []string{"red", "orange", "yellow", "green", "emerald", "aqua", "teal", "disruptor", "blue", "purple", "gray", "black", "white", "aliceblue", "rebeccapurple", "yellowgreen", "path/tag"}
	for i := 0; i < lines; i++ {
		sb.WriteString("dir/subdir/file_")
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString("..")
		sb.WriteString(tags[i%len(tags)])
		sb.WriteString("..")
		sb.WriteString(tags[(i+5)%len(tags)])
		sb.WriteString("..txt\n")
	}
	return sb.String()
}

type recordingWriter struct {
	writes []string
}

func (w *recordingWriter) Write(p []byte) (int, error) {
	w.writes = append(w.writes, string(p))
	return len(p), nil
}

type discardStringWriter struct{}

func (discardStringWriter) WriteString(s string) (int, error) {
	return len(s), nil
}

func renderLineString(line string) string {
	var out strings.Builder
	renderLineTo(&out, line)
	return out.String()
}

func styledTag(tag string, style resolvedTagStyle) string {
	return style.Prefix + "." + tag + "." + resetANSI
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close stdout pipe writer: %v", err)
	}
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read stdout pipe: %v", err)
	}
	return buf.String()
}
