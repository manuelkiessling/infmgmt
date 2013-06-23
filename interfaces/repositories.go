package interfaces

/*

- hier ist der befehlszeilen-zusammenbauer implementiert, der vom usecases layer benutzt wird um zB eine VM zu installieren
  oder puppet agent auszuführen usw.
- nutzt executor im infrastructure layer, um befehle tatsächlich auszuführen und deren ergebnis zu bekommen
- bekommt von aufrufer die infos als nackte daten, zB name der vm, größe arbeitsspeicher usw.
- webservice kann entscheiden einen json endpunkt anzubieten der alle vmhost infos auf einmal zurückgibt

*/

import (
	"fmt"
	"github.com/ManuelKiessling/infmgmt-backend/domain"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

type VmhostRepository struct {
	dbMap *gorp.DbMap
}

type vmhostModel struct {
	Id          string
	DnsName     string
}

func NewVmhostRepository(dbMap *gorp.DbMap) *VmhostRepository {
	// SetKeys(false) means we do have a primary key ("Id"), but we set it ourselves (no autoincrement)
	dbMap.AddTableWithName(vmhostModel{}, "vmhosts").SetKeys(false, "Id")
	repo := new(VmhostRepository)
	repo.dbMap = dbMap
	return repo
}

func (repo *VmhostRepository) Store(vmhost *domain.Vmhost) error {
	var mm *vmhostModel
	if vmhost.Vmhost == nil {
		mm = &vmhostModel{Id: vmhost.Id, DnsName: vmhost.DnsName}
	} else {
		repo.Store(vmhost.Vmhost)
		mm = &vmhostModel{Id: vmhost.Id, DnsName: vmhost.DnsName}
	}
	return repo.dbMap.Insert(mm)
}

func (repo *VmhostRepository) FindById(id string) (*domain.Vmhost, error) {
	var vmhost *domain.Vmhost
	var err error
	obj, err := repo.dbMap.Get(vmhostModel{}, id)
	if obj != nil {
		mm := obj.(*vmhostModel)
		vmhost = repo.getVmhostFromVmhostModel(mm)
	} else {
		vmhost = nil
		err = fmt.Errorf("No vmhost with id %v in repository", id)
	}
	return vmhost, err
}

func (repo *VmhostRepository) GetAll() (map[string]*domain.Vmhost, error) {
	var results []*vmhostModel
	vmhosts := make(map[string]*domain.Vmhost)
	query := "SELECT * FROM vmhosts ORDER BY Id"
	repo.dbMap.Select(&results, query)
	for _, result := range results {
		vmhosts[result.Id] = repo.getVmhostFromVmhostModel(result)
	}
	return vmhosts, nil
}

func (repo *VmhostRepository) getVmhostFromVmhostModel(mm *vmhostModel) *domain.Vmhost {
	var vmhost *domain.Vmhost
	if mm.VmhostId == "" {
		vmhost = nil
	} else {
		vmhost, _ = repo.FindById(mm.VmhostId)
	}
	return &domain.Vmhost{Id: mm.Id, DnsName: mm.DnsName}
}
