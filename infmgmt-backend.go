package main

import (
	"database/sql"
	"github.com/manuelkiessling/infmgmt-backend/infrastructure"
	"github.com/manuelkiessling/infmgmt-backend/interfaces"
	"github.com/manuelkiessling/infmgmt-backend/domain"
	"github.com/coopernurse/gorp"
	"log"
	"net/http"
	"os"
)

func main() {
	ce := new(infrastructure.DefaultCommandExecutor)
	oh := interfaces.NewDefaultVmhostOperationsHandler(ce)

	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-productiondb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt-backend:", log.Lmicroseconds))

	vglr := interfaces.NewVmguestLiveRepository(ce)
	vgcr := interfaces.NewVmguestCacheRepository(dbMap)
	vhr := interfaces.NewVmhostRepository(dbMap, vglr, vgcr)
	
	dbMap.DropTables()
	dbMap.CreateTables()

	vmhost, _ := domain.NewVmhost("1", "localhost", nil)
	vhr.Store(vmhost)

	mi := new(domain.VmhostsInteractor)
	mi.VmhostRepository = vhr
	mi.VmhostOperationsHandler = oh

	rh := interfaces.NewRequestHandler(mi)

	r := interfaces.NewRouter(rh)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
