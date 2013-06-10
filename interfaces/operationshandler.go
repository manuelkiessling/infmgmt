package interfaces

import (
)

type CommandExecutor interface {
	Run(commandline string) (output string, err error)
}

type DefaultMachineOperationsHandler struct {
	commandExecutor CommandExecutor
}

func NewDefaultMachineOperationsHandler(commandExecutor CommandExecutor) *DefaultMachineOperationsHandler {
	oh := new(DefaultMachineOperationsHandler)
	oh.commandExecutor = commandExecutor
	return oh
}

func (oh *DefaultMachineOperationsHandler) CreateGuestImageFromBaseImage(vmhostDnsName string, newImageName string) error {
	return nil	
}
