package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var verbose bool
	var textOnly bool
	flag.BoolVar(&verbose, "v", false, "verbose, useful debug")
	flag.BoolVar(&textOnly, "m", false, "message text only")
	flag.Usage = func() {
		fmt.Println("Parse cloud log in JSON format and output nice readable log lines.")
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
	if !verbose {
		log.SetOutput(ioutil.Discard)
	}
	inbytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatalln("read:", err)
	}
	objects, err := unmarshall(inbytes)
	if err != nil {
		fatalln("unmarshall:", err)
	}
	log.Println("unmarshalled objects:", len(objects))
	formatErrorCount := 0
	for _, obj := range objects {
		if s, err := format(obj, textOnly); err == nil {
			fmt.Println(s)
		} else {
			log.Println("format:", err)
			formatErrorCount++
		}
	}
	log.Println("format error count:", formatErrorCount)
	if formatErrorCount > 0 {
		fmt.Fprintf(os.Stderr, "There were %d format error(s)\n", formatErrorCount)
	}
}

func unmarshall(raw []byte) ([]any, error) {
	if obj, err := unmarshallSingle(raw); err == nil {
		if arr, ok := obj.([]any); ok {
			return arr, nil
		} else {
			return []any{obj}, nil
		}
	}
	return unmarshalOnePerLine(raw)
}

func unmarshalOnePerLine(raw []byte) ([]any, error) {
	objects := []any{}
	rawString := string(raw)
	rawString = strings.Trim(rawString, "\n")
	for _, line := range strings.Split(rawString, "\n") {
		obj, err := unmarshallSingle([]byte(line))
		if err != nil {
			return objects, err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

func unmarshallSingle(raw []byte) (any, error) {
	log.Printf("unmarshall single: %s", raw)
	var obj any
	err := json.Unmarshal(raw, &obj)
	return obj, err
}

func format(obj any, textOnly bool) (string, error) {
	if s, err := formatCloudLog(obj, textOnly); err == nil {
		return s, nil
	}
	return "", fmt.Errorf("cannot format")
}

func formatCloudLog(obj any, textOnly bool) (string, error) {
	timestamp, err := getField(obj, "timestamp")
	if err != nil {
		return "", err
	}
	severity, err := getField(obj, "severity")
	if err != nil {
		return "", err
	}
	text := ""
	if message, err := getField(obj, "jsonPayload", "message"); err != nil {
		if textPayload, err := getField(obj, "textPayload"); err != nil {
			return "", err
		} else {
			text = fmt.Sprint(textPayload)
		}
	} else {
		text = fmt.Sprint(message)
	}
	outputLine := ""
	if textOnly {
		outputLine = text
	} else {
		outputLine = fmt.Sprintf("%s\t%s\t%s", timestamp, severity, text)
	}
	return outputLine, nil
}

func getField(obj any, keys ...string) (any, error) {
	if len(keys) == 0 {
		return obj, nil
	}
	m, ok := obj.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("not a map: %s", obj)
	}
	return getField(m[keys[0]], keys[1:]...)
}

func fatalln(vals ...any) {
	fmt.Fprintln(os.Stderr, vals...)
	os.Exit(1)
}
