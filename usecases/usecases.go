package usecases

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
	"github.com/ManuelKiessling/infmgmt-backend/domain"
)

type OverviewListEntry struct {
	Id	string
	DnsName	string
}

type MachineOverviewInteractor struct {
	MachineRepository	domain.MachineRepository
}

func (interactor *MachineOverviewInteractor) List() (map[string]*OverviewListEntry, error) {
	var overviewListEntries map[string]*OverviewListEntry
  machines, nil := interactor.MachineRepository.GetAll()
	overviewListEntries = make(map[string] *OverviewListEntry, len(machines))
	for i, machine := range machines {
		overviewListEntries[i] = &OverviewListEntry{machine.Id, machine.DnsName}
	}
	return overviewListEntries, nil
}

