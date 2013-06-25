package domain

/*
- Bekommt Entitäts-Ids als Integer, setzt Domain Entitäten in Bwegung, indem es ggf. Funktionen der Entitäten mit
  Entitäten als Parameter aufruft statt mit Zahlen
- Triggert für den Case "erstelle neue VM" den Code im Interfaces Layer, der das Erstellen der VM auslöst
- endpunkt für events "vm installation abgeschlossen" oder "vmhost ist online", ändert entitäten entsprechend
- use case: umhängen einer vm auf andere pm
- starten, stoppen, zerstören von vms

*/

import (
	_ "fmt"
)

type VmhostRepository interface {
	Store(vmhost *Vmhost) error
	FindById(id string) (*Vmhost, error)
	GetAll() (map[string]*Vmhost, error)
}

type VmhostOperationsHandler interface {
	InitializeProcedure() string
	AddCommandsCreateVirtualmachine(procedureId string, vmhostDnsName string, vmguestName string) error
	ExecuteProcedure(procedureId string) chan int
	GetProcedureStatus(procedureId string) int
	//	SetIpAddressInGuestimage(vmhostDnsName string, vmguestName string, ipAddress string) error
	//	SetHostnameInGuestimage(vmhostDnsName string, vmguestName string, hostname string) error
}

type VmhostsListEntry struct {
	Id       string
	DnsName  string
	Vmguests map[string]*VmguestsListEntry
}

type VmguestsListEntry struct {
	Id    string
	Name  string
	State string
}

type VmhostsInteractor struct {
	VmhostRepository        VmhostRepository
	VmhostOperationsHandler VmhostOperationsHandler
}

func (interactor *VmhostsInteractor) GetList() (map[string]*VmhostsListEntry, error) {
	var list map[string]*VmhostsListEntry
	vmhosts, nil := interactor.VmhostRepository.GetAll()
	list = make(map[string]*VmhostsListEntry, len(vmhosts))
	for _, vmhost := range vmhosts {
		vmguestList, _ := interactor.GetVmguestsList(vmhost.Id)
		list[vmhost.Id] = &VmhostsListEntry{vmhost.Id, vmhost.DnsName, vmguestList}
	}
	return list, nil
}

func (interactor *VmhostsInteractor) GetVmguestsList(vmhostId string) (map[string]*VmguestsListEntry, error) {
	var list map[string]*VmguestsListEntry
	vmhost, nil := interactor.VmhostRepository.FindById(vmhostId)
	list = make(map[string]*VmguestsListEntry, len(vmhost.Vmguests))
	for _, vmguest := range vmhost.Vmguests {
		list[vmguest.Id] = &VmguestsListEntry{vmguest.Id, vmguest.Name, vmguest.State}
	}
	return list, nil
}

func (interactor *VmhostsInteractor) CreateVmguest(vmhostId string, vmguestName string) (output string, err error) {
	vmhost, err := interactor.VmhostRepository.FindById(vmhostId)
	if err == nil {
		pId := interactor.VmhostOperationsHandler.InitializeProcedure()
		interactor.VmhostOperationsHandler.AddCommandsCreateVirtualmachine(pId, vmhost.DnsName, vmguestName)
		interactor.VmhostOperationsHandler.ExecuteProcedure(pId)
	} else {
		return "", err
	}
	return "", nil
}
