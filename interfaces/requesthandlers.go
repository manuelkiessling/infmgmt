package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/ManuelKiessling/infmgmt-backend/domain"
	"github.com/gorilla/mux"
	"net/http"
)

type RequestHandler struct {
	machineRepository        domain.MachineRepository
	machineOperationsHandler domain.MachineOperationsHandler
}

func NewRequestHandler(machineRepository domain.MachineRepository, machineOperationsHandler domain.MachineOperationsHandler) *RequestHandler {
	requestHandler := new(RequestHandler)
	requestHandler.machineRepository = machineRepository
	requestHandler.machineOperationsHandler = machineOperationsHandler
	return requestHandler
}

func (rh *RequestHandler) HandleMachinesRequest(w http.ResponseWriter, r *http.Request) {
	interactor := new(domain.MachinesInteractor)
	interactor.MachineRepository = rh.machineRepository
	interactor.MachineOperationsHandler = rh.machineOperationsHandler

	allMachines, _ := interactor.ShowOverviewList()

	jsonResponse, _ := json.Marshal(allMachines)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", jsonResponse)
}

func (rh *RequestHandler) HandleMachineSetupRequest(w http.ResponseWriter, r *http.Request) {
	interactor := new(domain.MachinesInteractor)
	interactor.MachineRepository = rh.machineRepository
	interactor.MachineOperationsHandler = rh.machineOperationsHandler
	machineId := mux.Vars(r)["machineId"]
	interactor.SetupMachine(machineId)
//	fmt.Fprintf(w, "%+v", mux.Vars(r))
}
