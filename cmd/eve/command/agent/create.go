package agent

import (
	"github.com/concur/rohr"
	"github.com/concur/rohr/http"
	"github.com/concur/rohr/pkg/terraform"
	"github.com/concur/rohr/service"
	"github.com/spf13/cobra"
	"log"
	"runtime"
)

func CreateCmd(stateServer *http.ApiServer) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "To create infrastructure",
		Long:  `To create infrastructure based on user's credentials, quoin module and existing infrastructure state`,
		Run: func(cmd *cobra.Command, args []string) {
			infrastructureService := &service.InfrastructureService{}
			if err := infrastructureService.SubscribeAsyncProc(rohr.CREATE_INFRA, create(infrastructureService, stateServer)); err != nil {
				log.Fatalln(err)
			}
			log.Printf("Listening on [%s]\n", rohr.CREATE_INFRA)
			runtime.Goexit()
		},
	}
}

func create(infraSvc rohr.InfrastructureService, stateServer *http.ApiServer) rohr.InfrastructureAsyncHandler {
	var toFailStatus = func(name string) {
		if statusErr := infraSvc.UpdateInfrastructureStatus(name, rohr.FAILED); statusErr != nil {
			log.Println(statusErr)
		}
	}
	return func(infra *rohr.Infrastructure) {
		if infra == nil {
			log.Println("Empty infrastructure object detected.")
			return
		}
		log.Printf("Start infrastructure creation process for %s.\n", infra.Name)
		if err := infraSvc.UpdateInfrastructureStatus(infra.Name, rohr.RUNNING); err != nil {
			log.Println(err)
		}
		id := infraSvc.GetQuoinArchiveIdFromUri(infra.Quoin.ArchiveUri)
		quoinArchive, err := infraSvc.GetQuoinArchive(id)
		if err != nil {
			toFailStatus(infra.Name)
			log.Println(err)
			return
		}
		if quoinArchive == nil {
			toFailStatus(infra.Name)
			log.Println("Invalid Quoin Archive Id: ", id)
			return
		}
		log.Println("Infrastructure", infra.Name, "gets Quoin Archive:", id, quoinArchive.QuoinName)
		varfile := createVarFile(infra.Variables)
		remoteState := stateEndpoint(stateServer, infra.Name)
		if err := terraform.ApplyQuoin(infra.Name, quoinArchive.Modules, varfile, remoteState); err != nil {
			toFailStatus(infra.Name)
			log.Println(err)
			return
		}
		if err := infraSvc.UpdateInfrastructureStatus(infra.Name, rohr.DEPLOYED); err != nil {
			log.Println(err)
			return
		}
		log.Println("Creation Done!")
	}
}
