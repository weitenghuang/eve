package main

import (
	"fmt"
	"github.com/concur/rohr/http"
	"github.com/concur/rohr/http/httprouter"
	"os"
)

func main() {
	port := ":8080"
	if os.Getenv("ROHR_PORT") != "" {
		port = fmt.Sprintf(":%s", os.Getenv("ROHR_PORT"))
	}
	dns := "localhost"
	if os.Getenv("ROHR_DNS") != "" {
		dns = os.Getenv("ROHR_DNS")
	}
	os.Getenv("ROHR_PORT")
	apiServer := http.ApiServer{
		Addr:   port,
		DNS:    dns,
		Router: httprouter.NewRouter(),
	}
	apiServer.ListenAndServe()
}
