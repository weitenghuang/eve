package vault

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	eveHttp "github.com/concur/eve/http"
	"github.com/hashicorp/vault/api"
	"net/http"
	"net/http/httputil"
	"path"
	"sort"
	"strings"
)

const (
	PKI_STR = ":pki"
)

var routes = []string{
	"/sys/mounts",
	fmt.Sprint("/", PKI_STR, "/root/generate/internal"),
	fmt.Sprint("/", PKI_STR, "/ca_chain"),
	fmt.Sprint("/", PKI_STR, "/ca"),
	fmt.Sprint("/", PKI_STR, "/ca/pem"),
	fmt.Sprint("/", PKI_STR, "/config/urls"),
	fmt.Sprint("/", PKI_STR, "/roles"),
	fmt.Sprint("/", PKI_STR, "/issue"),
	fmt.Sprint("/", PKI_STR, "/config/ca"),
}

func init() {
	sort.Strings(routes)
}

type VaultHandler struct {
	Proxy *httputil.ReverseProxy
}

func (v *VaultHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req, err := eveHttp.Authentication(rw, req)
	if err != nil {
		log.Errorln("Authentication error", err)
		return
	}
	if err := routeFilter(rw, req); err != nil {
		log.Errorln("Invalid request error", err)
		return
	}
	log.Infoln("ServeHTTP: ", *req)
	v.Proxy.ServeHTTP(rw, req)
}

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) RegisterRoute(apiServer *eveHttp.ApiServer) eveHttp.Handler {
	vaultHandler := &VaultHandler{Proxy: &httputil.ReverseProxy{}}
	reverseProxy := vaultHandler.Proxy

	config := api.DefaultConfig()
	err := config.ReadEnvironment()
	if err != nil {
		log.Errorf("error reading Vault environment: %v", err)
		return reverseProxy
	}

	reverseProxy.Director = vaultRequestDirector(reverseProxy, config)

	log.Infof("TLSClientConfig: %#v ", config.HttpClient.Transport.(*http.Transport).TLSClientConfig)
	log.Infof("vault config: %#v", config)

	reverseProxy.Transport = config.HttpClient.Transport //customTransport
	reverseProxy.ModifyResponse = vaultResponseModifier
	return vaultHandler
}

func vaultRequestDirector(rProxy *httputil.ReverseProxy, config *api.Config) func(*http.Request) {
	target_director := rProxy.Director

	client, err := api.NewClient(config)
	if err != nil {
		log.Errorf("error creating Vault client: %v", err)
		return target_director
	}

	return func(req *http.Request) {
		log.Infoln("Start vault ServeHTTP: ", req.URL.Path)
		req.URL.Path = path.Join("/v1", req.URL.Path[len(eveHttp.VAULT)+1:])
		log.Infof("vaultRequest Path: %#v", req.URL.Path)

		vRequest := client.NewRequest(req.Method, req.URL.Path)
		log.Infof("vaultRequest Token: %v", vRequest.ClientToken)

		var data map[string]interface{}
		if req.Body != nil && http.MethodGet != req.Method && http.MethodDelete != req.Method {
			if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
				log.Errorf("error parsing Vault request body: %v", err)
				target_director(req)
			}
			vRequest.SetJSONBody(data)
			log.Infof("vaultRequest Data: %#v", data)
		}

		httpReq, err := vRequest.ToHTTP()
		if err != nil {
			log.Errorf("error converting Vault request: %v", err)
			target_director(req)
		}
		*req = *httpReq
	}
}

func vaultResponseModifier(resp *http.Response) error {
	log.Infoln("vaultResponseModifier: ", *resp)
	return nil
}

func routeFilter(w http.ResponseWriter, r *http.Request) error {
	route := r.URL.Path[len(eveHttp.VAULT)+1:]
	dirs := strings.Split(route, "/")
	var base string
	if len(dirs) > 2 {
		base = dirs[1]
		if base != "sys" {
			route = strings.Replace(route, base, PKI_STR, 1)
		}
		ind := sort.SearchStrings(routes, route)
		// Check if req path exactly matches available route
		if ind < len(routes) && route == routes[ind] {
			return nil
		}

		// Check if req path's dirctory matches available route
		if ind > 0 && path.Dir(route) == routes[ind-1] {
			return nil
		}
	}
	err := fmt.Errorf("Invalid Vault API Request Path.")
	http.Error(w, err.Error(), http.StatusBadRequest)
	return err
}
