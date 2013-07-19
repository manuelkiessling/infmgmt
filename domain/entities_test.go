package domain

import (
	"testing"
)

func TestCreateVmguest(t *testing.T) {
	vmguest := NewVmguest("12345", "virtual1", "shutdown")
	vmguest.AllocatedMemory = 1048576
}

func TestCreateVmhost(t *testing.T) {
	var vmguests map[string]*Vmguest
	vmhost := NewVmhost("12345", "kvmhost1", vmguests)
	vmhost.TotalMemory = 3944700
	vmhost.DiskFree = 69492960
	vmhost.DiskSize	= 114076280
}
