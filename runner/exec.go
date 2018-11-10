package runner

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

// Executor is a function capable of running the command and displaying its ouput.
type Executor func(*exec.Cmd, io.Writer) error

// Exec runs the command and prints stdout and stderr.
func Exec(command *exec.Cmd, writer io.Writer) error {
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
			fmt.Fprintln(writer, scanner.Text())
		}
	}

	go printPipe(stdout)
	go printPipe(stderr)

	return command.Wait()
}
