package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ReadInput reads and cleans user input from the terminal.
func (ss *SchoolSystem) ReadInput() string {
	if !ss.ConsoleScanner.Scan() {
		return ""
	}

	return strings.TrimSpace(ss.ConsoleScanner.Text())

}

// DisplaySystemHeader shows a formatted section title.
func (ss *SchoolSystem) DisplaySystemHeader(title string) {
	fmt.Println("\n==================================================")
	fmt.Printf("  %s\n", strings.ToUpper(title))
	fmt.Println("==================================================")
}

// ClearTerminal clears the console screen before showing new content.
func (ss *SchoolSystem) ClearTerminal() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	_ = cmd.Run()

}
