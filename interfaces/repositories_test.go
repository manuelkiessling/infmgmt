package interfaces

import (
	"github.com/ManuelKiessling/infmgmt-backend/domain"
	_ "log"
	_ "os"
	"testing"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/coopernurse/gorp"
)

func setupRepo() *MachineRepository {
	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-testdb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	//dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt-backend:", log.Lmicroseconds)) 
	repo := NewMachineRepository(dbMap)
	return repo
}

func TestMachineRepositoryStore(t *testing.T) {
	machine, _ := domain.NewMachine("foo", domain.P, nil)
	repo := setupRepo()
	err := repo.Store(machine)
	if err != nil {
		t.Errorf("Failed storing machine %+v in repo %+v", machine, repo)
	}
}

func TestMachineRepositoryFindById(t *testing.T) {
	machine, _ := domain.NewMachine("original", domain.P, nil)
	repo := setupRepo()
	repo.Store(machine)
	id := machine.Id
	dnsName := machine.DnsName
	machine.DnsName = "changed"
	repo = nil
	newRepo := setupRepo()
	retrievedMachine, err := newRepo.FindById(id)
	if err != nil {
		t.Errorf("Failed retrieving machine with id '%v' from repo %+v", id, repo)
		return
	}
	if retrievedMachine == nil {
		t.Errorf("Repo %+v did not return a machine", newRepo)
		return
	}
	if retrievedMachine.DnsName != dnsName {
		t.Errorf("Repo %+v did not return the correct machine: %+v", newRepo, retrievedMachine)
	}

}
