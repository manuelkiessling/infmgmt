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

func NewRouter(requestHandler *RequestHandler) *mux.Router {
	router := mux.NewRouter()

	// GET /machines returns a list of all machines
	router.HandleFunc("/machines", func(res http.ResponseWriter, req *http.Request) {
		requestHandler.HandleMachinesRequest(res, req)
	})

	// POST /machines creates a new machine entity with the given data

	// POST /machines/{machineId}/setup triggers the procedure that creates a virtual machine
	router.HandleFunc("/machines/{machineId}/setup", func(res http.ResponseWriter, req *http.Request) {
		requestHandler.HandleMachineSetupRequest(res, req)
	})

	return router
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
	if err != nil {
		http.Error(res, "Can't setup a non-virtual machine", 500)
	} else {
		fmt.Fprintf(res, "OK")
	}
}
