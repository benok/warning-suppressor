package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"golang.org/x/sys/windows"
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

func getOrigCommandLine() string {
	// https://github.com/golang/go/blob/f86f9a3038eb6db513a0ea36bc2af7a13b005e99/src/os/exec_windows.go#L96
	return windows.UTF16PtrToString(syscall.GetCommandLine())
}

func main() {
	args := os.Args
	myCmd := args[0]
	// make full path of _orig command (abc.exe -> abc_orig.exe)
	origCmd := filepath.Join(fileNameWithoutExt(myCmd) + "_orig" + filepath.Ext(myCmd))

	// execute abc_orig.exe
	cmd := exec.Command(origCmd /*,args[1:]...*/)

	// passing original command line is required for some tools
	// (see https://qiita.com/zetamatta/items/75ee1226f73113961f59,
	//	https://github.com/golang/go/blob/ea4631cc0cf301c824bd665a7980c13289ab5c9d/src/os/exec/exec.go#L373)
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: getOrigCommandLine()} // args[0] is not replaced, but it's OK.(not used)

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
