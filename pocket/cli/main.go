package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"slices"
	"strings"
)

const defaultStashFilename = ".pocket.stash"

const description = `To mark the files for further move or copy do

	ls | pocket yank 

To move the files later do

	pocket move .
`

var commandNamesCopy = []string{"c", "cp", "copy"}
var commandNamesMove = []string{"m", "mv", "move"}
var commandNamesYank = []string{"y", "yank"}

func main() {
	err := mainerr()
	if err != nil {
		log.Fatal(err)
	}
}
func mainerr() error {
	command := ""
	args := os.Args
	if len(args) > 1 {
		command = os.Args[1]
		args = args[2:]
	} else {
		args = args[1:]
	}
	if slices.Contains(commandNamesCopy, command) {
		return runCommandCopy(args)
	} else if slices.Contains(commandNamesMove, command) {
		return runCommandMove(args)
	} else if slices.Contains(commandNamesYank, command) {
		return runCommandYank(args)
	} else if command == "" {
		if isDataWaitingOnStdin() {
			return runCommandYank(args)
		} else {
			return runCommandPrint(args)
		}
	} else {
		return fmt.Errorf("unknown command %s", command)
	}
}

func runCommandYank(args []string) error {
	var paths []string
	if len(args) > 0 {
		paths = args
	} else {
		paths = readLines(os.Stdin)
	}
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to determine pwd: %w", err)
	}
	paths = normalizePaths(paths, pwd)
	stashPath, err := getStashPath()
	if err != nil {
		return fmt.Errorf("failed getting stash path: %w", err)
	}
	err = stashPaths(paths, stashPath)
	if err != nil {
		return fmt.Errorf("failed to save to stash: %w", err)
	}
	return nil
}

func runCommandPrint(args []string) error {
	lines, err := readStashedLines()
	if err != nil {
		return fmt.Errorf("failed reading stashed lines: %w", err)
	}
	fmt.Println(strings.Join(lines, "\n"))
	return nil
}

func runCommandCopy(args []string) error {
	return fmt.Errorf("not implemented")
}

func runCommandMove(args []string) error {
	return fmt.Errorf("not implemented")
}

func normalizePaths(paths []string, pwd string) []string {
	normalized := []string{}
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, string(os.PathSeparator)) {
			normalized = append(normalized, p)
		} else {
			normalized = append(normalized, path.Join(pwd, p))
		}
	}
	return normalized
}

func readStashedLines() ([]string, error) {
	p, err := getStashPath()
	if err != nil {
		return nil, fmt.Errorf("failed getting stash path: %w", err)
	}
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("failed opening stash: %w", err)
	}
	defer f.Close()
	return readLines(f), nil
}

func stashPaths(paths []string, stashFilePath string) error {
	f, err := os.Create(stashFilePath)
	if err != nil {
		return fmt.Errorf("failed opening stash file: %w", err)
	}
	defer f.Close()
	for _, p := range paths {
		fmt.Fprintln(f, p)
		fmt.Fprintln(os.Stdout, p)
	}
	return nil
}

func getStashPath() (string, error) {
	p, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(p, defaultStashFilename), nil
}

func readLines(r io.Reader) []string {
	lines := []string{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return lines
}

func isDataWaitingOnStdin() bool {
	stat, _ := os.Stdin.Stat()
	return ((stat.Mode() & os.ModeCharDevice) == 0)
}
