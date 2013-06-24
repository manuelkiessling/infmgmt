package interfaces

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/manuelkiessling/infmgmt-backend/domain"
	_ "github.com/mattn/go-sqlite3"
	_ "log"
	_ "os"
	"strconv"
	"strings"
	"testing"
)

type MockVmguestRepositoryCommandExecutor struct {
	Commandlines []string
}

func (ce *MockVmguestRepositoryCommandExecutor) Run(command string, arguments ...string) (output string, err error) {
	commandline := command + " " + strings.Join(arguments, " ")
	if commandline == "ssh -i /home/manuel.kiessling/.ssh/infmgmt.id_rsa root@vmhost1 virsh list --all | tail --lines=+3 | head --lines=-1 | wc -l" {
		return "1", nil
	}
	if commandline == "ssh -i /home/manuel.kiessling/.ssh/infmgmt.id_rsa root@vmhost1 virsh list --all | tail --lines=+3 | head --lines=1 | cut --bytes=8-38" {
		return "virtual1", nil
	}
	if commandline == "ssh -i /home/manuel.kiessling/.ssh/infmgmt.id_rsa root@vmhost1 virsh list --all | tail --lines=+3 | head --lines=1 | cut --bytes=39-52" {
		return "running", nil
	}
	if commandline == "ssh -i /home/manuel.kiessling/.ssh/infmgmt.id_rsa root@vmhost1 virsh dumpxml virtual1 | grep uuid | cut --bytes=9-44" {
		return "a0f39677-afda-f5bb-20b9-c5d8e3e06edf", nil
	}
	return "", nil
}

func setupVmguestRepo() *VmguestRepository {
	ce := new(MockVmguestRepositoryCommandExecutor)
	return NewVmguestRepository(ce)
}

func setupVmhostRepo() *VmhostRepository {
	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-testdb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	//dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt-backend:", log.Lmicroseconds))
	repo := NewVmhostRepository(dbMap, setupVmguestRepo())
	return repo
}

func (repo *VmhostRepository) reset() {
	repo.dbMap.DropTables()
	repo.dbMap.CreateTables()
}

func TestVmhostRepositoryStore(t *testing.T) {
	vmhost, _ := domain.NewVmhost("12345", "vmhost1", nil)
	repo := setupVmhostRepo()
	repo.reset()
	defer repo.reset()
	err := repo.Store(vmhost)
	if err != nil {
		t.Errorf("Failed storing vmhost %+v in repo %+v, error %+v", vmhost, repo, err)
	}
}

func TestVmhostRepositoryFindById(t *testing.T) {
	vmhost, _ := domain.NewVmhost("12345", "vmhost1", nil)
	repo := setupVmhostRepo()
	repo.reset()
	defer repo.reset()
	repo.Store(vmhost)
	id := vmhost.Id
	dnsName := vmhost.DnsName
	vmhost = nil
	repo = nil
	newRepo := setupVmhostRepo()
	retrievedVmhost, err := newRepo.FindById(id)
	if err != nil {
		t.Errorf("Failed retrieving vmhost with id '%v' from repo %+v", id, repo)
		return
	}
	if retrievedVmhost == nil {
		t.Errorf("Repo %+v did not return a vmhost", newRepo)
		return
	}
	if retrievedVmhost.DnsName != dnsName {
		t.Errorf("Repo %+v did not return the correct vmhost: %+v", newRepo, retrievedVmhost)
		return
	}
	vmguest := retrievedVmhost.Vmguests[0]
	if vmguest.Name != "virtual1" || vmguest.State != "running" || vmguest.Id != "a0f39677-afda-f5bb-20b9-c5d8e3e06edf" {
		t.Errorf("Repo %+v did not return a vmhost with correct vmguests: %+v", newRepo, retrievedVmhost.Vmguests[0])
		return
	}
}

func TestVmhostRepositoryGetAll(t *testing.T) {
	var vmhost *domain.Vmhost
	repo := setupVmhostRepo()
	repo.reset()
	defer repo.reset()
	for i := 0; i < 10; i++ {
		vmhost, _ = domain.NewVmhost(strconv.Itoa(i), "vmhost"+strconv.Itoa(i), nil)
		repo.Store(vmhost)
	}
	vmhosts, _ := repo.GetAll()
	if vmhosts[vmhost.Id].DnsName != "vmhost9" {
		t.Errorf("DnsName of retrieved vmhost %v (%+v) did not match DnsName of stored vmhost %+v", vmhost.Id, vmhosts[vmhost.Id], vmhost)
	}
}
