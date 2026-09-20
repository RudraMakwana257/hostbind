package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/RudraMakwana257/hostbind/internal/registry"
)

type Server struct {
	reg *registry.Registry
}

func New(reg *registry.Registry) *Server {
	return &Server{reg: reg}
}

func (p *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse Host (e.g. "api-myproject.localhost:8080" or "api.myproject.localhost:8080")
	host := strings.Split(r.Host, ":")[0]
	
	// Default to removing .localhost
	base := strings.TrimSuffix(host, ".localhost")

	// Support both "service-project" and "service.project" formats
	var serviceName, projectName string
	if strings.Contains(base, ".") {
		parts := strings.SplitN(base, ".", 2)
		serviceName = parts[0]
		projectName = parts[1]
	} else if strings.Contains(base, "-") {
		parts := strings.SplitN(base, "-", 2)
		serviceName = parts[0]
		projectName = parts[1]
	} else {
		http.Error(w, "HostBind Proxy: Invalid host format.\nUse: <service>.<project>.localhost or <service>-<project>.localhost", http.StatusBadRequest)
		return
	}

	allocs, err := p.reg.GetProjectAllocations(projectName, "main")
	if err != nil {
		http.Error(w, "HostBind Proxy: Database error.", http.StatusInternalServerError)
		return
	}

	var targetPort int
	for _, a := range allocs {
		if a.ServiceName == serviceName {
			targetPort = a.Port
			break
		}
	}

	if targetPort == 0 {
		msg := fmt.Sprintf("HostBind Proxy: Service '%s' is not running in project '%s'.", serviceName, projectName)
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	targetURL, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", targetPort))
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	
	// Add helpful headers
	r.Header.Set("X-Forwarded-Host", r.Host)
	proxy.ServeHTTP(w, r)
}
