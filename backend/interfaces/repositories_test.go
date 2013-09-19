package interfaces

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/manuelkiessling/infmgmt/backend/domain"
	"github.com/manuelkiessling/infmgmt/backend/infrastructure"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
	"fmt"
)

func setupVmguestLiveRepo() VmguestRepository {
	ce := new(infrastructure.MockCommandExecutor)
	return NewVmguestLiveRepository(ce)
}

func setupVmguestCacheRepo() VmguestRepository {
	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-testdb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt/backend:", log.Lmicroseconds))
	dbMap.TraceOff()
	repo := NewVmguestCacheRepository(dbMap)
	dbMap.DropTables()
	dbMap.CreateTables()
	return repo
}

func setupVmhostRepo() *VmhostRepository {
	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-testdb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt/backend:", log.Lmicroseconds))
	dbMap.TraceOff()
	repo := NewVmhostRepository(dbMap, setupVmguestLiveRepo(), setupVmguestCacheRepo())
	return repo
}

func (repo *VmhostRepository) reset() {
	repo.dbMap.DropTables()
	repo.dbMap.CreateTables()
}

func TestVmhostRepositoryStore(t *testing.T) {
	vmhost, _ := domain.NewVmhost("12345", "vmhost1")
	vmhost.SetTotalMemory(32918292)
	repo := setupVmhostRepo()
	repo.reset()
	defer repo.reset()
	err := repo.Store(vmhost)
	if err != nil {
		t.Errorf("Failed storing vmhost %+v in repo %+v, error %+v", vmhost, repo, err)
	}
	vmhost, _ = repo.FindById("12345")
	if vmhost.TotalMemory() != 32918292 {
		t.Errorf("Vmhost repository did not store TotalMemory value")
	}
	vmhost.SetTotalMemory(238934)
	err = repo.Store(vmhost)
	if err != nil {
		t.Errorf("Failed re-storing vmhost %+v in repo %+v, error %+v", vmhost, repo, err)
	}
}

func TestVmhostRepositoryStoreFailsWithEmptyId(t *testing.T) {
	vmhost, _ := domain.NewVmhost("", "vmhost1")
	vmhost.SetTotalMemory(32918292)
	repo := setupVmhostRepo()
	repo.reset()
	defer repo.reset()
	err := repo.Store(vmhost)
	if err == nil {
		t.Errorf("Could store vmhost %+v in repo %+v although its Id was empty", vmhost, repo)
	}
}

func TestVmhostRepositoryStoreFailsWithEmptyDnsName(t *testing.T) {
	vmhost, _ := domain.NewVmhost("12345", "")
	vmhost.SetTotalMemory(32918292)
	repo := setupVmhostRepo()
	repo.reset()
	defer repo.reset()
	err := repo.Store(vmhost)
	if err == nil {
		t.Errorf("Could store vmhost %+v in repo %+v although its DnsName was empty", vmhost, repo)
	}
}

func TestVmhostRepositoryFindById(t *testing.T) {
	vmhost, _ := domain.NewVmhost("12345", "vmhost1")
	vmhost.SetTotalMemory(32918292)
	repo := setupVmhostRepo()
	repo.reset()
	defer repo.reset()
	repo.Store(vmhost)
	id := vmhost.Id()
	dnsName := vmhost.DnsName()
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
	if retrievedVmhost.DnsName() != dnsName {
		t.Errorf("Repo %+v did not return the correct vmhost: %+v", newRepo, retrievedVmhost)
		return
	}
}

func TestVmhostRepositoryGetAll(t *testing.T) {
	var vmhost *domain.Vmhost
	repo := setupVmhostRepo()
	repo.reset()
	defer repo.reset()
	for i := 0; i < 10; i++ {
		vmhost, _ = domain.NewVmhost(strconv.Itoa(i), "vmhost"+strconv.Itoa(i))
		vmhost.SetTotalMemory(32918292)
		repo.Store(vmhost)
	}
	vmhosts, _ := repo.GetAll()
	if vmhosts[vmhost.Id()].DnsName() != "vmhost9" {
		t.Errorf("DnsName of retrieved vmhost %v (%+v) did not match DnsName of stored vmhost %+v", vmhost.Id(), vmhosts[vmhost.Id()], vmhost)
	}
}

func TestVmhostRepositoryUpdateCache(t *testing.T) {
	repo := setupVmhostRepo()
	repo.reset()
	defer repo.reset()

	vmhost, _ := domain.NewVmhost("12345", "vmhost1")
	vmhost.SetTotalMemory(32918292)
	repo.Store(vmhost)
	vmhost, _ = repo.FindById("12345")

	if len(vmhost.Vmguests()) != 0 {
		t.Errorf("Should be empty: %+v", vmhost.Vmguests())
	}

	repo.UpdateCache()

	vmhost, _ = repo.FindById("12345")
	vmguest := vmhost.Vmguests()["a0f39677-afda-f5bb-20b9-c5d8e3e06edf"]

	vmguestInfoUpdatedAt := vmguest.InfoUpdatedAt()
	durationSinceVmguestWasUpdated := time.Since(vmguestInfoUpdatedAt)

	if durationSinceVmguestWasUpdated.Seconds() > 5.0 {
		t.Errorf("The repo says that vmguest %+v wasn't recently updated. It was last updated %+v ago.", vmguest, durationSinceVmguestWasUpdated)
	}

	if vmguest.Name() != "virtual1" || vmguest.State() != "running" || vmguest.Id() != "a0f39677-afda-f5bb-20b9-c5d8e3e06edf" || vmguest.AllocatedMemory() != 1048576 {
		t.Errorf("Repo %+v did not return a vmhost with correct vmguests: %+v", repo, vmhost.Vmguests()["a0f39677-afda-f5bb-20b9-c5d8e3e06edf"])
		return
	}
}

func TestVmhostRepositoryUpdateCacheIsLocked(t *testing.T) {
	repo := setupVmhostRepo()
	repo.reset()
	defer repo.reset()

	vmhost, _ := domain.NewVmhost("12345", "vmhost1")
	vmhost.SetTotalMemory(32918292)
	repo.Store(vmhost)

	go func() {
		repo.UpdateCache()
	}()
	fmt.Printf("started UpdateCache in goroutine")
	err := repo.UpdateCache()
	if err == nil {
		t.Errorf("Expected repo.UpdateCache() to fail because another process is running")
	}
}
