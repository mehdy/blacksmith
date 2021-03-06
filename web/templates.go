package web

import (
	"fmt"
	"net"
	"net/http"
	"path"

	"github.com/cafebazaar/blacksmith/templating"
)

const (
	templatesDebugTag = "WEB-T"
)

func (ws *webServer) generateTemplateForMachine(templateName string, w http.ResponseWriter, r *http.Request) string {
	_, macStr := path.Split(r.URL.Path)

	mac, err := net.ParseMAC(macStr)
	if err != nil {
		http.Error(w, fmt.Sprintf(`Error while parsing the mac: %q`, err), 500)
		return ""
	}

	machine, exist := ws.ds.GetMachine(mac)
	if !exist {
		http.Error(w, "Machine not found", 500)
		return ""
	}

	cc, err := templating.ExecuteTemplateFolder(
		path.Join(ws.ds.WorkspacePath(), "config", templateName), ws.ds, machine, r.Host)
	if err != nil {
		http.Error(w, fmt.Sprintf(`Error while executing the template: %q`, err), 500)
		return ""
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(cc))

	return cc
}

// Cloudconfig generates and writes cloudconfig for the machine specified by the
// mac in the request url path
func (ws *webServer) Cloudconfig(w http.ResponseWriter, r *http.Request) {
	config := ws.generateTemplateForMachine("cloudconfig", w, r)

	if config != "" && r.FormValue("validate") != "" {
		w.Write([]byte(templating.ValidateCloudConfig(config)))
	}
}

// Ignition generates and writes ignition for the machine specified by the
// mac in the request url path
func (ws *webServer) Ignition(w http.ResponseWriter, r *http.Request) {
	ws.generateTemplateForMachine("ignition", w, r)
}

// Bootparams generates and writes bootparams for the machine specified by the
// mac in the request url path. (Just for validation purpose)
func (ws *webServer) Bootparams(w http.ResponseWriter, r *http.Request) {
	ws.generateTemplateForMachine("bootparams", w, r)
}
