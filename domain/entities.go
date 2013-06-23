package domain

/*

- implementiert business entities, zustände, rules: zB "eine VM vmhost darf nur einer VM-Host vmhost zugeordnet sein
- implementiert "setze zustand dieser vmhost auf 'online', getriggert durch use case

*/

import (
	_ "errors"
	"github.com/streadway/simpleuuid"
	"time"
)

const (
	P = 0
	V = 1
)

type Vmhost struct {
	Id          string
	DnsName     string
}

func NewVmhost(dnsName string) (newVmhost *Vmhost, err error) {
	uuid, _ := simpleuuid.NewTime(time.Now())
	id := uuid.String()

	vmhost := &Vmhost{id, dnsName}
	return vmhost, nil
}
