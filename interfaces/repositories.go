package interfaces

import (
	"fmt"
	"github.com/manuelkiessling/infmgmt-backend/domain"
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
	var vm *vmhostModel
	vm = &vmhostModel{Id: vmhost.Id, DnsName: vmhost.DnsName}
	return repo.dbMap.Insert(vm)
}

func (repo *VmhostRepository) FindById(id string) (*domain.Vmhost, error) {
	var vmhost *domain.Vmhost
	var err error
	obj, err := repo.dbMap.Get(vmhostModel{}, id)
	if obj != nil {
		vm := obj.(*vmhostModel)
		vmhost = repo.getVmhostFromVmhostModel(vm)
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

func (repo *VmhostRepository) getVmhostFromVmhostModel(vm *vmhostModel) *domain.Vmhost {
	return &domain.Vmhost{Id: vm.Id, DnsName: vm.DnsName}
}
