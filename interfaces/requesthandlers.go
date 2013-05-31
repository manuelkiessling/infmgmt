package interfaces

import (
	"github.com/ManuelKiessling/infmgmt-backend/usecases"
	"fmt"
	"database/sql"
	"github.com/coopernurse/gorp"
	"net/http"
	"log"
	"os"
)

func MachinesRequesthandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-testdb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt-backend:", log.Lmicroseconds)) 
	repo := NewMachineRepository(dbMap)
	interactor := &usecases.MachineOverviewInteractor{repo}
	allMachines, _ := interactor.List()
	fmt.Fprintf(w, "Machines: %+v", allMachines)
}

