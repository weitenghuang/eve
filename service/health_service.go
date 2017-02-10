package service

import (
	"fmt"
	"github.com/concur/rohr"
	"github.com/concur/rohr/pkg/config"
	"github.com/concur/rohr/service/health"
	"time"
)

var startTime = time.Now()

type HealthService struct {
}

func (h *HealthService) GetHealth() *rohr.HealthInfo {
	sysConfig := config.NewSystemConfig()

	healthInfo := &rohr.HealthInfo{
		Hostname: sysConfig.Hostname,
		Metadata: map[string]string{
			"Version":     sysConfig.Version,
			"Environment": sysConfig.Environment,
		},
	}

	errors := make([]rohr.Error, 0)

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

	if len(errors) > 0 {
		healthInfo.Errors = errors
	}

	healthInfo.Uptime = fmt.Sprintf("%s", time.Since(startTime))
	return healthInfo
}

func NewHealthService() *HealthService {
	return &HealthService{}
}
