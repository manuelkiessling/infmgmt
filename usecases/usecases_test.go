package usecases

import (
	"testing"
	"github.com/ManuelKiessling/infmgmt-backend/domain"
)

type MockRepository struct {
}

func (repo *MockRepository) Store(machine *domain.Machine) {}

func (repo *MockRepository) FindById(id int) *domain.Machine {
	return new(domain.Machine)
}

func (repo *MockRepository) GetAll() []*domain.Machine {
	machines := make([]*domain.Machine, 1)
	machines[0] = &domain.Machine{101, "Mocked machine #1"}
	machines[1] = &domain.Machine{102, "Mocked machine #2"}
	return machines
}


func TestList(t *testing.T) {
	interactor := new(MachineOverviewInteractor)
	interactor.MachineRepository = new(MockRepository)
	overviewListEntries, _ := interactor.List()
	if len(overviewListEntries) != 2 || overviewListEntries[1].DnsName != "Mocked machine #2" {
		t.Errorf("%+v", overviewListEntries)
	}
}
