package service

import (
	"fmt"
	"github.com/concur/rohr"
	"os"
	"time"
)

var startTime = time.Now()

type HealthService struct {
	*rohr.HealthInfo
}

func (h HealthService) GetHealth() (*rohr.HealthInfo, error) {
	h.HealthInfo.Uptime = fmt.Sprintf("%s", time.Since(startTime))
	return h.HealthInfo, nil
}

func NewHealthService() *HealthService {
	return &HealthService{
		HealthInfo: &rohr.HealthInfo{
			Verion:      "v0.0.1",
			Environment: os.Getenv("ENVIRONMENT"),
		},
	}
}
