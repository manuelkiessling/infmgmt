package infrastructure

import (
	"strings"
)

type MockCommandExecutor struct {
	Commandlines []string
}

func (ce *MockCommandExecutor) Run(command string, arguments ...string) (output string, err error) {
	commandline := command + " " + strings.Join(arguments, " ")
	if commandline == "/opt/infmgmt/bin/get_number_of_vmguests.sh vmhost1" {
		return "1", nil
	}
	if commandline == "/opt/infmgmt/bin/get_name_of_vmguest.sh vmhost1 0" {
		return "virtual1", nil
	}
	if commandline == "/opt/infmgmt/bin/get_state_of_vmguest.sh vmhost1 0" {
		return "running", nil
	}
	if commandline == "/opt/infmgmt/bin/get_uuid_of_vmguest.sh vmhost1 virtual1" {
		return "a0f39677-afda-f5bb-20b9-c5d8e3e06edf", nil
	}
	return "", nil
}
