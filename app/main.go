package main

import (
	"fmt"
	"io"
	// Uncomment this block to pass the first stage!
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func copyFile(srcPath, destPath string) error {
	binSrc, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer binSrc.Close()

	binDest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer binDest.Close()

	// copy the file
	if _, err = io.Copy(binDest, binSrc); err != nil {
		return err
	}

	// set the permissions
	if err = binDest.Chmod(0755); err != nil {
		return err
	}

	return nil

}

func makeChroot(binaryFilePath string) error {

	// create a new temp directory for the chroot
	tempDirPath, err := ioutil.TempDir("", "my-docker")
	if err != nil {
		log.Fatalf("Error creating temp directory: %v", err)
	}
	// remove the temp directory when the program exits
	defer os.RemoveAll(tempDirPath)

	// create a temp /dev/null file
	nullFilePath := filepath.Join(tempDirPath, "dev/null")
	err = os.MkdirAll(nullFilePath, 0755)
	if err != nil {
		log.Fatalf("Error creating temp /dev/null file: %v", err)
	}
	defer os.RemoveAll(filepath.Dir(nullFilePath))

	// destination path for the binary file
	destPath := filepath.Join(tempDirPath, binaryFilePath)

	// create the destination directory
	err = os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		log.Fatalf("Error creating destination directory: %v", err)
	}

	// copy the binary file to the destination path
	err = copyFile(binaryFilePath, destPath)
	if err != nil {
		log.Fatalf("Error copying file: %v", err)
	}

	// chroot to the temp directory
	err = syscall.Chroot(tempDirPath)
	if err != nil {
		log.Fatalf("Error chrooting: %v", err)
	}

	return nil
}

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {

	command := os.Args[3]
	err := makeChroot(command)

	args := os.Args[4:len(os.Args)]

	cmd := exec.Command(command, args...)
	// pipe the stdout to the parent process
	// pipe the stderr to the parent process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		// get the exit code from the child process
		// exit with the same code
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		} else {
			fmt.Println(err)
			os.Exit(1)
		}

	}

}
