package main

import (
	"database/sql"
	"github.com/ManuelKiessling/infmgmt-backend/infrastructure"
	"github.com/ManuelKiessling/infmgmt-backend/interfaces"
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
	repo := interfaces.NewMachineRepository(dbMap)

	rh := interfaces.NewRequestHandler(repo, oh)
	r := mux.NewRouter()

	r.HandleFunc("/machines", func(res http.ResponseWriter, req *http.Request) {
		rh.HandleMachinesRequest(res, req)
	})

	r.HandleFunc("/machines/{machineId}/setup", func(res http.ResponseWriter, req *http.Request) {
		rh.HandleMachineSetupRequest(res, req)
	})

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
