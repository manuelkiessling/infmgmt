package domain

import (
	"testing"
)

func TestCreateVmguest(t *testing.T) {
	NewVmguest("12345", "virtual1", "shutdown", 1048576)
}

func TestCreateVmhost(t *testing.T) {
	var vmguests map[string]*Vmguest
	NewVmhost("12345", "kvmhost1", 32918292, vmguests)
}
