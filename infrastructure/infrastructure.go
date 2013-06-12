package infrastructure

/*

- implementiert konkrete befehlsausführung, low level db zeugs etc.

*/

import (
	"fmt"
	"os/exec"
)

type DefaultCommandExecutor struct{}

func (ce *DefaultCommandExecutor) Run(commandline string) (output string, err error) {
	cmd := exec.Command(commandline)
	outputBytes, err := cmd.Output()
	output = fmt.Sprintf("%s", outputBytes)
	//fmt.Printf("%s", output)
	return output, err
}
