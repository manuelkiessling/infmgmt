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

type procedure struct {
	Id              string
	Status          int
	commandExecutor CommandExecutor
	commands        []*Command
}

func (p *procedure) Add(command *Command) {
	p.commands = append(p.commands, command)
}

func (p *procedure) Start() chan int {
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

type DefaultVmhostOperationsHandler struct {
	commandExecutor CommandExecutor
	procedures      map[string]*procedure
}

func NewDefaultVmhostOperationsHandler(commandExecutor CommandExecutor) *DefaultVmhostOperationsHandler {
	oh := new(DefaultVmhostOperationsHandler)
	oh.commandExecutor = commandExecutor
	oh.procedures = make(map[string]*procedure)
	return oh
}

func (oh *DefaultVmhostOperationsHandler) InitializeProcedure() string {
	procedure := new(procedure)
	uuid, _ := simpleuuid.NewTime(time.Now())
	procedure.Id = uuid.String()
	procedure.commandExecutor = oh.commandExecutor
	oh.procedures[procedure.Id] = procedure
	return procedure.Id
}

func (oh *DefaultVmhostOperationsHandler) AddCommandCreateVirtualmachine(procedureId string, vmhostDnsName string, virtualmachineName string) error {
	command := new(Command)
	command.Name = "/usr/bin/touch"
	command.Arguments = append(command.Arguments, "/tmp/testfile-"+vmhostDnsName+"-"+virtualmachineName)
	oh.procedures[procedureId].Add(command)
	//("ssh root@" + vmhostDnsName + " cp /var/lib/libvirt/images/infmgmgt-base.raw /var/lib/libvirt/images/" + newImageName + ".raw")
	return nil
}

func (oh *DefaultVmhostOperationsHandler) ExecuteProcedure(procedureId string) chan int {
	return oh.procedures[procedureId].Start()
}

func (oh *DefaultVmhostOperationsHandler) GetProcedureStatus(procedureId string) int {
	return oh.procedures[procedureId].Status
}
