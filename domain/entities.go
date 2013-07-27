package domain

/*

- implementiert business entities, zustände, rules: zB "eine VM vmhost darf nur einer VM-Host vmhost zugeordnet sein
- implementiert "setze zustand dieser vmhost auf 'online', getriggert durch use case

*/

import (
	_ "errors"
	"time"
)

type Vmhost struct {
	id          string
	dnsName     string
	totalMemory int
	vmguests    map[string]*Vmguest
}

type Vmguest struct {
	id              string
	name            string
	state           string
	allocatedMemory int
	infoUpdatedAt		time.Time
}

func NewVmhost(id string, dnsName string) (newVmhost *Vmhost, err error) {
	vmhost := &Vmhost{id, dnsName, 0, nil}
	return vmhost, nil
}

func (vh *Vmhost) Id() string {
	return vh.id
}

func (vh *Vmhost) DnsName() string {
	return vh.dnsName
}

func (vh *Vmhost) SetVmguests(vmguests map[string]*Vmguest) {
	vh.vmguests = vmguests
}

func (vh *Vmhost) Vmguests() map[string]*Vmguest {
	return vh.vmguests
}

func (vh *Vmhost) SetTotalMemory(totalMemory int) {
	vh.totalMemory = totalMemory
}

func (vh *Vmhost) TotalMemory() int {
	return vh.totalMemory
}


func NewVmguest(id string, name string) (newVmguest *Vmguest, err error) {
	zeroTime, _ := time.Parse("2013-07-27", "0000-01-01")
	vmguest := &Vmguest{id, name, "unknown", 0, zeroTime}
	return vmguest, nil
}

func (vg *Vmguest) Id() string {
	return vg.id
}

func (vg *Vmguest) Name() string {
	return vg.name
}

func (vg *Vmguest) SetState(state string) {
	vg.state = state
}

func (vg *Vmguest) State() string {
	return vg.state
}

func (vg *Vmguest) SetAllocatedMemory(allocatedMemory int) {
	vg.allocatedMemory = allocatedMemory
}

func (vg *Vmguest) AllocatedMemory() int {
	return vg.allocatedMemory
}

func (vg *Vmguest) SetInfoUpdatedAt(infoUpdatedAt time.Time) {
	vg.infoUpdatedAt = infoUpdatedAt
}

func (vg *Vmguest) InfoUpdatedAt() time.Time {
	return vg.infoUpdatedAt
}
