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

func (oh *MockOperationsHandler) AddCommandCreateVirtualMachine(procedureId string, vmhostDnsName string, machineName string) error {
	oh.Commands = append(oh.Commands, "CreateVirtualMachine " + vmhostDnsName + " " + machineName)
	return nil
}

func (oh *MockOperationsHandler) ExecuteProcedure(procedureId string) chan int {
	c := make(chan int)
	return c
}

func (oh *MockOperationsHandler) GetProcedureStatus(procedureId string) int {
	return 1
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
	expectedCommands := make([]string, 1)
	expectedCommands[0] = "CreateVirtualMachine Machine 1 Machine 2"

	oh := new(MockOperationsHandler)

	interactor := new(MachinesInteractor)
	interactor.MachineOperationsHandler = oh
	interactor.MachineRepository = new(MockRepository)
	interactor.SetupMachine("2")
	if !reflect.DeepEqual(expectedCommands, oh.Commands) {
		t.Errorf("Setup() did not execute OperationsHandler commands in the right order, calls were: %+v but I expected %+v", oh.Commands, expectedCommands)
	}
}

func TestSetupMachineSetsTheCorrectMachineState(t *testing.T) {

}

func TestSetupMachineFailsIfMachineIsNotVirtual(t *testing.T) {
	interactor := new(MachinesInteractor)
	interactor.MachineOperationsHandler = new(MockOperationsHandler)
	interactor.MachineRepository = new(MockRepository)
	_, err := interactor.SetupMachine("1")
	if err == nil {
		t.Errorf("Setting up a non-virtual machine should trigger an error, but it didn't")
	}
}
