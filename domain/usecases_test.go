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
	var vmguests map[string]*Vmguest
	vmhost, _ := NewVmhost(id, "vmhost" + id)
	vmhost.SetTotalMemory(32918292)
	vmhost.SetVmguests(vmguests)
	return vmhost, nil
}

func (repo *MockVmhostRepository) GetAll() (map[string]*Vmhost, error) {
	var vmguests map[string]*Vmguest
	vmhosts := make(map[string]*Vmhost)
	vmhosts["1"], _ = NewVmhost("1", "vmhost1")
	vmhosts["1"].SetTotalMemory(32918292)
	vmhosts["1"].SetVmguests(vmguests)
	vmhosts["2"], _ = NewVmhost("2", "vmhost2")
	vmhosts["2"].SetTotalMemory(32918292)
	vmhosts["2"].SetVmguests(vmguests)
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
