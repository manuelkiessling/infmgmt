package interfaces

import (
	_ "fmt"
	"github.com/streadway/simpleuuid"
	"time"
)

type CommandExecutor interface {
	Run(command string, arguments ...string) (output string, err error)
}

type Command struct {
	Name      string
	Arguments []string
}

type Procedure struct {
	Id              string
	Status          int
	commandExecutor CommandExecutor
	commands        []*Command
}

func (p *Procedure) Add(command *Command) {
	p.commands = append(p.commands, command)
}

func (p *Procedure) Start() chan int {
	c := make(chan int)
	go func() {
		for _, command := range p.commands {
			p.commandExecutor.Run(command.Name, command.Arguments...)
		}
		p.Status = 1
		c <- 1
	}()
	return c
}

type DefaultMachineOperationsHandler struct {
	commandExecutor CommandExecutor
	procedures      map[string]*Procedure
}

func NewDefaultMachineOperationsHandler(commandExecutor CommandExecutor) *DefaultMachineOperationsHandler {
	oh := new(DefaultMachineOperationsHandler)
	oh.commandExecutor = commandExecutor
	oh.procedures = make(map[string]*Procedure)
	return oh
}

func (oh *DefaultMachineOperationsHandler) NewProcedure() *Procedure {
	procedure := new(Procedure)
	uuid, _ := simpleuuid.NewTime(time.Now())
	procedure.Id = uuid.String()
	procedure.commandExecutor = oh.commandExecutor
	oh.procedures[procedure.Id] = procedure
	return procedure
}

func (oh *DefaultMachineOperationsHandler) CommandCreateVirtualMachine(vmhostDnsName string, machineName string) (*Command, error) {
	command := new(Command)
	command.Name = "/usr/bin/touch"
	command.Arguments = append(command.Arguments, "/tmp/testfile-"+vmhostDnsName+"-"+machineName)
	//("ssh root@" + vmhostDnsName + " cp /var/lib/libvirt/images/infmgmgt-base.raw /var/lib/libvirt/images/" + newImageName + ".raw")
	return command, nil
}
