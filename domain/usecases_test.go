package domain

import (
	"reflect"
	"testing"
)

type MockOperationsHandler struct {
	NumberOfCalls int
	Calls         map[string]int
}

func (oh *MockOperationsHandler) CreateGuestImageFromBaseImage(vmhostDnsName string, newImageName string) error {
	oh.Calls["CreateGuestImageFromBaseImage"] = oh.NumberOfCalls
	oh.NumberOfCalls++
	return nil
}

func (oh *MockOperationsHandler) SetIpAddressInGuestimage(vmhostDnsName string, vmguestName string, ipAddress string) error {
	oh.Calls["SetIpAddressInGuestimage"] = oh.NumberOfCalls
	oh.NumberOfCalls++
	return nil
}

func (oh *MockOperationsHandler) SetHostnameInGuestimage(vmhostDnsName string, vmguestName string, hostname string) error {
	oh.Calls["SetHostnameInGuestimage"] = oh.NumberOfCalls
	oh.NumberOfCalls++
	return nil
}

type MockRepository struct {
}

func (repo *MockRepository) Store(machine *Machine) error {
	return nil
}

func (repo *MockRepository) FindById(id string) (*Machine, error) {
	var machine *Machine
	if id == "1" { // asked for PM
		machine = &Machine{"1", "Machine 1", P, nil}
	}
	if id == "2" { // a VM with a PM
		vmhost := &Machine{"1", "Machine 1", P, nil}
		machine = &Machine{"2", "Machine 2", V, vmhost}
	}
	return machine, nil
}

func (repo *MockRepository) GetAll() (map[string]*Machine, error) {
	machines := make(map[string]*Machine)
	machines["1"] = &Machine{"1", "Machine 1", P, nil}
	machines["2"] = &Machine{"2", "Machine 2", V, machines["1"]}
	return machines, nil
}

func TestSetupMachineTriggersTheRightActions(t *testing.T) {
	expectedCalls := make(map[string]int)
	expectedCalls["CreateGuestImageFromBaseImage"] = 0
	//expectedCalls["SetIpAddressInGuestimage"] = 1
	//expectedCalls["SetHostnameInGuestimage"] = 2
	oh := new(MockOperationsHandler)
	oh.Calls = make(map[string]int)
	interactor := new(MachinesInteractor)
	interactor.MachineOperationsHandler = oh
	interactor.MachineRepository = new(MockRepository)
	interactor.SetupMachine("2")
	if !reflect.DeepEqual(expectedCalls, oh.Calls) {
		t.Errorf("Setup() did not execute OperationsHandler commands in the right order, calls were: %+v", oh.Calls)
	}
}

func TestSetupMachineSetsTheCorrectMachineState(t *testing.T) {

}

func TestSetupMachineFailsIfMachineIsNotVirtual(t *testing.T) {
	oh := new(MockOperationsHandler)
	oh.Calls = make(map[string]int)

	interactor := new(MachinesInteractor)
	interactor.MachineOperationsHandler = oh
	interactor.MachineRepository = new(MockRepository)
	_, err := interactor.SetupMachine("1")
	if err == nil {
		t.Errorf("Setting up a non-virtual machine should trigger an error, but it didn't")
	}
}
