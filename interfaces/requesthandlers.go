package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/manuelkiessling/infmgmt-backend/domain"
	"github.com/gorilla/mux"
	"net/http"
)

type RequestHandler struct {
	vmhostsInteractor *domain.VmhostsInteractor
}

func NewRouter(requestHandler *RequestHandler) *mux.Router {
	router := mux.NewRouter()

	// GET /vmhosts returns a list of all vmhosts
	router.HandleFunc("/vmhosts", func(res http.ResponseWriter, req *http.Request) {
		requestHandler.HandleVmhostsRequest(res, req)
	})

	// POST /vmhosts/{vmhostId}/create/vmguestName triggers the procedure that creates a virtual machine
	router.HandleFunc("/vmhosts/{vmhostId}/vmguests/{vmguestName}", func(res http.ResponseWriter, req *http.Request) {
		requestHandler.HandleCreateVmguestRequest(res, req)
	})

	return router
}

func NewRequestHandler(vmhostsInteractor *domain.VmhostsInteractor) *RequestHandler {
	requestHandler := new(RequestHandler)
	requestHandler.vmhostsInteractor = vmhostsInteractor
	return requestHandler
}

func (rh *RequestHandler) HandleVmhostsRequest(res http.ResponseWriter, req *http.Request) {
	overviewList, _ := rh.vmhostsInteractor.GetOverviewList()

	jsonResponse, _ := json.Marshal(overviewList)

	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, "%s", jsonResponse)
}

func (rh *RequestHandler) HandleCreateVmguestRequest(res http.ResponseWriter, req *http.Request) {
	vmhostId := mux.Vars(req)["vmhostId"]
	vmguestName := mux.Vars(req)["vmguestName"]
	_, err := rh.vmhostsInteractor.CreateVmguest(vmhostId, vmguestName)
	if err != nil {
		http.Error(res, "Could not create virtual guest '"+vmguestName+"' on vmhost "+vmhostId, 500)
	} else {
		fmt.Fprintf(res, "OK")
	}
}
