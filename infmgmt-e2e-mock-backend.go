package main

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/manuelkiessling/infmgmt-backend/domain"
	"github.com/manuelkiessling/infmgmt-backend/infrastructure"
	"github.com/manuelkiessling/infmgmt-backend/interfaces"
	"log"
	"net/http"
	"os"
)

func main() {
	ce := new(infrastructure.MockCommandExecutor)
	oh := interfaces.NewDefaultVmhostOperationsHandler(ce)

	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-productiondb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt-backend:", log.Lmicroseconds))

	vglr := interfaces.NewVmguestLiveRepository(ce)
	vgcr := interfaces.NewVmguestCacheRepository(dbMap)
	vhr := interfaces.NewVmhostRepository(dbMap, vglr, vgcr)

	dbMap.DropTables()
	dbMap.CreateTables()

	vmhost, _ := domain.NewVmhost("1", "vmhost1")
	vmhost.SetTotalMemory(32918292)
	vhr.Store(vmhost)
	/*
	vmhost, _ := domain.NewVmhost("1", "kvmhost1.meinauto.local", 32918292, nil)
	vhr.Store(vmhost)
	vmhost, _ = domain.NewVmhost("2", "kvmhost2.meinauto.local", 32918280, nil)
	vhr.Store(vmhost)
	vmhost, _ = domain.NewVmhost("3", "kvmhost3.meinauto.local", 24683540, nil)
	vhr.Store(vmhost)
  */
	vi := new(domain.VmhostsInteractor)
	vi.VmhostRepository = vhr
	vi.VmhostOperationsHandler = oh

	vi.UpdateCache()

	rh := interfaces.NewRequestHandler(vi)

	r := interfaces.NewRouter(rh)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
