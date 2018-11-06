package runner

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

// Executor is a function capable of running the named program.
type Executor func(string, ...string) error

// Exec runs the named program and prints stdout and stderr.
func Exec(name string, args ...string) error {
	command := exec.Command(name, args...)

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
