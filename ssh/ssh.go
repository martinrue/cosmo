package ssh

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

// Client wraps access to SSH.
type Client struct {
	Host string
}

// Exec executes the command on the remote host and returns the full response.
func (c *Client) Exec(cmd string) (string, error) {
	command := exec.Command("ssh", c.Host, cmd)

	buffer := &bytes.Buffer{}
	command.Stdout = buffer

	if err := command.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(buffer.String()), nil
}

// ExecStream executes the command on the remote host and streams the response.
func (c *Client) ExecStream(cmd string, onData func(string)) error {
	command := exec.Command("ssh", c.Host, cmd)

	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	command.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		onData(scanner.Text())
	}

	command.Wait()

	return nil
}
