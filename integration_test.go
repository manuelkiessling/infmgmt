package main

import (
	"database/sql"
	_ "fmt"
	"github.com/manuelkiessling/infmgmt-backend/domain"
	"github.com/manuelkiessling/infmgmt-backend/infrastructure"
	"github.com/manuelkiessling/infmgmt-backend/interfaces"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupRouter() *mux.Router {
	ce := new(infrastructure.DefaultCommandExecutor)
	oh := interfaces.NewDefaultVmhostOperationsHandler(ce)

	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-integrationtestdb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt-backend:", log.Lmicroseconds))
	dbMap.TraceOff()

	vgr := interfaces.NewVmguestRepository(dbMap, ce)
	vhr := interfaces.NewVmhostRepository(dbMap, vgr)

	dbMap.DropTables()
	dbMap.CreateTables()
	dbMap.Exec("INSERT INTO vmhosts (Id, DnsName) VALUES (?, ?)", "1", "localhost")

	mi := new(domain.VmhostsInteractor)
	mi.VmhostRepository = vhr
	mi.VmhostOperationsHandler = oh
	rh := interfaces.NewRequestHandler(mi)

	r := interfaces.NewRouter(rh)

	return r
}

func TestGetVmhosts(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/vmhosts", nil)
	if err != nil {
		log.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rec, req)

	expected := "{\"1\":{\"Id\":\"1\",\"DnsName\":\"localhost\",\"Vmguests\":{\"a0f39677-afda-f5bb-20b9-c5d8e3e06edf\":{\"Id\":\"a0f39677-afda-f5bb-20b9-c5d8e3e06edf\",\"Name\":\"wordpress\",\"State\":\"shut off\"}}}}"

	if expected != rec.Body.String() {
		t.Errorf("Expected response body %s, but got %s", expected, rec.Body.String())
	}
}

func TestGetVmguests(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/vmhosts/1/vmguests", nil)
	if err != nil {
		log.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rec, req)

	expected := "{\"a0f39677-afda-f5bb-20b9-c5d8e3e06edf\":{\"Id\":\"a0f39677-afda-f5bb-20b9-c5d8e3e06edf\",\"Name\":\"wordpress\",\"State\":\"shut off\"}}"

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
