package domain

import (
	"reflect"
	"testing"
)

type MockOperationsHandler struct {
	Commands []string
}

func (oh *MockOperationsHandler) InitializeProcedure() string {
	return "123"
}

func (oh *MockOperationsHandler) AddCommandsCreateVirtualmachine(procedureId string, vmhostDnsName string, virtualmachineName string) error {
	oh.Commands = append(oh.Commands, "CreateVirtualmachine "+vmhostDnsName+" "+virtualmachineName)
	return nil
}

func (oh *MockOperationsHandler) ExecuteProcedure(procedureId string) chan int {
	c := make(chan int)
	return c
}

func (oh *MockOperationsHandler) GetProcedureStatus(procedureId string) int {
	return 1
}

type MockVmhostRepository struct {
}

func (repo *MockVmhostRepository) Store(vmhost *Vmhost) error {
	return nil
}

func (repo *MockVmhostRepository) FindById(id string) (*Vmhost, error) {
	var vmhost *Vmhost
	if id == "1" {
		vmhost = &Vmhost{"1", "vmhost1"}
	}
	if id == "2" {
		vmhost = &Vmhost{"2", "vmhost2"}
	}
	return vmhost, nil
}

func (repo *MockVmhostRepository) GetAll() (map[string]*Vmhost, error) {
	vmhosts := make(map[string]*Vmhost)
	vmhosts["1"] = &Vmhost{"1", "vmhost1"}
	vmhosts["2"] = &Vmhost{"2", "vmhost2"}
	return vmhosts, nil
}

func TestSetupVmguestTriggersTheRightActions(t *testing.T) {
	expectedCommands := make([]string, 1)
	expectedCommands[0] = "CreateVirtualmachine vmhost1 vmguest1"

	oh := new(MockOperationsHandler)

	interactor := new(VmhostsInteractor)
	interactor.VmhostOperationsHandler = oh
	interactor.VmhostRepository = new(MockVmhostRepository)

	interactor.CreateVmguest("1", "vmguest1")
	if !reflect.DeepEqual(expectedCommands, oh.Commands) {
		t.Errorf("Setup() did not execute OperationsHandler commands in the right order, calls were: %+v but I expected %+v", oh.Commands, expectedCommands)
	}
}
