package main

import (
	"database/sql"
	_ "fmt"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
	"github.com/manuelkiessling/infmgmt-backend/domain"
	"github.com/manuelkiessling/infmgmt-backend/interfaces"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var numberOfCommandCalls int

type MockVmguestRepositoryCommandExecutor struct {
	Commandlines []string
}

func (ce *MockVmguestRepositoryCommandExecutor) Run(command string, arguments ...string) (output string, err error) {
	numberOfCommandCalls++
	commandline := command + " " + strings.Join(arguments, " ")
	if commandline == "ssh -i /home/manuel.kiessling/.ssh/infmgmt.id_rsa root@vmhost1 virsh list --all | tail --lines=+3 | head --lines=-1 | wc -l" {
		return "1", nil
	}
	if commandline == "ssh -i /home/manuel.kiessling/.ssh/infmgmt.id_rsa root@vmhost1 virsh list --all | tail --lines=+3 | head --lines=1 | sed 's/ \\+/ /g' | cut -d' ' -f3" {
		return "virtual1", nil
	}
	if commandline == "ssh -i /home/manuel.kiessling/.ssh/infmgmt.id_rsa root@vmhost1 virsh list --all | tail --lines=+3 | head --lines=1 | sed 's/ \\+/ /g' | cut -d' ' -f4-" {
		return "running", nil
	}
	if commandline == "ssh -i /home/manuel.kiessling/.ssh/infmgmt.id_rsa root@vmhost1 virsh dumpxml virtual1 | grep uuid | cut --bytes=9-44" {
		return "a0f39677-afda-f5bb-20b9-c5d8e3e06edf", nil
	}
	return "", nil
}

func setupRouter() *mux.Router {
	ce := new(MockVmguestRepositoryCommandExecutor)
	oh := interfaces.NewDefaultVmhostOperationsHandler(ce)

	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-integrationtestdb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt-backend:", log.Lmicroseconds))
//	dbMap.TraceOff()

	vglr := interfaces.NewVmguestLiveRepository(ce)
	vgcr := interfaces.NewVmguestCacheRepository(dbMap)
	vhr := interfaces.NewVmhostRepository(dbMap, vglr, vgcr)

	dbMap.DropTables()
	dbMap.CreateTables()
	dbMap.Exec("INSERT INTO vmhosts (Id, DnsName) VALUES (?, ?)", "1", "vmhost1")

	mi := new(domain.VmhostsInteractor)
	mi.VmhostRepository = vhr
	mi.VmhostOperationsHandler = oh
	rh := interfaces.NewRequestHandler(mi)

	r := interfaces.NewRouter(rh)

	return r
}

func TestGetVmhosts(t *testing.T) {
	updateCacheReq, err := http.NewRequest("POST", "http://example.com/cacheupdate", nil)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "http://example.com/vmhosts", nil)
	if err != nil {
		log.Fatal(err)
	}

	router := setupRouter()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, updateCacheReq)

	time.Sleep(1000)

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	expected := "{\"1\":{\"Id\":\"1\",\"DnsName\":\"vmhost1\",\"Vmguests\":{\"a0f39677-afda-f5bb-20b9-c5d8e3e06edf\":{\"Id\":\"a0f39677-afda-f5bb-20b9-c5d8e3e06edf\",\"Name\":\"virtual1\",\"State\":\"running\"}}}}"

	if expected != rec.Body.String() {
		t.Errorf("Expected response body %s, but got %s", expected, rec.Body.String())
	}
}

func TestGetVmhostsFromEmptyCache(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/vmhosts", nil)
	if err != nil {
		log.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rec, req)

	expected := "{\"1\":{\"Id\":\"1\",\"DnsName\":\"vmhost1\",\"Vmguests\":{}}}"

	if expected != rec.Body.String() {
		t.Errorf("Expected response body %s, but got %s", expected, rec.Body.String())
	}
}

func TestGetVmguests(t *testing.T) {
	updateCacheReq, err := http.NewRequest("POST", "http://example.com/cacheupdate", nil)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "http://example.com/vmhosts/1/vmguests", nil)
	if err != nil {
		log.Fatal(err)
	}

	router := setupRouter()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, updateCacheReq)

	time.Sleep(1000)
	
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	expected := "{\"a0f39677-afda-f5bb-20b9-c5d8e3e06edf\":{\"Id\":\"a0f39677-afda-f5bb-20b9-c5d8e3e06edf\",\"Name\":\"virtual1\",\"State\":\"running\"}}"

	if expected != rec.Body.String() {
		t.Errorf("Expected response body %s, but got %s", expected, rec.Body.String())
	}
}

func TestGetVmguestsFromEmptyCache(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/vmhosts/1/vmguests", nil)
	if err != nil {
		log.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rec, req)

	expected := "{}"

	if expected != rec.Body.String() {
		t.Errorf("Expected response body %s, but got %s", expected, rec.Body.String())
	}
}

func TestSetupWorkingForVirtualMachine(t *testing.T) {
	return
	req, err := http.NewRequest("POST", "http://example.com/vmhosts/1/vmguests", nil)
	if err != nil {
		log.Fatal(err)
	}

	rec := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Expected response code 200, but got %+v from request %+v", rec, req)
	}
}
