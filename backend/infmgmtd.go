package main

import (
	"database/sql"
	_ "fmt"
	"github.com/coopernurse/gorp"
	"github.com/manuelkiessling/infmgmt/backend/domain"
	"github.com/manuelkiessling/infmgmt/backend/infrastructure"
	"github.com/manuelkiessling/infmgmt/backend/interfaces"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ce := new(infrastructure.DefaultCommandExecutor)
	oh := interfaces.NewDefaultVmhostOperationsHandler(ce)

	db, _ := sql.Open("sqlite3", "/tmp/infmgmt-productiondb.sqlite")
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "infmgmt/backend:", log.Lmicroseconds))

	vglr := interfaces.NewVmguestLiveRepository(ce)
	vgcr := interfaces.NewVmguestCacheRepository(dbMap)
	vhr := interfaces.NewVmhostRepository(dbMap, vglr, vgcr)

	dbMap.DropTables()
	dbMap.CreateTables()

	vmhost, _ := domain.NewVmhost("1", "kvmhost1.meinauto.local")
	vmhost.SetTotalMemory(32918292)
	vhr.Store(vmhost)
	vmhost, _ = domain.NewVmhost("2", "kvmhost2.meinauto.local")
	vmhost.SetTotalMemory(32918280)
	vhr.Store(vmhost)
	vmhost, _ = domain.NewVmhost("3", "kvmhost3.meinauto.local")
	vmhost.SetTotalMemory(24683540)
	vhr.Store(vmhost)

	mi := new(domain.VmhostsInteractor)
	mi.VmhostRepository = vhr
	mi.VmhostOperationsHandler = oh

	rh := interfaces.NewRequestHandler(mi)

	r := interfaces.NewRouter(rh)

	go func() {
		for true {
			mi.UpdateCache()
			time.Sleep(1 * time.Minute)
		}
	}()

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
