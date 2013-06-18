package domain

import (
	"reflect"
	"testing"
)

type MockCommand struct {
	Command
	TheCommand string
}

type MockProcedure struct {
	Id                string
	Status            int
	OperationsHandler *MockOperationsHandler
	Started           bool
}

func (p *MockProcedure) Add(command Command) {
	p.OperationsHandler.Commands = append(p.OperationsHandler.Commands, command.TheCommand)
}

func (p *MockProcedure) Start() chan int {
	c := make(chan int)
	p.Status = 1
	p.Started = true
	c <- 1
	return c
}

type MockOperationsHandler struct {
	Commands []string
}

func (oh *MockOperationsHandler) NewProcedure() Procedure {
	p := new(MockProcedure)
	p.OperationsHandler = oh
	return p
}

func (oh *MockOperationsHandler) CommandCreateVirtualMachine(vmhostDnsName string, machineName string) (Command, error) {
	command := new(MockCommand)
	command.TheCommand = "CommandCreateVirtualMachine " + vmhostDnsName + " " + machineName
	return command, nil
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
	expectedCommands = append(expectedCommands, "CreateVirtualMachine Machine 1 Machine 2")

	oh := new(MockOperationsHandler)

	interactor := new(MachinesInteractor)
	interactor.MachineOperationsHandler = oh
	interactor.MachineRepository = new(MockRepository)
	interactor.SetupMachine("2")
	if !reflect.DeepEqual(expectedCommands, oh.Commands) {
		t.Errorf("Setup() did not execute OperationsHandler commands in the right order, calls were: %+v", oh.Commands)
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
