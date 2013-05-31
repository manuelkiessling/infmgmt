package interfaces

import (
	"github.com/ManuelKiessling/infmgmt-backend/usecases"
	"github.com/ManuelKiessling/infmgmt-backend/domain"
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

	machine, _ := domain.NewMachine("blah", domain.P, nil)
	repo.Store(machine)

	interactor := new(usecases.MachineOverviewInteractor)
	interactor.MachineRepository = repo
	allMachines, _ := interactor.List()
	fmt.Fprintf(w, "Machines: %+v", allMachines)
}

