package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/manuelkiessling/infmgmt-backend/domain"
	"net/http"
)

type vmhostResource struct {
	Id          string
	DnsName     string
	TotalMemory int
	Vmguests    map[string]*vmguestResource
}

type vmguestResource struct {
	Id              string
	Name            string
	State           string
	AllocatedMemory int
	InfoUpdatedAt   int64
}

type RequestHandler struct {
	vmhostsInteractor *domain.VmhostsInteractor
}

func convertVmhostsListEntriesToResources(vmhostsListEntries map[string]*domain.VmhostsListEntry) map[string]*vmhostResource {
	vmhostResources := make(map[string]*vmhostResource)
	for k, vmhostsListEntry := range vmhostsListEntries {
		vmhostResources[k] = convertVmhostsListEntryToResource(vmhostsListEntry)
	}
	return vmhostResources
}

func convertVmguestsListEntriesToResources(vmguestsListEntries map[string]*domain.VmguestsListEntry) map[string]*vmguestResource {
	vmguestResources := make(map[string]*vmguestResource)
	for k, vmguestsListEntry := range vmguestsListEntries {
		vmguestResources[k] = convertVmguestsListEntryToResource(vmguestsListEntry)
	}
	return vmguestResources
}

func convertVmhostsListEntryToResource(vmhostsListEntry *domain.VmhostsListEntry) *vmhostResource {
	vmhostResource := new(vmhostResource)
	vmhostResource.Id = vmhostsListEntry.Id
	vmhostResource.DnsName = vmhostsListEntry.DnsName
	vmhostResource.TotalMemory = vmhostsListEntry.TotalMemory
	vmhostResource.Vmguests = convertVmguestsListEntriesToResources(vmhostsListEntry.Vmguests)
	return vmhostResource
}

func convertVmguestsListEntryToResource(vmguestsListEntry *domain.VmguestsListEntry) *vmguestResource {
	vmguestResource := new(vmguestResource)
	vmguestResource.Id = vmguestsListEntry.Id
	vmguestResource.Name = vmguestsListEntry.Name
	vmguestResource.State = vmguestsListEntry.State
	vmguestResource.AllocatedMemory = vmguestsListEntry.AllocatedMemory
	vmguestResource.InfoUpdatedAt = vmguestsListEntry.InfoUpdatedAt.Unix()
	return vmguestResource
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
	router.PathPrefix("/app").Handler(http.StripPrefix("/app", http.FileServer(http.Dir("/usr/share/infmgmt/htdocs"))))

	return router
}

func NewRequestHandler(vmhostsInteractor *domain.VmhostsInteractor) *RequestHandler {
	requestHandler := new(RequestHandler)
	requestHandler.vmhostsInteractor = vmhostsInteractor
	return requestHandler
}

func (rh *RequestHandler) HandleListVmhostsRequest(res http.ResponseWriter, req *http.Request) {
	list, _ := rh.vmhostsInteractor.GetList()
	resources := convertVmhostsListEntriesToResources(list)
	jsonResponse, _ := json.Marshal(resources)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, "%s", jsonResponse)
}

func (rh *RequestHandler) HandleListVmguestsRequest(res http.ResponseWriter, req *http.Request) {
	vmhostId := mux.Vars(req)["vmhostId"]
	list, _ := rh.vmhostsInteractor.GetVmguestsList(vmhostId)
	resources := convertVmguestsListEntriesToResources(list)
	jsonResponse, _ := json.Marshal(resources)
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
