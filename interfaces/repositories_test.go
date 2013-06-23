package interfaces

import (
	"database/sql"
	"github.com/manuelkiessling/infmgmt-backend/domain"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
	"testing"
	"reflect"
)

type MockVmguestRepositoryCommandExecutor struct {
	Commandlines []string
}

func (ce *MockVmguestRepositoryCommandExecutor) Run(command string, arguments ...string) (output string, err error) {
	var expectedArguments []string
	expectedArguments = make([]string, 0)
	expectedArguments = append(expectedArguments, "root@kvmhost1")
	expectedArguments = append(expectedArguments, "'virsh list --all'")
	if command == "ssh" && reflect.DeepEqual(arguments, expectedArguments) {
		return " Id    Name                           State\n----------------------------------------------------\n 2     virtual1                           running", nil
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
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt-backend:", log.Lmicroseconds))
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
