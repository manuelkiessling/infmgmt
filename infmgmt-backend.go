package main

import (
	"database/sql"
	"github.com/manuelkiessling/infmgmt-backend/infrastructure"
	"github.com/manuelkiessling/infmgmt-backend/interfaces"
	"github.com/manuelkiessling/infmgmt-backend/domain"
	"github.com/coopernurse/gorp"
	"log"
	_ "net/http"
	"os"
	"fmt"
)

func main() {
	ce := new(infrastructure.DefaultCommandExecutor)
//	oh := interfaces.NewDefaultVmhostOperationsHandler(ce)

	vgr := interfaces.NewVmguestRepository(ce)

	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-testdb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt-backend:", log.Lmicroseconds))
	vhr := interfaces.NewVmhostRepository(dbMap, vgr)
	
	vmhost, _ := domain.NewVmhost("1", "localhost", nil)
	vhr.Store(vmhost)
	vmhost, _ = vhr.FindById("1")
	fmt.Printf("Guests: %+v", vmhost.Vmguests[0])

//	mi := new(domain.VmhostsInteractor)
//	mi.VmhostRepository = vhr
//	mi.VmhostOperationsHandler = oh
//
//	rh := interfaces.NewRequestHandler(mi)
//
//	r := interfaces.NewRouter(rh)
//
//	http.Handle("/", r)
//	http.ListenAndServe(":8080", nil)
}
