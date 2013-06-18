package domain

/*
- Bekommt Entitäts-Ids als Integer, setzt Domain Entitäten in Bwegung, indem es ggf. Funktionen der Entitäten mit
  Entitäten als Parameter aufruft statt mit Zahlen
- Triggert für den Case "erstelle neue VM" den Code im Interfaces Layer, der das Erstellen der VM auslöst
- endpunkt für events "vm installation abgeschlossen" oder "machine ist online", ändert entitäten entsprechend
- use case: umhängen einer vm auf andere pm
- starten, stoppen, zerstören von vms

*/

import (
	_ "fmt"
	"errors"
)

type Command interface {}

type MachineRepository interface {
	Store(machine *Machine) error
	FindById(id string) (*Machine, error)
	GetAll() (map[string]*Machine, error)
}

type Procedure interface {
	Add(command Command)
	Start() chan int
}

type MachineOperationsHandler interface {
	NewProcedure() Procedure
	CommandCreateVirtualMachine(vmhostDnsName string, machineName string) (Command, error)
	//	SetIpAddressInGuestimage(vmhostDnsName string, vmguestName string, ipAddress string) error
	//	SetHostnameInGuestimage(vmhostDnsName string, vmguestName string, hostname string) error
}

type MachineOverviewListEntry struct {
	Id      string
	DnsName string
}

type MachinesInteractor struct {
	MachineRepository        MachineRepository
	MachineOperationsHandler MachineOperationsHandler
}

func (interactor *MachinesInteractor) ShowOverviewList() (map[string]*MachineOverviewListEntry, error) {
	var machineOverviewListEntries map[string]*MachineOverviewListEntry
	machines, nil := interactor.MachineRepository.GetAll()
	machineOverviewListEntries = make(map[string]*MachineOverviewListEntry, len(machines))
	for i, machine := range machines {
		machineOverviewListEntries[i] = &MachineOverviewListEntry{machine.Id, machine.DnsName}
	}
	return machineOverviewListEntries, nil
}

func (interactor *MachinesInteractor) SetupMachine(machineId string) (output string, err error) {
	machine, err := interactor.MachineRepository.FindById(machineId)
	if err == nil {
		if (machine.MachineType != V) {
			return "", errors.New("Can't setup a non-virtual machine")
		}
		p := interactor.MachineOperationsHandler.NewProcedure()
		command,_ := interactor.MachineOperationsHandler.CommandCreateVirtualMachine(machine.Vmhost.DnsName, machine.DnsName)
		p.Add(command)
		p.Start()
	} else {
		return "", err
	}
	return "", nil
}
