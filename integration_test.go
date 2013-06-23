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

	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-testdb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt-backend:", log.Lmicroseconds))
	dbMap.TraceOff()

	mr := interfaces.NewVmhostRepository(dbMap)

	dbMap.DropTables()
	dbMap.CreateTables()
	dbMap.Exec("INSERT INTO vmhosts (Id, DnsName) VALUES (?, ?)", "1", "kvmhost1")
	dbMap.Exec("INSERT INTO vmhosts (Id, DnsName) VALUES (?, ?)", "2", "kvmhost2")

	mi := new(domain.VmhostsInteractor)
	mi.VmhostRepository = mr
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

	expected := "{\"1\":{\"Id\":\"1\",\"DnsName\":\"kvmhost1\"},\"2\":{\"Id\":\"2\",\"DnsName\":\"kvmhost2\"}}"

	if expected != rec.Body.String() {
		t.Errorf("Expected response body %s, but got %s", expected, rec.Body.String())
	}
}

func TestSetupWorkingForVirtualMachine(t *testing.T) {
	req, err := http.NewRequest("POST", "http://example.com/vmhosts/1/vmguests/foo", nil)
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
