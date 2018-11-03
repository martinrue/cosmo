package runner

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

// Runner is a type capable of running a script.
type Runner interface {
	Run(script string) error
}

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
