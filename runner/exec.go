package runner

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

// Executor is a function capable of running the command.
type Executor func(*exec.Cmd) error

// Exec runs the command and prints stdout and stderr.
func Exec(command *exec.Cmd) error {
	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	command.Start()

	printPipe := func(pipe io.ReadCloser) {
		scanner := bufio.NewScanner(pipe)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}

	go printPipe(stdout)
	go printPipe(stderr)

	return command.Wait()
}
