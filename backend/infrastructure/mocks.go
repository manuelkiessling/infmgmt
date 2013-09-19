package infrastructure

import (
	"strings"
	"time"
)

type MockCommandExecutor struct {
	Commandlines []string
}

func (ce *MockCommandExecutor) Run(command string, arguments ...string) (output string, err error) {
	commandline := command + " " + strings.Join(arguments, " ")
  time.Sleep(100 * time.Millisecond)
	if commandline == "/usr/share/infmgmt/shellscripts/vmhostoperations/get_number_of_vmguests vmhost1" {
		return "1", nil
	}
	if commandline == "/usr/share/infmgmt/shellscripts/vmhostoperations/get_name_of_vmguest vmhost1 0" {
		return "virtual1", nil
	}
	if commandline == "/usr/share/infmgmt/shellscripts/vmhostoperations/get_state_of_vmguest vmhost1 0" {
		return "running", nil
	}
	if commandline == "/usr/share/infmgmt/shellscripts/vmhostoperations/get_uuid_of_vmguest vmhost1 virtual1" {
		return "a0f39677-afda-f5bb-20b9-c5d8e3e06edf", nil
	}
	if commandline == "/usr/share/infmgmt/shellscripts/vmhostoperations/get_allocated_memory_of_vmguest vmhost1 virtual1" {
		return "1048576", nil
	}
	if commandline == "/usr/share/infmgmt/shellscripts/vmhostoperations/get_memory_size_of_vmhost vmhost1" {
		return "32918292", nil
	}
	return "", nil
}
