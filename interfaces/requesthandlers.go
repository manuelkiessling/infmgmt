package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/manuelkiessling/infmgmt-backend/domain"
	"net/http"
)

type RequestHandler struct {
	vmhostsInteractor *domain.VmhostsInteractor
}

func NewRouter(requestHandler *RequestHandler) *mux.Router {
	router := mux.NewRouter()

	// GET /vmhosts returns a list of all vmhosts
	router.HandleFunc("/webservice/vmhosts", func(res http.ResponseWriter, req *http.Request) {
		requestHandler.HandleListVmhostsRequest(res, req)
	})

	// GET /vmhosts/{vmhostId}/vmguests shows list of virtual machines on vmhost
	router.HandleFunc("/webservice/vmhosts/{vmhostId}/vmguests", func(res http.ResponseWriter, req *http.Request) {
		requestHandler.HandleListVmguestsRequest(res, req)
	}).Methods("GET")

	// GET /vmhosts/{vmhostId}/vmguests/{vmguestName} shows details of named virtual machine
	//	router.HandleFunc("/vmhosts/{vmhostId}/vmguests/{vmguestName}", func(res http.ResponseWriter, req *http.Request) {
	//		requestHandler.HandleShowVmguestRequest(res, req)
	//	}).Methods("GET")

	// POST /vmhosts/{vmhostId}/vmguests triggers the procedure that creates a virtual machine
	router.HandleFunc("/webservice/vmhosts/{vmhostId}/vmguests", func(res http.ResponseWriter, req *http.Request) {
		requestHandler.HandleCreateVmguestRequest(res, req)
	}).Methods("POST")

	// POST /cacherefresh triggers update of the vmguest cache
	router.HandleFunc("/webservice/cacheupdate", func(res http.ResponseWriter, req *http.Request) {
		requestHandler.HandleUpdateCache(res, req)
	}).Methods("POST")

	// Application
	router.PathPrefix("/app").Handler(http.StripPrefix("/app", http.FileServer(http.Dir("/home/manuel.kiessling/Dropbox/Projects/infmgmt-frontend/htdocs"))))

	return router
}

func NewRequestHandler(vmhostsInteractor *domain.VmhostsInteractor) *RequestHandler {
	requestHandler := new(RequestHandler)
	requestHandler.vmhostsInteractor = vmhostsInteractor
	return requestHandler
}

func (rh *RequestHandler) HandleListVmhostsRequest(res http.ResponseWriter, req *http.Request) {
	list, _ := rh.vmhostsInteractor.GetList()
	jsonResponse, _ := json.Marshal(list)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, "%s", jsonResponse)
}

func (rh *RequestHandler) HandleListVmguestsRequest(res http.ResponseWriter, req *http.Request) {
	vmhostId := mux.Vars(req)["vmhostId"]
	list, _ := rh.vmhostsInteractor.GetVmguestsList(vmhostId)
	jsonResponse, _ := json.Marshal(list)
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

func (rh *RequestHandler) HandleUpdateCache(res http.ResponseWriter, req *http.Request) {
	err := rh.vmhostsInteractor.UpdateCache()
	if err == nil {
		fmt.Fprintf(res, "OK")
	} else {
		fmt.Fprintf(res, "%s", err)
	}
}
