package main

import (
	"bytes"
	"io"
	"os"
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
