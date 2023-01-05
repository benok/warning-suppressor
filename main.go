package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// https://stackoverflow.com/questions/27576902/reading-stdout-from-a-subprocess
func filterOutput(scanner *bufio.Scanner, c color.Attribute) {
	for scanner.Scan() {
		line := scanner.Text()
		// TODO: read config file
		if strings.Contains(line, "LINN32:") {
			continue
		}
		colorize := color.New(c).SprintFunc()
		fmt.Println(colorize(line))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "warning-suppressor: scan failed(%s).\n", err)
	}
}

func fileNameWithoutExt(fn string) string {
	return fn[:len(fn)-len(filepath.Ext(fn))]
}

func main() {
	args := os.Args
	myCmd := args[0]
	// make full path of _orig command (abc.exe -> abc_orig.exe)
	origCmd := filepath.Join(fileNameWithoutExt(myCmd) + "_orig" + filepath.Ext(myCmd))

	cmd := exec.Command(origCmd, args[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning-suppressor: %s\n", err.Error())
		os.Exit(1)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning-suppressor: %s\n", err.Error())
		os.Exit(1)
	}

	err = cmd.Start()
	defer cmd.Wait()

	if err != nil {
		fmt.Fprintf(os.Stderr, "warning-suppressor: %s\n", err.Error())
		os.Exit(1)
	}

	s_out := bufio.NewScanner(stdout)
	s_err := bufio.NewScanner(stderr)

	go filterOutput(s_out, color.FgCyan)
	go filterOutput(s_err, color.FgRed)
}
