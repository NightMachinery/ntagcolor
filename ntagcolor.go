/// 2>/dev/null ; exec gorun "$0" "$@"

package main

import (
	. "fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func PrintColored(text string, fgR, fgG, fgB, bgR, bgG, bgB int) {
	Printf("\x1b[1m\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%s\x1B[00m", fgR, fgG, fgB, bgR, bgG, bgB, text)
}
func main() {
	inBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err.Error())
	}
	input := string(inBytes)
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		//if strings.TrimSpace(line) == "" {
		if (line) == "" {
			continue
		}
		tokens := strings.Split(line, "..")
		for i, t := range tokens {
			if i == 0 {
				if i == (len(tokens) - 1) {
					Printf("%s", t)
				} else {
					Printf("%s.", t)
				}
			} else if i == (len(tokens) - 1) {
				Printf(".%s", t)
			} else {
				if t == "red" {
					PrintColored(Sprintf(".%s.", t),255, 255, 255, 255, 0, 0)
				} else if t == "blue" {
					PrintColored(Sprintf(".%s.", t),255, 255, 255, 0, 0, 255)
				} else if t == "green" {
					PrintColored(Sprintf(".%s.", t),0, 0, 0, 0, 255, 0)
				} else if t == "orange" {
					PrintColored(Sprintf(".%s.", t),255, 255, 255, 255, 120, 0)
				} else if t == "yellow" {
					PrintColored(Sprintf(".%s.", t),0, 0, 0, 255, 255, 0)
				} else if t == "purple" {
					PrintColored(Sprintf(".%s.", t),255, 255, 255, 100, 10, 255)
				} else if t == "gray" || t == "grey" {
					PrintColored(Sprintf(".%s.", t),255, 255, 255, 100, 100, 100)
				} else if t == "black" {
					PrintColored(Sprintf(".%s.", t),255, 255, 255, 0, 0, 0)
				} else if t == "aqua" {
					PrintColored(Sprintf(".%s.", t),0, 0, 0, 0, 255, 255)
				} else if t == "teal" {
					PrintColored(Sprintf(".%s.", t),255, 255, 255, 0, 128, 128)
				} else if ! strings.ContainsRune(t, '/') {
					PrintColored(Sprintf(".%s.", t),255, 120, 0, 255, 255, 255)
				} else {
					Printf(".%s.", t)
				}
			}
		}
		Print("\n")
	}
}

