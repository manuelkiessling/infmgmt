package main

import (
	"database/sql"
	"github.com/ManuelKiessling/infmgmt-backend/infrastructure"
	"github.com/ManuelKiessling/infmgmt-backend/interfaces"
	"github.com/ManuelKiessling/infmgmt-backend/domain"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	ce := new(infrastructure.DefaultCommandExecutor)
	oh := interfaces.NewDefaultMachineOperationsHandler(ce)

	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-testdb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt-backend:", log.Lmicroseconds))
	mr := interfaces.NewMachineRepository(dbMap)

	mi := new(domain.MachinesInteractor)
	mi.MachineRepository = mr
	mi.MachineOperationsHandler = oh

	rh := interfaces.NewRequestHandler(mi)
	r := mux.NewRouter()

	// GET /machines returns a list of all machines
	r.HandleFunc("/machines", func(res http.ResponseWriter, req *http.Request) {
		rh.HandleMachinesRequest(res, req)
	})

	// POST /machines creates a new machine entity with the given data

	// POST /machines/{machineId}/setup triggers the procedure that creates a virtual machine
	r.HandleFunc("/machines/{machineId}/setup", func(res http.ResponseWriter, req *http.Request) {
		rh.HandleMachineSetupRequest(res, req)
	})

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
