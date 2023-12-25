package parse

import (
	"log"
	"strings"
)

var LogEnabled = false
var indent = ""
var logIndent = 0

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
		indent = strings.Repeat(".", logIndent)
	}
}
