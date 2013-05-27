package usecases

import (
	"testing"
	"github.com/ManuelKiessling/infmgmt-backend/domain"
)

type MockRepository struct {
}

func (repo *MockRepository) Store(machine *domain.Machine) error {
	return nil
}

func (repo *MockRepository) FindById(id string) *domain.Machine {
	return new(domain.Machine)
}

func (repo *MockRepository) GetAll() map[string]*domain.Machine {
	machines := make(map[string] *domain.Machine)
	machines["101"] = &domain.Machine{"101", "Mocked machine #1", domain.P, nil}
	machines["102"] = &domain.Machine{"102", "Mocked machine #2", domain.P, nil}
	return machines
}


func TestList(t *testing.T) {
	interactor := new(MachineOverviewInteractor)
	interactor.MachineRepository = new(MockRepository)
	overviewListEntries, _ := interactor.List()
	if len(overviewListEntries) != 2 || overviewListEntries["102"].DnsName != "Mocked machine #2" {
		t.Errorf("%+v", overviewListEntries)
	}
}
