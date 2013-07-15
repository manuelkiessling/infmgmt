package domain

/*

- implementiert business entities, zustände, rules: zB "eine VM vmhost darf nur einer VM-Host vmhost zugeordnet sein
- implementiert "setze zustand dieser vmhost auf 'online', getriggert durch use case

*/

import (
	_ "errors"
)

type Vmhost struct {
	Id          string
	DnsName     string
	TotalMemory int
	Vmguests    map[string]*Vmguest
}

type Vmguest struct {
	Id              string
	Name            string
	State           string
	AllocatedMemory int
}

func NewVmhost(id string, dnsName string, totalMemory int, vmguests map[string]*Vmguest) (newVmhost *Vmhost, err error) {
	vmhost := &Vmhost{id, dnsName, totalMemory, vmguests}
	return vmhost, nil
}

func NewVmguest(id string, name string, state string, allocatedMemory int) (newVmguest *Vmguest, err error) {
	vmguest := &Vmguest{id, name, state, allocatedMemory}
	return vmguest, nil
}
