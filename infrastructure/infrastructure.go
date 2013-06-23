package infrastructure

/*

- implementiert konkrete befehlsausführung, low level db zeugs etc.

*/

import (
	"fmt"
	"os/exec"
)

type DefaultCommandExecutor struct{}

func (ce *DefaultCommandExecutor) Run(command string, arguments ...string) (output string, err error) {
	fmt.Printf("Now running: %s with args %+v", command, arguments)
	cmd := exec.Command(command, arguments...)

	outputBytes, err := cmd.Output()
	fmt.Printf("OutputBytes and error: %s - %+v", outputBytes, err)

	output = fmt.Sprintf("%s", outputBytes)
	fmt.Printf("Output: %s", output)
	return output, err
}
