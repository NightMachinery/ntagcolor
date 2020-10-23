/// 2>/dev/null ; exec gorun "$0" "$@"

package main

import (
	. "fmt"
	"io/ioutil"
	"log"
	"os"
	//"regexp"
	//"strconv"
	"strings"
	//"bufio"
)

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
				Printf(".%s.", t)
			}
		}
		Print("\n")
	}
}

