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
	"fmt"
	"time"
)

type VmhostRepository interface {
	Store(vmhost *Vmhost) error
	FindById(id string) (*Vmhost, error)
	GetAll() (map[string]*Vmhost, error)
	UpdateCache() error
}

type VmhostOperationsHandler interface {
	InitializeProcedure() string
	AddCommandsCreateVirtualmachine(procedureId string, vmhostDnsName string, vmguestName string) error
	ExecuteProcedure(procedureId string) (chan int, error)
	IsProcedureFinished(procedureId string) bool
	//	SetIpAddressInGuestimage(vmhostDnsName string, vmguestName string, ipAddress string) error
	//	SetHostnameInGuestimage(vmhostDnsName string, vmguestName string, hostname string) error
}

type VmhostsListEntry struct {
	Id          string
	DnsName     string
	TotalMemory int
	Vmguests    map[string]*VmguestsListEntry
}

type VmguestsListEntry struct {
	Id              string
	Name            string
	State           string
	AllocatedMemory int
	InfoUpdatedAt		time.Time
}

type VmhostsInteractor struct {
	VmhostRepository        VmhostRepository
	VmhostOperationsHandler VmhostOperationsHandler
}

func (interactor *VmhostsInteractor) GetList() (map[string]*VmhostsListEntry, error) {
	var list map[string]*VmhostsListEntry
	vmhosts, err := interactor.VmhostRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve list of vmhosts from repository")
	}
	list = make(map[string]*VmhostsListEntry, len(vmhosts))
	for _, vmhost := range vmhosts {
		vmguestList, err := interactor.GetVmguestsList(vmhost.Id())
		if err != nil {
			return nil, fmt.Errorf("Error while trying to retrieve list of vmguests for vmhost %+v from repository", vmhost)
		}
		list[vmhost.Id()] = &VmhostsListEntry{vmhost.Id(), vmhost.DnsName(), vmhost.TotalMemory(), vmguestList}
	}
	return list, nil
}

func (interactor *VmhostsInteractor) GetVmguestsList(vmhostId string) (map[string]*VmguestsListEntry, error) {
	var list map[string]*VmguestsListEntry
	vmhost, err := interactor.VmhostRepository.FindById(vmhostId)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to retrieve vmhost with Id %s from repository", vmhostId)
	}
	list = make(map[string]*VmguestsListEntry, len(vmhost.Vmguests()))
	for _, vmguest := range vmhost.Vmguests() {
		list[vmguest.Id()] = &VmguestsListEntry{vmguest.Id(), vmguest.Name(), vmguest.State(), vmguest.AllocatedMemory(), vmguest.InfoUpdatedAt()}
	}
	return list, nil
}

func (interactor *VmhostsInteractor) CreateVmguest(vmhostId string, vmguestName string) (output string, err error) {
	vmhost, err := interactor.VmhostRepository.FindById(vmhostId)
	if err == nil {
		pId := interactor.VmhostOperationsHandler.InitializeProcedure()
		interactor.VmhostOperationsHandler.AddCommandsCreateVirtualmachine(pId, vmhost.DnsName(), vmguestName)
		// TODO: create vmguest entity as soon as all procedures are finished,
		//       add it to the vmguest cache (or get this by running the vmguest live repo?)
		//       maybe we should pass the channel into each ExecuteProcedure, this way it's
		//       simple to check if all procedures are finished?
		//       also, makes it simple to call procedures in order, not in parallel
		//       just give the channel a length of 1, as soon as one proc fills it, the next can start
		interactor.VmhostOperationsHandler.ExecuteProcedure(pId)
	} else {
		return "", err
	}
	return "", nil
}

func (interactor *VmhostsInteractor) UpdateCache() error {
	return interactor.VmhostRepository.UpdateCache()
}
