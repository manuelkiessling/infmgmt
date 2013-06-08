package domain

import (
	"testing"
)

type MockOperationsHandler struct {
	SetupWasCalled bool
	RebootWasCalled bool
}

func (oh *MockOperationsHandler) Setup(machine *Machine) error {
	oh.SetupWasCalled = true
	return nil
}

func (oh *MockOperationsHandler) Reboot(machine *Machine) error {
	oh.RebootWasCalled = true
	return nil
}

func TestAddVirtualMachine(t *testing.T) {
	pm, err := NewMachine("kvmhost3", P, nil, nil)
	vm1, err := NewMachine("vm1", V, pm, nil)
	if err != nil {
		t.Errorf("Could not create vm1 attached to kvmhost3, but it should be possible")
	}

	_, err = NewMachine("vm2", V, vm1, nil)
	if err == nil {
		t.Errorf("Could create vm2 attached to vm1, but it should not be possible")
	}

	_, err = NewMachine("vm3", V, nil, nil)
	if err == nil {
		t.Errorf("Could create vm2 attached to no other machine, but it should not be possible")
	}

	_, err = NewMachine("pm2", P, vm1, nil)
	if err == nil {
		t.Errorf("Could create pm2 attached to virtual machine, but it should not be possible")
	}

	_, err = NewMachine("pm2", P, pm, nil)
	if err == nil {
		t.Errorf("Could create pm2 attached to other physical machine, but it should not be possible")
	}
}

func TestSetupVirtualMachine(t *testing.T) {
	oh := new(MockOperationsHandler)
	pm, _ := NewMachine("kvmhost3", P, nil, oh)
	vm1, _ := NewMachine("vm1", V, pm, oh)
	vm1.Setup()
	if !oh.SetupWasCalled {
		t.Errorf("Setup() of machines's OperationsHandler was not called.")
	}
}

func TestRebootVirtualMachine(t *testing.T) {
	oh := new(MockOperationsHandler)
	pm, _ := NewMachine("kvmhost3", P, nil, oh)
	vm1, _ := NewMachine("vm1", V, pm, oh)
	vm1.Reboot()
	if !oh.RebootWasCalled {
		t.Errorf("Reboot() of machines's OperationsHandler was not called.")
	}
}
