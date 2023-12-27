package parse

import (
	"fmt"
	"log"
	"strings"
)

var LogEnabled = false
var indent = ""
var logIndent = 0

func init() {
	updateIndent()
}

func Logf(fmt string, args ...any) {
	if !LogEnabled {
		return
	}
	log.Printf(indent+fmt, args...)
}

func LogIndentInc() {
	logIndent += 1
	updateIndent()
}

func LogIndentDec() {
	logIndent -= 1
	updateIndent()
}
func updateIndent() {
	if logIndent >= 0 {
		indent = fmt.Sprintf("%2d|%s", logIndent, strings.Repeat(" .", logIndent/2))
		if logIndent%2 == 1 {
			indent += " "
		}
	}
}
