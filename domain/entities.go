package domain

/*

- implementiert business entities, zustände, rules: zB "eine VM vmhost darf nur einer VM-Host vmhost zugeordnet sein
- implementiert "setze zustand dieser vmhost auf 'online', getriggert durch use case

*/

import (
	_ "errors"
)

type Vmhost struct {
	Id      string
	DnsName string
	Vmguests []*Vmguest
}

type Vmguest struct {
	Id      string
	Name    string
	State   string
}

func NewVmhost(id string, dnsName string, vmguests []*Vmguest) (newVmhost *Vmhost, err error) {
	vmhost := &Vmhost{id, dnsName, vmguests}
	return vmhost, nil
}

func NewVmguest(id string, name string, state string) (newVmguest *Vmguest, err error) {
	vmguest := &Vmguest{id, name, state}
	return vmguest, nil
}
