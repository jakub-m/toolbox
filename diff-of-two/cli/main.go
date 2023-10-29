package main

// Large part generated with ChatGPT-3.5

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

//func main() {
//	edits := editDistance(str1, str2)
//	for _, edit := range edits {
//		fmt.Println(edit)
//	}
//
//}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lineCounter := 0
	lines := []string{"", ""}
	for ; lineCounter < len(lines); scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		log.Printf("line: %s", scanner.Text())
		lines[lineCounter%len(lines)] = scanner.Text()
		lineCounter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error reading input", err)
	}
	edits := editDistance(lines[0], lines[1])
	for _, line := range lines {
		fmt.Println(line)
	}
	// for _, edit := range edits {
	// 	fmt.Println(edit)
	// }
	fmt.Println(formatEdits(edits))
}

type EditType string

const (
	EditAdd    EditType = "add"
	EditRemove EditType = "remove"
	EditChange EditType = "change"
)

type Edit struct {
	editType EditType
	pos      int
	char     byte // should be rune
}

func editDistance(str1, str2 string) []Edit {
	// Create a matrix to store the edit distances between substrings
	matrix := make([][]int, len(str1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(str2)+1)
	}

	// Initialize the matrix with values for the base cases
	for i := 0; i <= len(str1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(str2); j++ {
		matrix[0][j] = j
	}

	// Fill in the matrix using dynamic programming
	for i := 1; i <= len(str1); i++ {
		for j := 1; j <= len(str2); j++ {
			cost := 0
			if str1[i-1] != str2[j-1] {
				cost = 1
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // Deletion
				matrix[i][j-1]+1,      // Insertion
				matrix[i-1][j-1]+cost, // Substitution
			)
		}
	}

	// Backtrack to find the edits
	edits := make([]Edit, 0)
	i, j := len(str1), len(str2)
	for i > 0 || j > 0 {
		currentCost := matrix[i][j]
		if i > 0 && matrix[i-1][j] == currentCost-1 {
			edits = append(edits, Edit{editType: EditRemove, pos: i - 1, char: str1[i-1]})
			i--
		} else if j > 0 && matrix[i][j-1] == currentCost-1 {
			edits = append(edits, Edit{editType: EditAdd, pos: j, char: str2[j-1]})
			j--
		} else {
			if currentCost > 0 {
				//edits = append(edits, fmt.Sprintf("Replace '%c' at position %d with '%c'", str1[i-1], i-1, str2[j-1]))
				edits = append(edits, Edit{editType: EditChange, pos: i - 1, char: str2[j-1]})
			}
			i--
			j--
		}
	}

	// Reverse the list of edits to get the correct order
	reverseEditSlice(edits)

	return edits
}

func min(a, b, c int) int {
	if a <= b && a <= c {
		return a
	} else if b <= a && b <= c {
		return b
	}
	return c
}

func reverseEditSlice(slice []Edit) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func formatEdits(edits []Edit) string {
	maxPos := 0
	for _, edit := range edits {
		if maxPos < edit.pos {
			maxPos = edit.pos
		}
	}
	buffer := make([]string, maxPos+1)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = " "
	}
	for _, edit := range edits {
		if edit.editType == EditAdd {
			buffer[edit.pos] = "+"
		} else if edit.editType == EditRemove {
			buffer[edit.pos] = "-"
		} else if edit.editType == EditChange {
			buffer[edit.pos] = "x"
		} else {
			panic(fmt.Sprintf("bad edit type: %+v", edit.editType))
		}
	}
	return strings.Join(buffer, "")
}
