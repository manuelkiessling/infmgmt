package domain

import (
	"testing"
)

func TestAddVirtualMachine(t *testing.T) {
	pm, err := NewMachine("kvmhost3", P, nil)
	vm1, err := NewMachine("vm1", V, pm)
	if err != nil {
		t.Errorf("Could not create vm1 attached to kvmhost3, but it should be possible")
	}

	_, err = NewMachine("vm2", V, vm1)
	if err == nil {
		t.Errorf("Could create vm2 attached to vm1, but it should not be possible")
	}

	_, err = NewMachine("vm3", V, nil)
	if err == nil {
		t.Errorf("Could create vm2 attached to no other machine, but it should not be possible")
	}

	_, err = NewMachine("pm2", P, vm1)
	if err == nil {
		t.Errorf("Could create pm2 attached to virtual machine, but it should not be possible")
	}

	_, err = NewMachine("pm2", P, pm)
	if err == nil {
		t.Errorf("Could create pm2 attached to other physical machine, but it should not be possible")
	}
}

