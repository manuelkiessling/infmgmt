package domain

import (
	"reflect"
	"testing"
)

type MockOperationsHandler struct {
	NumberOfCalls int
	Calls         map[string]int
}

func (oh *MockOperationsHandler) CopyBaseImage(kvmHostDnsName string, newImageName string) error {
	oh.Calls["CopyBaseImage"] = oh.NumberOfCalls
	oh.NumberOfCalls++
	return nil
}

func (oh *MockOperationsHandler) SetIpAddressInImage(kvmHostDnsName string, ipAddress string) error {
	oh.Calls["SetIpAddressInImage"] = oh.NumberOfCalls
	oh.NumberOfCalls++
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
	expectedCalls := make(map[string]int)
	expectedCalls["CopyBaseImage"] = 0
	expectedCalls["SetIpAddressInImage"] = 1
	oh := new(MockOperationsHandler)
	oh.Calls = make(map[string]int)
	pm, _ := NewMachine("kvmhost3", P, nil, oh)
	vm1, _ := NewMachine("vm1", V, pm, oh)
	vm1.Setup()
	if !reflect.DeepEqual(expectedCalls, oh.Calls) {
		t.Errorf("Setup() did not execute OperationsHandle commands in the right order, calls were: %+v", oh.Calls)
	}
}
