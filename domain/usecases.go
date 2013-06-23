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
	AddCommandCreateVirtualmachine(procedureId string, vmhostDnsName string, vmguestName string) error
	ExecuteProcedure(procedureId string) chan int
	GetProcedureStatus(procedureId string) int
	//	SetIpAddressInGuestimage(vmhostDnsName string, vmguestName string, ipAddress string) error
	//	SetHostnameInGuestimage(vmhostDnsName string, vmguestName string, hostname string) error
}

type VmhostsOverviewListEntry struct {
	Id      string
	DnsName string
}

type VmhostsInteractor struct {
	VmhostRepository        VmhostRepository
	VmhostOperationsHandler VmhostOperationsHandler
}

func (interactor *VmhostsInteractor) GetOverviewList() (map[string]*VmhostsOverviewListEntry, error) {
	var vmhostsOverviewListEntries map[string]*VmhostsOverviewListEntry
	vmhosts, nil := interactor.VmhostRepository.GetAll()
	vmhostsOverviewListEntries = make(map[string]*VmhostsOverviewListEntry, len(vmhosts))
	for i, vmhost := range vmhosts {
		vmhostsOverviewListEntries[i] = &VmhostsOverviewListEntry{vmhost.Id, vmhost.DnsName}
	}
	return vmhostsOverviewListEntries, nil
}

func (interactor *VmhostsInteractor) CreateVmguest(vmhostId string, vmguestName string) (output string, err error) {
	vmhost, err := interactor.VmhostRepository.FindById(vmhostId)
	if err == nil {
		pId := interactor.VmhostOperationsHandler.InitializeProcedure()
		interactor.VmhostOperationsHandler.AddCommandCreateVirtualmachine(pId, vmhost.DnsName, vmguestName)
		interactor.VmhostOperationsHandler.ExecuteProcedure(pId)
	} else {
		return "", err
	}
	return "", nil
}
