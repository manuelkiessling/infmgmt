package domain

/*

- implementiert business entities, zustände, rules: zB "eine VM machine darf nur einer VM-Host machine zugeordnet sein
- implementiert "setze zustand dieser machine auf 'online', getriggert durch use case

*/

import (
	"errors"
	"github.com/streadway/simpleuuid"
	"time"
)

const (
	P = 0
	V = 1
)

type MachineRepository interface {
	Store(machine *Machine) error
	FindById(id string) (*Machine, error)
	GetAll() (map[string]*Machine, error)
}

type MachineOperationsHandler interface {
	Setup(machine *Machine) error
	Reboot(machine *Machine) error
}

type Machine struct {
	Id			string
	DnsName	string
	MachineType int
	Vmhost *Machine
	operationsHandler MachineOperationsHandler
}

func NewMachine(name string, machineType int, vmhost *Machine, operationsHandler MachineOperationsHandler) (newMachine *Machine, err error) {
	if machineType == V {
		if vmhost == nil {
			return nil, errors.New("Must attach virtual machine to existing physical machine")
		}
		if vmhost.MachineType == V {
			return nil, errors.New("Cannot attach virtual machine to other virtual machine")
		}
	}
	if machineType == P {
		if vmhost != nil {
			return nil, errors.New("Cannot attach physical machine to other machine")
		}
	}
	
	uuid, _ := simpleuuid.NewTime(time.Now())
	id := uuid.String()

	machine := &Machine{id, name, machineType, vmhost, operationsHandler}
	return machine, nil
}

func (machine *Machine) Setup() error {
	machine.operationsHandler.Setup(machine)
	return nil
}

func (machine *Machine) Reboot() error {
	machine.operationsHandler.Reboot(machine)
	return nil
}
