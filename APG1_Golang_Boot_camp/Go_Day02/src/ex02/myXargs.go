package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	args := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		args = append(args, line)
	}

	flags := os.Args[2:]
	command := os.Args[1]
	args = append(flags, args...)

	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error executing command:", err)
		os.Exit(1)
	}

}
