package domain

import (
	"reflect"
	"testing"
)

var vmguestsStoredInRepo []string

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
	var vmguests map[string]*Vmguest
	vmhost = &Vmhost{id, "vmhost" + id, vmguests}
	return vmhost, nil
}

func (repo *MockVmhostRepository) GetAll() (map[string]*Vmhost, error) {
	var vmguests map[string]*Vmguest
	vmhosts := make(map[string]*Vmhost)
	vmhosts["1"] = &Vmhost{"1", "vmhost1", vmguests}
	vmhosts["2"] = &Vmhost{"2", "vmhost2", vmguests}
	return vmhosts, nil
}

//func TestUpdateVmguestCache(t *testing.T) {
//	var expected []string
//	expected = append(expected, "1")
//	expected = append(expected, "2")
//	interactor := new(VmhostsInteractor)
//	interactor.UpdateVmguestCache() // muss für jeden vmhost einmal UpdateVmguestCache() aufrufen
//	if !reflect.DeepEqual(expected, updateVmguestCacheCalls) {
//	if allVmguestsWereStoredToRepo == false {
//		t.Errorf("Use case did not ")
//	}
//}

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
