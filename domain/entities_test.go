package domain

import (
	"testing"
)

func TestCreateVmguest(t *testing.T) {
	NewVmguest("12345", "virtual1", "shutdown")
}

func TestCreateVmhost(t *testing.T) {
	var vmguests map[string]*Vmguest
	NewVmhost("12345", "kvmhost1", vmguests)
}
