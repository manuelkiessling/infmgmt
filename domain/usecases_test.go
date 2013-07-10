package domain

import (
	"reflect"
	"testing"
)

var vmhostRepositoryUpdateCacheWasCalled bool

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

func (oh *MockOperationsHandler) ExecuteProcedure(procedureId string) (chan int, error) {
	c := make(chan int)
	return c, nil
}

func (oh *MockOperationsHandler)	IsProcedureFinished(procedureId string) bool {
	return true
}

type MockVmhostRepository struct {
}

func (repo *MockVmhostRepository) Store(vmhost *Vmhost) error {
	return nil
}

func (repo *MockVmhostRepository) UpdateCache() error {
	vmhostRepositoryUpdateCacheWasCalled = true
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

func TestUpdateCache(t *testing.T) {
	oh := new(MockOperationsHandler)

	interactor := new(VmhostsInteractor)
	interactor.VmhostOperationsHandler = oh
	interactor.VmhostRepository = new(MockVmhostRepository)

	interactor.UpdateCache()
	if vmhostRepositoryUpdateCacheWasCalled == false {
		t.Errorf("interactor.UpdateCache did not call repo's UpdateCache")
	}
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
