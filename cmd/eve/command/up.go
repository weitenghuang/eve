package command

import (
	"fmt"
	"github.com/concur/rohr/http"
	"github.com/concur/rohr/http/httprouter"
	"github.com/spf13/cobra"
	"os"
)

const (
	DEFAULT_PORT   = "8088"
	DEFAULT_DNS    = "localhost"
	DEFAULT_SCHEME = "http"
)

var apiServer *http.ApiServer

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "To start eve api server",
	Long:  `To start eve api server`,
	Run: func(cmd *cobra.Command, args []string) {
		apiServer.Router = httprouter.NewRouter()
		apiServer.ListenAndServe()
	},
}

func init() {
	port := os.Getenv("EVE_PORT")
	if port == "" {
		port = DEFAULT_PORT
	}
	dns := os.Getenv("EVE_DNS")
	if dns == "" {
		dns = DEFAULT_DNS
	}
	scheme := os.Getenv("EVE_SCHEME")
	if scheme == "" {
		scheme = DEFAULT_SCHEME
	}
	apiServer = &http.ApiServer{
		Addr:   fmt.Sprintf(":%s", port),
		DNS:    dns,
		Scheme: scheme,
		// Router: httprouter.NewRouter(),
	}
}
