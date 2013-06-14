package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/ManuelKiessling/infmgmt-backend/domain"
	"github.com/gorilla/mux"
	"net/http"
)

type RequestHandler struct {
	machinesInteractor *domain.MachinesInteractor
}

func NewRequestHandler(machinesInteractor *domain.MachinesInteractor) *RequestHandler {
	requestHandler := new(RequestHandler)
	requestHandler.machinesInteractor = machinesInteractor
	return requestHandler
}

func (rh *RequestHandler) HandleMachinesRequest(res http.ResponseWriter, req *http.Request) {
	allMachines, _ := rh.machinesInteractor.ShowOverviewList()

	jsonResponse, _ := json.Marshal(allMachines)

	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, "%s", jsonResponse)
}

func (rh *RequestHandler) HandleMachineSetupRequest(res http.ResponseWriter, req *http.Request) {
	machineId := mux.Vars(req)["machineId"]
	_, err := rh.machinesInteractor.SetupMachine(machineId)
	if (err != nil) {
		http.Error(res, "Can't setup a non-virtual machine", 500)
	} else {
		fmt.Fprintf(res, "OK")
	}
}
