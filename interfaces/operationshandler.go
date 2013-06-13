package interfaces

import ()

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
	oh.commandExecutor.Run("touch testfile-" + vmhostDnsName + "-" + newImageName)
	//oh.commandExecutor.Run("ssh root@" + vmhostDnsName + " cp /var/lib/libvirt/images/infmgmgt-base.raw /var/lib/libvirt/images/" + newImageName + ".raw")
	return nil
}
