package interfaces

import (
	"testing"
)

type MockCommandExecutor struct {
	Commandline string
}

func (ce *MockCommandExecutor) Run(commandline string) (output string, err error) {
	ce.Commandline = commandline
	return "", nil
}

func TestCreateGuestImageFromBaseImage(t *testing.T) {
	//expected := "ssh root@kvmhost cp /var/lib/libvirt/images/infmgmgt-base.raw /var/lib/libvirt/images/newimage.raw"
	expected := "/usr/bin/touch /tmp/testfile-kvmhost-newimage"
	commandExecutor := new(MockCommandExecutor)
	oh := NewDefaultMachineOperationsHandler(commandExecutor)
	oh.CreateGuestImageFromBaseImage("kvmhost", "newimage")
	if commandExecutor.Commandline != expected {
		t.Errorf("OperationsHandler created commandline %+v, expected %+v", commandExecutor.Commandline, expected)
	}
}
