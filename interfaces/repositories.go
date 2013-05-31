package interfaces
/*

- hier ist der befehlszeilen-zusammenbauer implementiert, der vom usecases layer benutzt wird um zB eine VM zu installieren
  oder puppet agent auszuführen usw.
- nutzt executor im infrastructure layer, um befehle tatsächlich auszuführen und deren ergebnis zu bekommen
- bekommt von aufrufer die infos als nackte daten, zB name der vm, größe arbeitsspeicher usw.
- webservice kann entscheiden einen json endpunkt anzubieten der alle machine infos auf einmal zurückgibt

*/

import (
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/ManuelKiessling/infmgmt-backend/domain"
)

type Db interface {
	
}

type MachineRepository struct {
	dbMap *gorp.DbMap
}

func NewMachineRepository(dbMap *gorp.DbMap) *MachineRepository {
	// SetKeys(false) means we have a primary key, but we set it ourselves (no autoincrement)
	dbMap.AddTableWithName(domain.Machine{}, "machines").SetKeys(false, "Id")
	repo := new(MachineRepository)
	repo.dbMap = dbMap
	return repo
}

func (repo *MachineRepository) Store(machine *domain.Machine) error {
	return repo.dbMap.Insert(machine)
}

func (repo *MachineRepository) FindById(id string) (*domain.Machine, error) {
	var machine *domain.Machine
	var err error
	obj, err := repo.dbMap.Get(domain.Machine{}, id)
	if obj != nil {
		machine = obj.(*domain.Machine)
	} else {
		machine = nil
		err = fmt.Errorf("No machine with id %v in repository", id)
	}
	return machine, err
}
