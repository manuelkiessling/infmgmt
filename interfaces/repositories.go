package interfaces

import (
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/manuelkiessling/infmgmt-backend/domain"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
)

var cacheUpdateRunning = false

type VmguestRepository interface {
	Store(vmhostDnsName string, vmguest *domain.Vmguest) error
	GetAll(vmhostDnsName string) (map[string]*domain.Vmguest, error)
}

type VmguestLiveRepository struct {
	commandExecutor CommandExecutor
}

type VmguestCacheRepository struct {
	dbMap *gorp.DbMap
}

type vmguestModel struct {
	Id            string
	VmhostDnsName string
	Name          string
	State         string
}

type VmhostRepository struct {
	dbMap                  *gorp.DbMap
	vmguestLiveRepository  VmguestRepository
	vmguestCacheRepository VmguestRepository
}

type vmhostModel struct {
	Id      string
	DnsName string
}

func NewVmguestLiveRepository(commandExecutor CommandExecutor) VmguestRepository {
	return &VmguestLiveRepository{commandExecutor}
}

func NewVmguestCacheRepository(dbMap *gorp.DbMap) VmguestRepository {
	dbMap.AddTableWithName(vmguestModel{}, "vmguests").SetKeys(false, "Id")
	return &VmguestCacheRepository{dbMap}
}

func (repo *VmguestLiveRepository) Store(vmhostDnsName string, vmguest *domain.Vmguest) error {
	return fmt.Errorf("Live Vmguest Repo cannot store")
}

func (repo *VmguestLiveRepository) GetAll(vmhostDnsName string) (map[string]*domain.Vmguest, error) {
	var output string
	var machineCount int
	var id, name, state string
	var command string
	var arguments []string
	var vmguests map[string]*domain.Vmguest

	vmguests = make(map[string]*domain.Vmguest)

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
		vmguests[id] = vmguest
	}
	return vmguests, nil
}

func (repo *VmguestCacheRepository) Store(vmhostDnsName string, vmguest *domain.Vmguest) error {
	var vm *vmguestModel
	vm = &vmguestModel{Id: vmguest.Id, VmhostDnsName: vmhostDnsName, Name: vmguest.Name, State: vmguest.State}
	return repo.dbMap.Insert(vm)
}

func (repo *VmguestCacheRepository) GetAll(vmhostDnsName string) (map[string]*domain.Vmguest, error) {
	var results []*vmguestModel
	vmguests := make(map[string]*domain.Vmguest)
	query := "SELECT * FROM vmguests WHERE VmhostDnsName = ?"
	repo.dbMap.Select(&results, query, vmhostDnsName)
	for _, result := range results {
		vmguests[result.Id] = repo.getVmguestFromVmguestModel(result)
	}
	return vmguests, nil
}

func (repo *VmguestCacheRepository) getVmguestFromVmguestModel(vm *vmguestModel) *domain.Vmguest {
	return &domain.Vmguest{Id: vm.Id, Name: vm.Name, State: vm.State}
}

func NewVmhostRepository(dbMap *gorp.DbMap, vmguestLiveRepository VmguestRepository, vmguestCacheRepository VmguestRepository) *VmhostRepository {
	// SetKeys(false) means we do have a primary key ("Id"), but we set it ourselves (no autoincrement)
	dbMap.AddTableWithName(vmhostModel{}, "vmhosts").SetKeys(false, "Id")
	repo := new(VmhostRepository)
	repo.dbMap = dbMap
	repo.vmguestLiveRepository = vmguestLiveRepository
	repo.vmguestCacheRepository = vmguestCacheRepository
	return repo
}

func (repo *VmhostRepository) Store(vmhost *domain.Vmhost) error {
	var vm *vmhostModel
	if vmhost.Id == "" {
		return errors.New("Cannot store vmhosts with an empty Id")
	}
	if vmhost.DnsName == "" {
		return errors.New("Cannot store vmhosts with an empty DnsName")
	}
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
		vmhost.Vmguests, _ = repo.vmguestCacheRepository.GetAll(vmhost.DnsName)
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
		vmguests, err := repo.vmguestCacheRepository.GetAll(vmhost.DnsName)
		if err != nil {
			return nil, fmt.Errorf("Error loading vmguests for vmhost (%+v)", err)
		}
		vmhost.Vmguests = vmguests
		vmhosts[result.Id] = vmhost
	}
	return vmhosts, nil
}

func (repo *VmhostRepository) UpdateCache() error {
	if cacheUpdateRunning {
		return fmt.Errorf("A cache update is already running, try again later")
	}

	vmhosts, err := repo.GetAll()
	if err != nil {
		fmt.Errorf("Error getting vmhosts for UpdateCache")
	}

	cacheUpdateRunning = true

	go func() {
		for _, vmhost := range vmhosts {
			vmguests, err := repo.vmguestLiveRepository.GetAll(vmhost.DnsName)
			if err != nil {
				fmt.Errorf("Error getting vmguests for UpdateCache")
			}
			for _, vmguest := range vmguests {
				repo.vmguestCacheRepository.Store(vmhost.DnsName, vmguest)
			}
		}
		cacheUpdateRunning = false
	}()

	return nil
}

func (repo *VmhostRepository) getVmhostFromVmhostModel(vm *vmhostModel) *domain.Vmhost {
	return &domain.Vmhost{Id: vm.Id, DnsName: vm.DnsName}
}
