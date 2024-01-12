package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

func main() {
	yamlData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading YAML from stdin: %v\n", err)
		os.Exit(1)
	}

	var yamlMap any
	err = yaml.Unmarshal(yamlData, &yamlMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	if err := marshall(yamlMap, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error converting YAML to JSON: %v\n", err)
		os.Exit(1)
	}
}

func marshall(d any, out io.Writer) error {
	if d == nil {
		fmt.Fprintf(out, "null")
		return nil
	}
	switch value := d.(type) {
	case map[any]any:
		fmt.Fprintf(out, "{")
		i := 0
		for k, v := range value {
			marshall(k, out)
			fmt.Fprintf(out, ":")
			marshall(v, out)
			i++
			if i < len(value) {
				fmt.Fprintf(out, ",")
			}
		}
		fmt.Fprintf(out, "}")
	case []any:
		fmt.Fprintf(out, "[")
		for i, v := range value {
			marshall(v, out)
			if i < len(value)-1 {
				fmt.Fprintf(out, ",")
			}
		}
		fmt.Fprintf(out, "]")
	case string:
		fmt.Fprintf(out, "%s", strconv.Quote(value))
	case bool:
		if value {
			fmt.Fprintf(out, "true")
		} else {
			fmt.Fprintf(out, "false")
		}
	default:
		fmt.Fprintf(out, "%v", d)
	}
	return nil
}
