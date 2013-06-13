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
	vmhost, _ := NewMachine("Vmhost 12345", P, nil)
	machine, _ := NewMachine("Machine " + id, V, vmhost)
	return machine, nil
}

func (repo *MockRepository) GetAll() (map[string]*Machine, error) {
	machines := make(map[string] *Machine)
	machines["101"] = &Machine{"101", "Mocked machine #1", P, nil}
	machines["102"] = &Machine{"102", "Mocked machine #2", P, nil}
	return machines, nil
}

func TestSetupMachine(t *testing.T) {
	expectedCalls := make(map[string]int)
	expectedCalls["CreateGuestImageFromBaseImage"] = 0
//	expectedCalls["SetIpAddressInGuestimage"] = 1
//	expectedCalls["SetHostnameInGuestimage"] = 2
	oh := new(MockOperationsHandler)
	oh.Calls = make(map[string]int)
	interactor := new(MachinesInteractor)
	interactor.MachineOperationsHandler = oh
	interactor.MachineRepository = new(MockRepository)
	interactor.SetupMachine("111")
	if !reflect.DeepEqual(expectedCalls, oh.Calls) {
		t.Errorf("Setup() did not execute OperationsHandle commands in the right order, calls were: %+v", oh.Calls)
	}
}

