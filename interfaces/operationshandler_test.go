package interfaces

import (
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
	//expected := "ssh root@kvmhost cp /var/lib/libvirt/images/infmgmgt-base.raw /var/lib/libvirt/images/newimage.raw"
	expected := "/usr/bin/touch /tmp/testfile-kvmhost1-virtual1"
	commandExecutor := new(MockCommandExecutor)
	oh := NewDefaultMachineOperationsHandler(commandExecutor)
	p := oh.NewProcedure()
	command, _ := oh.CommandCreateVirtualMachine("kvmhost1", "virtual1")
	p.Add(command)
	p.Start()
	for p.Status == 0 {}
	if commandExecutor.Commandlines[0] != expected {
		t.Errorf("OperationsHandler created commandline %+v, expected %+v", commandExecutor.Commandlines[0], expected)
	}
}

//func Test() {
//	operationsHandler := NewOperationsHandler()
//	procedure := operationsHandler.NewProcedure()
//	procedure.Add(operationsHandler.OperationCreateVirtualMachine("kvmhost1", "virtual1"))
//	procedure.Add(operationsHandler.OperationSetVirtualMachineHostname("virtual1"))
//	procedure.Run()
//	id := procedure.Id()
//	procedure := operationsHandler.getProcedure(Id)
//	status := procedure.Status()
//	output := procedure.Output()
//}
