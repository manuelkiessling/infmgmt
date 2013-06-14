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

type Machine struct {
	Id          string
	DnsName     string
	MachineType int
	Vmhost      *Machine
}

func NewMachine(name string, machineType int, vmhost *Machine) (newMachine *Machine, err error) {
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

	machine := &Machine{id, name, machineType, vmhost}
	return machine, nil
}