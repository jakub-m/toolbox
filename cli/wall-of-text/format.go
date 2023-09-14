package main

import (
	"regexp"
	"strings"
)

var regexLineWithListedFields *regexp.Regexp

func init() {
	regexLineWithListedFields = regexp.MustCompile(`(\w+:\s+\w+)+`)
}

const newLine = "\n"

func FormatLine(line string) string {
	var output strings.Builder
	latestConsumedIndex := 0
	if allSubmatchIndices := regexLineWithListedFields.FindAllStringSubmatchIndex(line, -1); allSubmatchIndices != nil {
		for _, submatchIndices := range allSubmatchIndices {
			for i := 2; i < len(submatchIndices); i += 2 {
				indexStart := submatchIndices[i]
				indexEnd := submatchIndices[i+1]
				if indexStart > latestConsumedIndex {
					output.WriteString(line[latestConsumedIndex:indexStart])
				}
				latestConsumedIndex = indexEnd
				output.WriteString(newLine)
				output.WriteString(line[indexStart:indexEnd])
			}

		}
	}
	return output.String()
}
