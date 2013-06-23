package domain

import (
	"testing"
)

func TestCreateVmguest(t *testing.T) {
	NewVmguest("12345", "virtual1", "virtual1.example.com", "shutdown")
}

func TestCreateVmhost(t *testing.T) {
	vmguests := make([]*Vmguest, 0)
	NewVmhost("12345", "kvmhost1", vmguests)
}

func TestAppendVmguestToVmhost(t *testing.T) {
	vmguests := make([]*Vmguest, 0)
	vmhost, _ := NewVmhost("12345", "kvmhost1", vmguests)
	vmguest, _ := NewVmguest("67890", "virtual1", "virtual1.example.com", "shutdown")
	vmhost.Vmguests = append(vmhost.Vmguests, vmguest)
}
