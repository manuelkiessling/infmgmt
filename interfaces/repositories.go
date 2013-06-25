package interfaces

import (
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/manuelkiessling/infmgmt-backend/domain"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
)

type VmguestRepository struct {
	commandExecutor CommandExecutor
}

type VmhostRepository struct {
	dbMap             *gorp.DbMap
	vmguestRepository *VmguestRepository
}

type vmhostModel struct {
	Id      string
	DnsName string
}

func NewVmguestRepository(commandExecutor CommandExecutor) *VmguestRepository {
	return &VmguestRepository{commandExecutor}
}

func (repo *VmguestRepository) GetAll(vmhostDnsName string) ([]*domain.Vmguest, error) {
	var output string
	var machineCount int
	var id, name, state string
	var command string
	var arguments []string
	var vmguests []*domain.Vmguest

	command = "ssh"
	arguments = append(arguments, "-i /home/manuel.kiessling/.ssh/infmgmt.id_rsa")
	arguments = append(arguments, "root@"+vmhostDnsName)
	arguments = append(arguments, strings.Split("virsh list --all | tail --lines=+3 | head --lines=-1 | wc -l", " ")...)
	output, _ = repo.commandExecutor.Run(command, arguments...)
	machineCount, _ = strconv.Atoi(strings.TrimSpace(output))

	for i := 0; i < machineCount; i++ {
		arguments = nil
		arguments = append(arguments, "-i /home/manuel.kiessling/.ssh/infmgmt.id_rsa")
		arguments = append(arguments, "root@"+vmhostDnsName)
		arguments = append(arguments, strings.Split("virsh list --all | tail --lines=+"+strconv.Itoa(3+i)+" | head --lines=1 | sed 's/ \\+/ /g' | cut -d' ' -f3", " ")...)
		output, _ = repo.commandExecutor.Run(command, arguments...)
		name = strings.TrimSpace(output)

		arguments = nil
		arguments = append(arguments, "-i /home/manuel.kiessling/.ssh/infmgmt.id_rsa")
		arguments = append(arguments, "root@"+vmhostDnsName)
		arguments = append(arguments, strings.Split("virsh list --all | tail --lines=+"+strconv.Itoa(3+i)+" | head --lines=1 | sed 's/ \\+/ /g' | cut -d' ' -f4-", " ")...)
		output, _ = repo.commandExecutor.Run(command, arguments...)
		state = strings.TrimSpace(output)

		arguments = nil
		arguments = append(arguments, "-i /home/manuel.kiessling/.ssh/infmgmt.id_rsa")
		arguments = append(arguments, "root@"+vmhostDnsName)
		arguments = append(arguments, strings.Split("virsh dumpxml "+name+" | grep uuid | cut --bytes=9-44", " ")...)
		output, _ = repo.commandExecutor.Run(command, arguments...)
		id = strings.TrimSpace(output)

		vmguest, _ := domain.NewVmguest(id, name, state)
		vmguests = append(vmguests, vmguest)
	}
	return vmguests, nil
}

func NewVmhostRepository(dbMap *gorp.DbMap, vmguestRepository *VmguestRepository) *VmhostRepository {
	// SetKeys(false) means we do have a primary key ("Id"), but we set it ourselves (no autoincrement)
	dbMap.AddTableWithName(vmhostModel{}, "vmhosts").SetKeys(false, "Id")
	repo := new(VmhostRepository)
	repo.dbMap = dbMap
	repo.vmguestRepository = vmguestRepository
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
		vmhost.Vmguests, _ = repo.vmguestRepository.GetAll(vmhost.DnsName)
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
		vmhost := repo.getVmhostFromVmhostModel(result)
		vmhost.Vmguests, _ = repo.vmguestRepository.GetAll(vmhost.DnsName)
		vmhosts[result.Id] = vmhost
	}
	return vmhosts, nil
}

func (repo *VmhostRepository) getVmhostFromVmhostModel(vm *vmhostModel) *domain.Vmhost {
	return &domain.Vmhost{Id: vm.Id, DnsName: vm.DnsName}
}
