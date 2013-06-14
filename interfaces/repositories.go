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
	"github.com/ManuelKiessling/infmgmt-backend/domain"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

type MachineRepository struct {
	dbMap *gorp.DbMap
}

type machineModel struct {
	Id          string
	DnsName     string
	MachineType int
	VmhostId    string
}

func NewMachineRepository(dbMap *gorp.DbMap) *MachineRepository {
	// SetKeys(false) means we do have a primary key ("Id"), but we set it ourselves (no autoincrement)
	dbMap.AddTableWithName(machineModel{}, "machines").SetKeys(false, "Id")
	repo := new(MachineRepository)
	repo.dbMap = dbMap
	return repo
}

func (repo *MachineRepository) Store(machine *domain.Machine) error {
	var mm *machineModel
	if machine.Vmhost == nil {
		mm = &machineModel{Id: machine.Id, DnsName: machine.DnsName, MachineType: machine.MachineType, VmhostId: ""}
	} else {
		repo.Store(machine.Vmhost)
		mm = &machineModel{Id: machine.Id, DnsName: machine.DnsName, MachineType: machine.MachineType, VmhostId: machine.Vmhost.Id}
	}
	return repo.dbMap.Insert(mm)
}

func (repo *MachineRepository) FindById(id string) (*domain.Machine, error) {
	var machine *domain.Machine
	var err error
	obj, err := repo.dbMap.Get(machineModel{}, id)
	if obj != nil {
		mm := obj.(*machineModel)
		machine = repo.getMachineFromMachineModel(mm)
	} else {
		machine = nil
		err = fmt.Errorf("No machine with id %v in repository", id)
	}
	return machine, err
}

func (repo *MachineRepository) GetAll() (map[string]*domain.Machine, error) {
	var results []*machineModel
	machines := make(map[string]*domain.Machine)
	query := "SELECT * FROM machines ORDER BY Id"
	repo.dbMap.Select(&results, query)
	for _, result := range results {
		machines[result.Id] = repo.getMachineFromMachineModel(result)
	}
	return machines, nil
}

func (repo *MachineRepository) getMachineFromMachineModel(mm *machineModel) *domain.Machine {
	var vmhost *domain.Machine
	if mm.VmhostId == "" {
		vmhost = nil
	} else {
		vmhost, _ = repo.FindById(mm.VmhostId)
	}
	return &domain.Machine{Id: mm.Id, DnsName: mm.DnsName, MachineType: mm.MachineType, Vmhost: vmhost}
}
