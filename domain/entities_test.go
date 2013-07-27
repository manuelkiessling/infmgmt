package domain

import (
	"testing"
	"time"
)

func TestCreateVmguest(t *testing.T) {
	vmguest, err := NewVmguest("12345", "virtual1")
	if err != nil {
		t.Errorf("Could not instantiate new vmguest: %+v", err)
	}
	if vmguest.InfoUpdatedAt().IsZero() != true {
		t.Errorf("Vmguest was not instantiated with zero infoUpdatedAt time: %+v", vmguest.InfoUpdatedAt())
	}
	vmguest.SetState("shut down")
	vmguest.SetAllocatedMemory(1048576)
	now := time.Now()
	vmguest.SetInfoUpdatedAt(now)
	if vmguest.Id() != "12345" ||
		vmguest.Name() != "virtual1" ||
		vmguest.State() != "shut down" ||
		vmguest.AllocatedMemory() != 1048576 ||
		vmguest.InfoUpdatedAt().Equal(now) != true {
		t.Errorf("Could not create and read vmguest %+v", vmguest)
	}
}

func TestCreateVmhost(t *testing.T) {
	var vmguests map[string]*Vmguest
	vmhost, err := NewVmhost("12345", "kvmhost1")
	if err != nil {
		t.Errorf("Could not instantiate new vmhost: %+v", err)
	}
	vmhost.SetVmguests(vmguests)
	vmhost.SetTotalMemory(3944700)
	if vmhost.Id() != "12345" ||
		vmhost.DnsName() != "kvmhost1" ||
		vmhost.TotalMemory() != 3944700 {
		t.Errorf("Could not create and read vmhost %+v", vmhost)	
	}
}
