package service

import (
	"fmt"
	"github.com/scipian/eve"
	"github.com/scipian/eve/pkg/config"
	"github.com/scipian/eve/service/health"
	"time"
)

var startTime = time.Now()

type HealthService struct {
}

func (h *HealthService) GetHealth() *eve.HealthInfo {
	sysConfig := config.NewSystemConfig()

	healthInfo := &eve.HealthInfo{
		Hostname: sysConfig.Hostname,
		Metadata: map[string]string{
			"Version":     sysConfig.Version,
			"Environment": sysConfig.Environment,
		},
	}

	errors := make([]eve.Error, 0)

	rethinkChecker := health.NewRethinkdbChecker()
	if err := rethinkChecker.Ping(); err != nil {
		errors = append(errors, *err)
	}
	if err := rethinkChecker.DbReady(); err != nil {
		errors = append(errors, *err)
	}
	if err := rethinkChecker.TableReady(); err != nil {
		errors = append(errors, *err)
	}

	natsChecker := health.NewNatsChecker()
	if err := natsChecker.Ping(); err != nil {
		errors = append(errors, *err)
	}

	vaultChecker, err := health.NewVaultChecker()
	if err != nil {
		errors = append(errors, *err)
	} else {
		if err := vaultChecker.InitStatus(); err != nil {
			errors = append(errors, *err)
		}
		if err := vaultChecker.SealStatus(); err != nil {
			errors = append(errors, *err)
		}
	}

	if len(errors) > 0 {
		healthInfo.Errors = errors
	}

	healthInfo.Uptime = fmt.Sprintf("%s", time.Since(startTime))
	return healthInfo
}

func NewHealthService() *HealthService {
	return &HealthService{}
}
