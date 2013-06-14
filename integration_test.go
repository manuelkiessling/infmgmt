package main

import (
	"database/sql"
	_ "fmt"
	"github.com/ManuelKiessling/infmgmt-backend/domain"
	"github.com/ManuelKiessling/infmgmt-backend/infrastructure"
	"github.com/ManuelKiessling/infmgmt-backend/interfaces"
	"github.com/coopernurse/gorp"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupRequestHandler() *interfaces.RequestHandler {
	ce := new(infrastructure.DefaultCommandExecutor)
	oh := interfaces.NewDefaultMachineOperationsHandler(ce)

	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-testdb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt-backend:", log.Lmicroseconds))
	dbMap.TraceOff()

	mr := interfaces.NewMachineRepository(dbMap)

	dbMap.DropTables()
	dbMap.CreateTables()
	dbMap.Exec("INSERT INTO machines (Id, DnsName, MachineType, VmhostId) VALUES (?, ?, ?, ?)", "1", "kvmhost1", 0, "")
	dbMap.Exec("INSERT INTO machines (Id, DnsName, MachineType, VmhostId) VALUES (?, ?, ?, ?)", "2", "virtual1", 1, "1")

	mi := new(domain.MachinesInteractor)
	mi.MachineRepository = mr
	mi.MachineOperationsHandler = oh

	return interfaces.NewRequestHandler(mi)
}

func TestGetMachines(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/machines", nil)
	if err != nil {
		log.Fatal(err)
	}

	rec := httptest.NewRecorder()
	requestHandler := setupRequestHandler()
	requestHandler.HandleMachinesRequest(rec, req)

	expected := "{\"1\":{\"Id\":\"1\",\"DnsName\":\"kvmhost1\"},\"2\":{\"Id\":\"2\",\"DnsName\":\"virtual1\"}}"

	if expected != rec.Body.String() {
		t.Errorf("Expected response body %s, but got %s", expected, rec.Body.String())
	}
}

func TestSetupFailingForPhysicalMachine(t *testing.T) {
	req, err := http.NewRequest("POST", "http://example.com/machine/1/setup", nil)
	if err != nil {
		log.Fatal(err)
	}

	rec := httptest.NewRecorder()

	requestHandler := setupRequestHandler()

	requestHandler.HandleMachineSetupRequest(rec, req)

	if rec.Code != 500 {
		t.Errorf("Expected response code 500, but got %d", rec.Code)
	}
}
