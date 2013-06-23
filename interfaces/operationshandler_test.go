package interfaces

import (
	_ "fmt"
	"testing"
)

type MockCommandExecutor struct {
	Commandlines []string
}

func (ce *MockCommandExecutor) Run(command string, arguments ...string) (output string, err error) {
	for _, argument := range arguments {
		command = command + " " + argument
	}
	ce.Commandlines = append(ce.Commandlines, command)
	return "", nil
}

func TestCreateGuestImageFromBaseImage(t *testing.T) {
	expected0 := "ssh root@kvmhost1 'cp /var/lib/libvirt/images/infmgmt-base.raw /var/lib/libvirt/images/virtual1.raw'"
	expected1 := "ssh root@kvmhost1 'virt-install virtual1'"
	commandExecutor := new(MockCommandExecutor)
	oh := NewDefaultVmhostOperationsHandler(commandExecutor)
	procedureId := oh.InitializeProcedure()
	oh.AddCommandsCreateVirtualmachine(procedureId, "kvmhost1", "virtual1")
	c := oh.ExecuteProcedure(procedureId)
	status := 0
	for status == 0 {
		<-c
		status = oh.GetProcedureStatus(procedureId)
	}
	if commandExecutor.Commandlines[0] != expected0 {
		t.Errorf("OperationsHandler created commandline %+s, expected %+s", commandExecutor.Commandlines[0], expected0)
	}
	if commandExecutor.Commandlines[1] != expected1 {
		t.Errorf("OperationsHandler created commandline %+s, expected %+s", commandExecutor.Commandlines[1], expected1)
	}
}
