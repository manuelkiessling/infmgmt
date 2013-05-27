package interfaces
/*

- hier ist der befehlszeilen-zusammenbauer implementiert, der vom usecases layer benutzt wird um zB eine VM zu installieren
  oder puppet agent auszuführen usw.
- nutzt executor im infrastructure layer, um befehle tatsächlich auszuführen und deren ergebnis zu bekommen
- bekommt von aufrufer die infos als nackte daten, zB name der vm, größe arbeitsspeicher usw.
- webservice kann entscheiden einen json endpunkt anzubieten der alle machine infos auf einmal zurückgibt

*/

import (
	"github.com/ManuelKiessling/infmgmt-backend/domain"
)

type MachineRepository struct {
	Entries map[string]*domain.Machine
}

func NewMachineRepository() *MachineRepository {
	machineRepository := new(MachineRepository)
	machineRepository.Entries = make(map[string] *domain.Machine)
	return machineRepository
}

func (repo *MachineRepository) Store(machine *domain.Machine) error {
	repo.Entries[machine.Id] = machine
	return nil
}

