package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	flagVerbose := false
	flag.BoolVar(&flagVerbose, "v", false, "verbose")
	flag.Parse()
	if flagVerbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(FormatLine(line))
	}
}
