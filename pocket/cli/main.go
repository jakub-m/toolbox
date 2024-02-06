package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"slices"
	"strings"

	cp "github.com/otiai10/copy"
)

const defaultStashFilename = ".pocket.stash"

const helpString = `
pocket helps copying and moving files around. The basic usage is as follows

	ls | pocket
	pocket cp /some/destination
	pocket mv /other/destination

The possible commands are:

yank
	Stash the file names.

copy
	Copy the stashed files to destination.

move
	Move the stashed files to destination.
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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\n", strings.TrimSpace(helpString))
		flag.PrintDefaults()
	}
	flag.Parse()

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
	paths, err := readStashedPaths()
	if err != nil {
		return fmt.Errorf("failed reading stashed paths: %w", err)
	}
	fmt.Println(strings.Join(paths, "\n"))
	return nil
}

func runCommandCopy(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("expected destination path as the positional argument")
	}
	dirTo := args[0]
	copy2 := func(from, to string) error {
		return cp.Copy(from, to)
	}
	return forEachStashedPath(dirTo, copy2)
}

func runCommandMove(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("expected destination path as the positional argument")
	}
	dirTo := args[0]
	return forEachStashedPath(dirTo, os.Rename)
}

func forEachStashedPath(destinationDir string, fn func(sourcePath, destinationPath string) error) error {
	if err := ensurePathIsDirectory(destinationDir); err != nil {
		return err
	}
	paths, err := readStashedPaths()
	if err != nil {
		return fmt.Errorf("failed reading stashed paths: %w", err)
	}
	for _, pathFrom := range paths {
		pathTo := path.Join(destinationDir, path.Base(pathFrom))
		err := fn(pathFrom, pathTo)
		fmt.Fprintf(os.Stderr, "%s -> %s\n", pathFrom, pathTo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s -> %s failed: %s\n", pathFrom, destinationDir, err)
			//... but don't abort, continue.
		}
	}
	return nil
}

func ensurePathIsDirectory(path string) error {
	if info, err := os.Stat(path); err != nil {
		return fmt.Errorf("bad destination %s: %w", path, err)
	} else {
		if !info.IsDir() {
			return fmt.Errorf("destination %s is not a directory", path)
		}
	}
	return nil
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

func readStashedPaths() ([]string, error) {
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
