package main

import (
	"fmt"
	// Uncomment this block to pass the first stage!
	"os"
	"os/exec"
)

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {

	// Uncomment this block to pass the first stage!
	//
	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	cmd := exec.Command(command, args...)
	// pipe the stdout to the parent process
	// pipe the stderr to the parent process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}

}
