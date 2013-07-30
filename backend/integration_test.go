package main

import (
	"database/sql"
	_ "fmt"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
	"github.com/manuelkiessling/infmgmt/backend/domain"
	"github.com/manuelkiessling/infmgmt/backend/infrastructure"
	"github.com/manuelkiessling/infmgmt/backend/interfaces"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"strconv"
)

func setupRouter() *mux.Router {
	ce := new(infrastructure.MockCommandExecutor)
	oh := interfaces.NewDefaultVmhostOperationsHandler(ce)

	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-integrationtestdb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt/backend:", log.Lmicroseconds))
	dbMap.TraceOff()

	vglr := interfaces.NewVmguestLiveRepository(ce)
	vgcr := interfaces.NewVmguestCacheRepository(dbMap)
	vhr := interfaces.NewVmhostRepository(dbMap, vglr, vgcr)

	dbMap.DropTables()
	dbMap.CreateTables()
	dbMap.Exec("INSERT INTO vmhosts (Id, DnsName, TotalMemory) VALUES (?, ?, ?)", "1", "vmhost1", 32918292)

	mi := new(domain.VmhostsInteractor)
	mi.VmhostRepository = vhr
	mi.VmhostOperationsHandler = oh
	rh := interfaces.NewRequestHandler(mi)

	r := interfaces.NewRouter(rh)

	return r
}

func TestGetVmhosts(t *testing.T) {
	updateCacheReq, err := http.NewRequest("POST", "http://example.com/webservice/cacheupdate", nil)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "http://example.com/webservice/vmhosts", nil)
	if err != nil {
		log.Fatal(err)
	}

	router := setupRouter()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, updateCacheReq)

	expectedTimeString := strconv.FormatInt(time.Now().Unix(), 10)

	time.Sleep(1000 * time.Millisecond)

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	expected := "{\"1\":{\"Id\":\"1\",\"DnsName\":\"vmhost1\",\"TotalMemory\":32918292,\"Vmguests\":{\"a0f39677-afda-f5bb-20b9-c5d8e3e06edf\":{\"Id\":\"a0f39677-afda-f5bb-20b9-c5d8e3e06edf\",\"Name\":\"virtual1\",\"State\":\"running\",\"AllocatedMemory\":1048576,\"InfoUpdatedAt\":"+expectedTimeString+"}}}}"

	if expected != rec.Body.String() {
		t.Errorf("Expected response body %s, but got %s", expected, rec.Body.String())
	}
}

func TestGetVmhostsFromEmptyCache(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/webservice/vmhosts", nil)
	if err != nil {
		log.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rec, req)

	expected := "{\"1\":{\"Id\":\"1\",\"DnsName\":\"vmhost1\",\"TotalMemory\":32918292,\"Vmguests\":{}}}"

	if expected != rec.Body.String() {
		t.Errorf("Expected response body %s, but got %s", expected, rec.Body.String())
	}
}

func TestGetVmguests(t *testing.T) {
	updateCacheReq, err := http.NewRequest("POST", "http://example.com/webservice/cacheupdate", nil)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "http://example.com/webservice/vmhosts/1/vmguests", nil)
	if err != nil {
		log.Fatal(err)
	}

	router := setupRouter()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, updateCacheReq)

	expectedTimeString := strconv.FormatInt(time.Now().Unix(), 10)

	time.Sleep(2000)

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	expected := "{\"a0f39677-afda-f5bb-20b9-c5d8e3e06edf\":{\"Id\":\"a0f39677-afda-f5bb-20b9-c5d8e3e06edf\",\"Name\":\"virtual1\",\"State\":\"running\",\"AllocatedMemory\":1048576,\"InfoUpdatedAt\":"+expectedTimeString+"}}"

	if expected != rec.Body.String() {
		t.Errorf("Expected response body %s, but got %s", expected, rec.Body.String())
	}
}

func TestGetVmguestsFromEmptyCache(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/webservice/vmhosts/1/vmguests", nil)
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
	req, err := http.NewRequest("POST", "http://example.com/webservice/vmhosts/1/vmguests", nil)
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
