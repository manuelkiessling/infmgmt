package infrastructure

/*

- implementiert konkrete befehlsausführung, low level db zeugs etc.

*/

import (
	"fmt"
	"os/exec"
	"log"
)

type DefaultCommandExecutor struct{}

func (ce *DefaultCommandExecutor) Run(command string, arguments ...string) (output string, err error) {
	log.Printf("Now running: %s with args %+v\n", command, arguments)
	cmd := exec.Command(command, arguments...)

	outputBytes, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	output = fmt.Sprintf("%s", outputBytes)
	log.Printf("Output: %s\n", output)

	return output, err
}
