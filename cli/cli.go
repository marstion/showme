package cli

import (
	"bytes"
	"os"
	"os/exec"
)

func Cli(command string, params []string) (outStr string, err error) {
	cmd := exec.Command(command, params...)

	var out bytes.Buffer

	cmd.Stderr = os.Stderr
	cmd.Stdout = &out

	err = cmd.Start()
	if err != nil {
		return
	}
	err = cmd.Wait()
	return out.String(), err
}
